package domain

import (
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

// FilterPredicate is a function that returns true if a flight passes the filter.
type FilterPredicate func(*entity.Flight) bool

// BuildFilterPredicates builds a list of filter predicates from the given SearchFilter.
func BuildFilterPredicates(filter *entity.SearchFilter) []FilterPredicate {
	var predicates []FilterPredicate

	if filter == nil {
		return predicates
	}

	if filter.MinPrice != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Price >= *filter.MinPrice
		})
	}

	if filter.MaxPrice != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Price <= *filter.MaxPrice
		})
	}

	if filter.MaxStops != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Stops <= *filter.MaxStops
		})
	}

	if filter.DepartureStart != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return !f.DepartureTime.Before(*filter.DepartureStart)
		})
	}

	if filter.DepartureEnd != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return !f.DepartureTime.After(*filter.DepartureEnd)
		})
	}

	if filter.ArrivalStart != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return !f.ArrivalTime.Before(*filter.ArrivalStart)
		})
	}

	if filter.ArrivalEnd != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return !f.ArrivalTime.After(*filter.ArrivalEnd)
		})
	}

	if filter.MaxDuration != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Duration <= *filter.MaxDuration
		})
	}

	if len(filter.AirlineCodes) > 0 {
		airlineMap := make(map[string]bool, len(filter.AirlineCodes))
		for _, code := range filter.AirlineCodes {
			airlineMap[code] = true
		}
		predicates = append(predicates, func(f *entity.Flight) bool {
			return airlineMap[f.Provider]
		})
	}

	if filter.CabinClass != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.CabinClass == *filter.CabinClass
		})
	}

	return predicates
}

// ApplyFilters filters flights in-place using the given predicates.
func ApplyFilters(flights []*entity.Flight, predicates []FilterPredicate) []*entity.Flight {
	if len(predicates) == 0 {
		return flights
	}

	n := 0
	for _, f := range flights {
		keep := true
		for _, predicate := range predicates {
			if !predicate(f) {
				keep = false
				break
			}
		}
		if keep {
			flights[n] = f
			n++
		}
	}

	for i := n; i < len(flights); i++ {
		flights[i] = nil
	}

	return flights[:n]
}
