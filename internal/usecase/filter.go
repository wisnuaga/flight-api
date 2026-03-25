package usecase

import (
	"github.com/wisnuaga/flight-api/internal/domain"
)

type FilterPredicate func(*domain.Flight) bool

func BuildFilterPredicates(filter *domain.SearchFilter) []FilterPredicate {
	var predicates []FilterPredicate

	if filter == nil {
		return predicates
	}

	if filter.MinPrice != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return f.Price >= *filter.MinPrice
		})
	}

	if filter.MaxPrice != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return f.Price <= *filter.MaxPrice
		})
	}

	if filter.MaxStops != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return f.Stops <= *filter.MaxStops
		})
	}

	if filter.DepartureStart != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return !f.DepartureTime.Before(*filter.DepartureStart)
		})
	}

	if filter.DepartureEnd != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return !f.DepartureTime.After(*filter.DepartureEnd)
		})
	}

	if filter.ArrivalStart != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return !f.ArrivalTime.Before(*filter.ArrivalStart)
		})
	}

	if filter.ArrivalEnd != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return !f.ArrivalTime.After(*filter.ArrivalEnd)
		})
	}

	if filter.MaxDuration != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return f.Duration <= *filter.MaxDuration
		})
	}

	if len(filter.AirlineCodes) > 0 {
		airlineMap := make(map[string]bool, len(filter.AirlineCodes))
		for _, code := range filter.AirlineCodes {
			airlineMap[code] = true
		}
		predicates = append(predicates, func(f *domain.Flight) bool {
			return airlineMap[f.Provider]
		})
	}

	if filter.CabinClass != nil {
		predicates = append(predicates, func(f *domain.Flight) bool {
			return f.CabinClass == *filter.CabinClass
		})
	}

	return predicates
}

func ApplyFilters(flights []*domain.Flight, predicates []FilterPredicate) []*domain.Flight {
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
