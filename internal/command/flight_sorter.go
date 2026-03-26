package command

import (
	"sort"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/port"
)

type flightSortCommandImpl struct{}

// NewFlightSortCommand creates a new FlightSortCommand instance.
func NewFlightSortCommand() port.FlightSortCommand {
	return &flightSortCommandImpl{}
}

// Execute sorts flights in-place according to sortParam.
func (c *flightSortCommandImpl) Execute(flights []*entity.Flight, sortParam entity.SearchSort) {
	field := sortParam.Field
	if field == "" {
		field = entity.SortByPrice
	}

	order := sortParam.Order
	if order == "" {
		order = entity.SortAsc
	}

	// Compute boundaries for best-value normalisation
	var minPrice, maxPrice float64
	var minDuration, maxDuration int64

	if field == entity.SortByBestValue && len(flights) > 0 {
		minPrice = flights[0].Price
		maxPrice = flights[0].Price
		minDuration = int64(flights[0].Duration)
		maxDuration = int64(flights[0].Duration)

		for _, f := range flights {
			if f.Price < minPrice {
				minPrice = f.Price
			}
			if f.Price > maxPrice {
				maxPrice = f.Price
			}
			fDur := int64(f.Duration)
			if fDur < minDuration {
				minDuration = fDur
			}
			if fDur > maxDuration {
				maxDuration = fDur
			}
		}
	}

	priceWeight := sortParam.PriceWeight
	if priceWeight == 0 {
		priceWeight = 1.0
	}

	durationWeight := sortParam.DurationWeight
	if durationWeight == 0 {
		durationWeight = 1.0
	}

	sort.SliceStable(flights, func(i, j int) bool {
		f1, f2 := flights[i], flights[j]

		var isLess bool
		switch field {
		case entity.SortByPrice:
			isLess = f1.Price < f2.Price
		case entity.SortByDuration:
			isLess = f1.Duration < f2.Duration
		case entity.SortByDeparture:
			isLess = f1.DepartureTime.Before(f2.DepartureTime)
		case entity.SortByArrival:
			isLess = f1.ArrivalTime.Before(f2.ArrivalTime)
		case entity.SortByBestValue:
			score1 := c.calculateBestValueScore(f1, minPrice, maxPrice, minDuration, maxDuration, priceWeight, durationWeight)
			score2 := c.calculateBestValueScore(f2, minPrice, maxPrice, minDuration, maxDuration, priceWeight, durationWeight)
			isLess = score1 < score2
		default:
			isLess = f1.Price < f2.Price
		}

		if order == entity.SortDesc {
			return !isLess
		}
		return isLess
	})
}

// calculateBestValueScore returns a normalised 0-1 score for a flight (lower = better value).
func (c *flightSortCommandImpl) calculateBestValueScore(f *entity.Flight, minPrice, maxPrice float64, minDur, maxDur int64, weightPrice, weightDur float64) float64 {
	priceRange := maxPrice - minPrice
	durRange := float64(maxDur - minDur)

	var normPrice float64
	if priceRange > 0 {
		normPrice = (f.Price - minPrice) / priceRange
	}

	var normDur float64
	if durRange > 0 {
		normDur = float64(int64(f.Duration)-minDur) / durRange
	}

	return (normPrice * weightPrice) + (normDur * weightDur)
}
