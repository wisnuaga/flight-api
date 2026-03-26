package command

import (
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

// RoundTripCombiner combines outbound and return flights into round-trip itineraries
type RoundTripCombiner struct {
	minLayoverMinutes  int // Minimum connection time between outbound arrival and return departure
	maxLayoverMinutes  int // Maximum connection time between outbound arrival and return departure
	preferSameAirline  bool
	preferSameTerminal bool
}

// NewRoundTripCombiner creates a new round-trip combiner with default settings
func NewRoundTripCombiner() *RoundTripCombiner {
	return &RoundTripCombiner{
		minLayoverMinutes:  90,   // Default minimum 90 minutes layover
		maxLayoverMinutes:  1440, // Default maximum 24 hours layover
		preferSameAirline:  false,
		preferSameTerminal: false,
	}
}

// Combine creates round-trip itineraries from outbound and return flights
// Returns only valid combinations where layover time is within acceptable range
func (rc *RoundTripCombiner) Combine(outboundFlights, returnFlights []*entity.Flight) []*entity.RoundTripItinerary {
	var itineraries []*entity.RoundTripItinerary

	for _, outbound := range outboundFlights {
		for _, returnFlight := range returnFlights {
			// Check if combination is valid (layover time constraints)
			if !rc.isValidCombination(outbound, returnFlight) {
				continue
			}

			// Calculate total price and duration
			totalPrice := outbound.Price.Add(returnFlight.Price)
			totalDuration := outbound.TotalTripDuration() + returnFlight.TotalTripDuration()

			itinerary := &entity.RoundTripItinerary{
				OutboundFlight: outbound,
				ReturnFlight:   returnFlight,
				TotalPrice:     totalPrice,
				TotalDuration:  totalDuration,
			}

			itineraries = append(itineraries, itinerary)
		}
	}

	return itineraries
}

// isValidCombination checks if an outbound-return flight pair is valid
// Validates layover time between outbound arrival and return departure
func (rc *RoundTripCombiner) isValidCombination(outbound, returnFlight *entity.Flight) bool {
	// Outbound arrival time
	outboundArrival := outbound.Destination.Time

	// Return departure time
	returnDeparture := returnFlight.Origin.Time

	// Calculate layover duration
	layoverDuration := returnDeparture.Sub(outboundArrival)
	layoverMinutes := int(layoverDuration.Minutes())

	// Check minimum layover constraint
	if layoverMinutes < rc.minLayoverMinutes {
		return false
	}

	// Check maximum layover constraint
	if layoverMinutes > rc.maxLayoverMinutes {
		return false
	}

	return true
}
