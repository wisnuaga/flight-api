package domain

import "time"

type SearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate time.Time
	Passengers    int
	CabinClass    string
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
