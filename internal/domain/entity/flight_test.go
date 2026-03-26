package entity_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

func TestNormalizeFlight(t *testing.T) {
	dep := time.Now()
	arr := dep.Add(2 * time.Hour)

	f := entity.Flight{
		DepartureTime:  dep,
		ArrivalTime:    arr,
		CabinClass:     "economy",
		AvailableSeats: 1,
	}

	norm := entity.NormalizeFlight(f)
	if norm.Duration != 2*time.Hour {
		t.Errorf("Expected duration 2h, got %v", norm.Duration)
	}
}

func TestIsValidFlight_Invalid(t *testing.T) {
	dep := time.Now()

	f := entity.Flight{
		DepartureTime: dep,
		ArrivalTime:   dep.Add(-2 * time.Hour), // arrival before departure
		Price:         decimal.NewFromInt(500),
	}

	if entity.IsValidFlight(f) {
		t.Errorf("Expected negative duration to fail validation")
	}
}

func TestIsValidFlight_Valid(t *testing.T) {
	dep := time.Now()

	f := entity.Flight{
		Origin:        "CGK",
		Destination:   "DPS",
		FlightNumber:  "FL100",
		DepartureTime: dep,
		ArrivalTime:   dep.Add(2 * time.Hour),
		Duration:      2 * time.Hour,
		Price:         decimal.NewFromInt(500),
	}

	if !entity.IsValidFlight(f) {
		t.Errorf("Failed: Origin=%v, Dest=%v, Price=%v, Dur=%v",
			f.Origin, f.Destination, f.Price, f.Duration)
	}
}
