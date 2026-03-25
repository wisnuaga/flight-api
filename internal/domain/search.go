package domain

import "time"

type SearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate time.Time
	Passengers    int

	Filter SearchFilter
	Sort   SearchSort
}

type SearchFilter struct {
	// Price filters
	MinPrice *float64
	MaxPrice *float64

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
	Flights []*Flight
	Meta    *SearchMeta
}

type SearchMeta struct {
	TotalFlights int
	Providers    int
	SuccessCount int
	FailedCount  int
}
