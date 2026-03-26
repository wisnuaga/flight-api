package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Flight struct {
	// Basic info
	ID           string
	Airline      AirlineName // The airline/provider name (e.g., "Garuda Indonesia")
	FlightNumber string
	AirlineCode  string // IATA airline code (e.g., "GA" for Garuda)

	// Route with timezone-aware locations
	Origin      Location
	Destination Location

	// Schedule and pricing
	Duration       time.Duration
	Price          decimal.Decimal
	Currency       string
	CabinClass     string
	AvailableSeats int

	// Routing info
	Stops    int
	Layovers []*Layover
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
// and recomputing duration from departure/arrival times (UTC-based).
func NormalizeFlight(f Flight) Flight {
	if f.CabinClass == "" {
		f.CabinClass = "economy"
	}

	if f.AvailableSeats == 0 {
		f.AvailableSeats = 1 // minimum default
	}

	// Ensure times are in UTC for consistent handling across providers
	if !f.Origin.Time.IsZero() {
		f.Origin.Time = f.Origin.Time.UTC()
	}
	if !f.Destination.Time.IsZero() {
		f.Destination.Time = f.Destination.Time.UTC()
	}

	// Default timezone to UTC if not set
	if f.Origin.Timezone == nil {
		f.Origin.Timezone = time.UTC
	}
	if f.Destination.Timezone == nil {
		f.Destination.Timezone = time.UTC
	}

	// Compute duration from UTC times
	if !f.Destination.Time.IsZero() && !f.Origin.Time.IsZero() {
		f.Duration = f.Destination.Time.Sub(f.Origin.Time)
	}

	// Apply entity-level normalization (uppercase codes, default currency)
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

	// Both times must be present and valid
	if f.Origin.Time.IsZero() || f.Destination.Time.IsZero() {
		return false
	}

	// Arrival must be after departure
	if !f.Destination.Time.After(f.Origin.Time) {
		return false
	}

	// Duration must be positive
	if f.Duration <= 0 {
		return false
	}

	return true
}

// TotalTripDuration returns sum of all in-air duration + layovers
func (f *Flight) TotalTripDuration() time.Duration {
	total := f.Duration
	for _, l := range f.Layovers {
		total += l.Duration
	}
	return total
}
