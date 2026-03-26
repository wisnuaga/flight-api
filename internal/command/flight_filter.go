package command

import (
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/port"
)

type flightFilterCommandImpl struct{}

// NewFlightFilterCommand creates a new FlightFilterCommand instance.
func NewFlightFilterCommand() port.FlightFilterCommand {
	return &flightFilterCommandImpl{}
}

type filterPredicate func(*entity.Flight) bool

func (c *flightFilterCommandImpl) buildPredicates(filter *entity.SearchFilter) []filterPredicate {
	var predicates []filterPredicate

	if filter == nil {
		return predicates
	}

	if filter.MinPrice != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Price.GreaterThanOrEqual(*filter.MinPrice)
		})
	}

	if filter.MaxPrice != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Price.LessThanOrEqual(*filter.MaxPrice)
		})
	}

	if filter.MaxStops != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			return f.Stops <= *filter.MaxStops
		})
	}

	if filter.DepartureStart != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			// Reject zero times — they would compare as before any real start boundary
			if f.Origin.Time.IsZero() {
				return false
			}
			// Use UTC time from Origin for consistent timezone-safe comparison
			return !f.Origin.Time.Before(*filter.DepartureStart)
		})
	}

	if filter.DepartureEnd != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			// Reject zero times — they have no meaningful departure to compare
			if f.Origin.Time.IsZero() {
				return false
			}
			// Use UTC time from Origin for consistent timezone-safe comparison
			return !f.Origin.Time.After(*filter.DepartureEnd)
		})
	}

	if filter.ArrivalStart != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			// Reject zero times — they would compare as before any real start boundary
			if f.Destination.Time.IsZero() {
				return false
			}
			// Use UTC time from Destination for consistent timezone-safe comparison
			return !f.Destination.Time.Before(*filter.ArrivalStart)
		})
	}

	if filter.ArrivalEnd != nil {
		predicates = append(predicates, func(f *entity.Flight) bool {
			// Reject zero times — they have no meaningful arrival to compare
			if f.Destination.Time.IsZero() {
				return false
			}
			// Use UTC time from Destination for consistent timezone-safe comparison
			return !f.Destination.Time.After(*filter.ArrivalEnd)
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

// Execute filters flights in-place using the given filter parameters.
func (c *flightFilterCommandImpl) Execute(flights []*entity.Flight, filter *entity.SearchFilter) []*entity.Flight {
	predicates := c.buildPredicates(filter)
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
