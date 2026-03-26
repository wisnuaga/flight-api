package port

import "github.com/wisnuaga/flight-api/internal/domain/entity"

// FlightFilterCommand encapsulates the flight filtering logic.
type FlightFilterCommand interface {
	Execute(flights []*entity.Flight, filter *entity.SearchFilter) []*entity.Flight
}

// FlightSortCommand encapsulates the flight sorting logic.
type FlightSortCommand interface {
	Execute(flights []*entity.Flight, sortParam entity.SearchSort)
}

// RoundTripCombinerCommand encapsulates the round-trip combination logic.
type RoundTripCombinerCommand interface {
	Combine(outboundFlights, returnFlights []*entity.Flight) []*entity.RoundTripItinerary
}
