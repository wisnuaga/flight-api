package domain

import (
	"strings"
	"time"
)

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

	// Routing info
	Stops int

	// Pricing
	Price    float64
	Currency string
}

func (f *Flight) Normalize() {
	f.Origin = strings.ToUpper(f.Origin)
	f.Destination = strings.ToUpper(f.Destination)

	if f.Currency == "" {
		f.Currency = "IDR"
	}
}
