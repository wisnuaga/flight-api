package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Flight struct {
	// Basic info
	ID           string
	Provider     string
	FlightNumber string

	// Route
	Origin      Location
	Destination Location

	// Schedule
	DepartureTime  time.Time
	ArrivalTime    time.Time
	Duration       time.Duration
	Price          decimal.Decimal
	Currency       string
	CabinClass     string
	AvailableSeats int

	// Routing info
	Stops int
}

// Normalize applies basic field normalisations on the Flight value.
func (f *Flight) Normalize() {
	f.Origin.Normalize()
	f.Destination.Normalize()

	if f.Currency == "" {
		f.Currency = "IDR"
	}
}

// NormalizeFlight returns a fully normalised copy of f, filling in defaults
// and recomputing duration from departure/arrival times.
func NormalizeFlight(f Flight) Flight {
	if f.CabinClass == "" {
		f.CabinClass = "economy"
	}

	if f.AvailableSeats == 0 {
		f.AvailableSeats = 1 // minimum default
	}

	// Always store times in UTC for consistent handling across providers
	if !f.DepartureTime.IsZero() {
		f.DepartureTime = f.DepartureTime.UTC()
	}
	if !f.ArrivalTime.IsZero() {
		f.ArrivalTime = f.ArrivalTime.UTC()
	}

	// Compute duration from times
	if !f.ArrivalTime.IsZero() && !f.DepartureTime.IsZero() {
		f.Duration = f.ArrivalTime.Sub(f.DepartureTime)
	}

	// Apply entity-level normalisation (uppercase codes, default currency)
	f.Normalize()

	return f
}

// IsValidFlight validates that a Flight meets minimum business rules.
func IsValidFlight(f Flight) bool {
	if f.Origin.Airport == "" || f.Destination.Airport == "" {
		return false
	}

	if f.Price.LessThanOrEqual(decimal.Zero) {
		return false
	}

	if f.DepartureTime.IsZero() || f.ArrivalTime.IsZero() {
		return false
	}

	if !f.ArrivalTime.After(f.DepartureTime) {
		return false
	}

	if f.Duration <= 0 {
		return false
	}

	return true
}
