package usecase

import (
	"sort"

	"github.com/wisnuaga/flight-api/internal/domain"
)

func ApplySorting(flights []*domain.Flight, sortParam domain.SearchSort) {
	field := sortParam.Field
	if field == "" {
		field = domain.SortByPrice
	}

	order := sortParam.Order
	if order == "" {
		order = domain.SortAsc
	}

	// Compute max/min boundaries for normalization when calculating best value
	var minPrice, maxPrice float64
	var minDuration, maxDuration int64 // store duration in int64 nanos or millis for diffs

	if field == domain.SortByBestValue && len(flights) > 0 {
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
		priceWeight = 1.0 // default
	}

	durationWeight := sortParam.DurationWeight
	if durationWeight == 0 {
		durationWeight = 1.0 // default
	}

	sort.SliceStable(flights, func(i, j int) bool {
		f1, f2 := flights[i], flights[j]

		var isLess bool
		switch field {
		case domain.SortByPrice:
			isLess = f1.Price < f2.Price
		case domain.SortByDuration:
			isLess = f1.Duration < f2.Duration
		case domain.SortByDeparture:
			isLess = f1.DepartureTime.Before(f2.DepartureTime)
		case domain.SortByArrival:
			isLess = f1.ArrivalTime.Before(f2.ArrivalTime)
		case domain.SortByBestValue:
			score1 := calculateBestValueScore(f1, minPrice, maxPrice, minDuration, maxDuration, priceWeight, durationWeight)
			score2 := calculateBestValueScore(f2, minPrice, maxPrice, minDuration, maxDuration, priceWeight, durationWeight)
			isLess = score1 < score2
		default:
			isLess = f1.Price < f2.Price
		}

		if order == domain.SortDesc {
			return !isLess
		}
		return isLess
	})
}

// calculateBestValueScore gives a 0-1 mapped score for price and duration. Lowest is best.
func calculateBestValueScore(f *domain.Flight, minPrice, maxPrice float64, minDur, maxDur int64, weightPrice, weightDur float64) float64 {
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
