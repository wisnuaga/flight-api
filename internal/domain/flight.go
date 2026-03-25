package domain

import "time"

type SearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate time.Time
	Passengers    int
	CabinClass    string
}

type Flight struct {
	// Basic info
	ID           string
	Provider     string
	FlightNumber string

	// Route
	Origin      string
	Destination string

	// Schedule
	DepartureTime time.Time
	ArrivalTime   time.Time
	Duration      time.Duration

	// Seat info
	CabinClass     string
	AvailableSeats int

	// Pricing
	Price    float64
	Currency string
}
