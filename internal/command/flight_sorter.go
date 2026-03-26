package command

import (
	"sort"

	"github.com/shopspring/decimal"
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
	var minPrice, maxPrice decimal.Decimal
	var minDuration, maxDuration int64

	if field == entity.SortByBestValue && len(flights) > 0 {
		minPrice = flights[0].Price
		maxPrice = flights[0].Price
		minDuration = int64(flights[0].TotalTripDuration())
		maxDuration = int64(flights[0].TotalTripDuration())

		for _, f := range flights {
			if f.Price.LessThan(minPrice) {
				minPrice = f.Price
			}
			if f.Price.GreaterThan(maxPrice) {
				maxPrice = f.Price
			}
			fDur := int64(f.TotalTripDuration())
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
			isLess = f1.Price.LessThan(f2.Price)
		case entity.SortByDuration:
			isLess = f1.TotalTripDuration() < f2.TotalTripDuration()
		case entity.SortByDeparture:
			// Push zero-time flights to the end; use UTC for timezone-safe comparison
			if f1.Origin.Time.IsZero() {
				isLess = false
			} else if f2.Origin.Time.IsZero() {
				isLess = true
			} else {
				isLess = f1.Origin.Time.Before(f2.Origin.Time)
			}
		case entity.SortByArrival:
			// Push zero-time flights to the end; use UTC for timezone-safe comparison
			if f1.Destination.Time.IsZero() {
				isLess = false
			} else if f2.Destination.Time.IsZero() {
				isLess = true
			} else {
				isLess = f1.Destination.Time.Before(f2.Destination.Time)
			}
		case entity.SortByBestValue:
			score1 := c.calculateBestValueScore(f1, minPrice, maxPrice, minDuration, maxDuration, priceWeight, durationWeight)
			score2 := c.calculateBestValueScore(f2, minPrice, maxPrice, minDuration, maxDuration, priceWeight, durationWeight)
			isLess = score1 < score2
		default:
			isLess = f1.Price.LessThan(f2.Price)
		}

		if order == entity.SortDesc {
			return !isLess
		}
		return isLess
	})
}

// calculateBestValueScore returns a normalised 0-1 score for a flight (lower = better value).
func (c *flightSortCommandImpl) calculateBestValueScore(f *entity.Flight, minPrice, maxPrice decimal.Decimal, minDur, maxDur int64, weightPrice, weightDur float64) float64 {
	priceRange := maxPrice.Sub(minPrice)
	durRange := float64(maxDur - minDur)

	var normPrice float64
	if priceRange.GreaterThan(decimal.Zero) {
		normPrice = f.Price.Sub(minPrice).Div(priceRange).InexactFloat64()
	}

	var normDur float64
	if durRange > 0 {
		normDur = float64(int64(f.TotalTripDuration())-minDur) / durRange
	}

	return (normPrice * weightPrice) + (normDur * weightDur)
}
