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
		default:
			isLess = f1.Price < f2.Price
		}

		if order == domain.SortDesc {
			return !isLess
		}
		return isLess
	})
}
