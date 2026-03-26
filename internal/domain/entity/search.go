package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type SearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate time.Time
	ReturnDate    *time.Time // nil for one-way, set for round-trip
	Passengers    int

	Filter SearchFilter
	Sort   SearchSort
}

type SearchFilter struct {
	// Price filters
	MinPrice *decimal.Decimal
	MaxPrice *decimal.Decimal

	// Stops filters
	MaxStops *int

	// Time filters
	DepartureStart *time.Time
	DepartureEnd   *time.Time
	ArrivalStart   *time.Time
	ArrivalEnd     *time.Time
	MaxDuration    *time.Duration

	// Airline filters
	AirlineCodes []string

	// Cabin class filters
	CabinClass *string
}

type SearchResult struct {
	Flights              []*Flight
	RoundTripItineraries []*RoundTripItinerary
	Meta                 *SearchMeta
}

// RoundTripItinerary represents a complete round-trip itinerary with outbound and return flights
type RoundTripItinerary struct {
	OutboundFlight *Flight
	ReturnFlight   *Flight
	TotalPrice     decimal.Decimal
	TotalDuration  time.Duration
}

type SearchMeta struct {
	TotalFlights int
	Providers    int
	SuccessCount int
	FailedCount  int
	SearchTimeMs int
	CacheHit     bool
}
