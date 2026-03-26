package entity_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

func TestNormalizeFlight(t *testing.T) {
	dep := time.Now().UTC()
	arr := dep.Add(2 * time.Hour)

	f := entity.Flight{
		Origin: entity.Location{
			Airport: "CGK",
			Time:    dep,
		},
		Destination: entity.Location{
			Airport: "DPS",
			Time:    arr,
		},
		Price:          decimal.NewFromInt(500),
		CabinClass:     "economy",
		AvailableSeats: 1,
	}

	norm := entity.NormalizeFlight(f)
	if norm.Duration != 2*time.Hour {
		t.Errorf("Expected duration 2h, got %v", norm.Duration)
	}
}

func TestIsValidFlight_Invalid(t *testing.T) {
	dep := time.Now().UTC()

	f := entity.Flight{
		Origin: entity.Location{
			Airport: "CGK",
			Time:    dep,
		},
		Destination: entity.Location{
			Airport: "DPS",
			Time:    dep.Add(-2 * time.Hour), // arrival before departure
		},
		Price: decimal.NewFromInt(500),
	}

	if entity.IsValidFlight(f) {
		t.Errorf("Expected negative duration to fail validation")
	}
}

func TestIsValidFlight_Valid(t *testing.T) {
	dep := time.Now().UTC()

	f := entity.Flight{
		Origin: entity.Location{
			Airport: "CGK",
			Time:    dep,
		},
		Destination: entity.Location{
			Airport:  "DPS",
			Time:     dep.Add(2 * time.Hour),
			Timezone: time.UTC,
		},
		FlightNumber: "FL100",
		Duration:     2 * time.Hour,
		Price:        decimal.NewFromInt(500),
	}

	if !entity.IsValidFlight(f) {
		t.Errorf("Failed: Origin=%v, Dest=%v, Price=%v, Dur=%v",
			f.Origin, f.Destination, f.Price, f.Duration)
	}
}

func TestNormalizeFlight_DefaultsNilTimezoneToUTC(t *testing.T) {
	dep := time.Date(2025, 1, 1, 6, 0, 0, 0, time.UTC)
	arr := dep.Add(2 * time.Hour)

	f := entity.Flight{
		Origin:         entity.Location{Airport: "CGK", Time: dep}, // Timezone is nil
		Destination:    entity.Location{Airport: "DPS", Time: arr}, // Timezone is nil
		Price:          decimal.NewFromInt(500),
		CabinClass:     "economy",
		AvailableSeats: 1,
	}

	norm := entity.NormalizeFlight(f)

	if norm.Origin.Timezone != time.UTC {
		t.Errorf("expected Origin.Timezone = UTC when nil, got %v", norm.Origin.Timezone)
	}
	if norm.Destination.Timezone != time.UTC {
		t.Errorf("expected Destination.Timezone = UTC when nil, got %v", norm.Destination.Timezone)
	}
}

func TestNormalizeFlight_PreservesNonUTCTimezone(t *testing.T) {
	jkt, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("failed to load Asia/Jakarta: %v", err)
	}

	// Provider supplies times already in UTC, with the original TZ stored separately
	depUTC := time.Date(2025, 1, 1, 6, 0, 0, 0, time.UTC) // 06:00 UTC = 13:00 WIB
	arrUTC := depUTC.Add(2 * time.Hour)

	f := entity.Flight{
		Origin: entity.Location{
			Airport:  "CGK",
			Time:     depUTC,
			Timezone: jkt,
		},
		Destination: entity.Location{
			Airport:  "DPS",
			Time:     arrUTC,
			Timezone: jkt,
		},
		Price:          decimal.NewFromInt(1250000),
		CabinClass:     "economy",
		AvailableSeats: 28,
	}

	norm := entity.NormalizeFlight(f)

	// Times must remain in UTC after normalization
	if norm.Origin.Time != depUTC {
		t.Errorf("Origin.Time changed: got %v, want %v", norm.Origin.Time, depUTC)
	}
	// Original timezone must be preserved
	if norm.Origin.Timezone.String() != "Asia/Jakarta" {
		t.Errorf("Origin.Timezone overwritten: got %v, want Asia/Jakarta", norm.Origin.Timezone)
	}
	// Duration computed from UTC times
	if norm.Duration != 2*time.Hour {
		t.Errorf("Duration mismatch: got %v, want 2h", norm.Duration)
	}
	// Local datetime must reflect WIB offset (+07:00)
	localDep := norm.Origin.Time.In(norm.Origin.Timezone)
	if localDep.Hour() != 13 {
		t.Errorf("Expected local departure hour 13 (WIB), got %d", localDep.Hour())
	}
}
