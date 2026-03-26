package command

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

func TestNewRoundTripCombiner(t *testing.T) {
	combiner := NewRoundTripCombiner()
	if combiner == nil {
		t.Fatal("NewRoundTripCombiner returned nil")
	}
	if combiner.minLayoverMinutes != 90 {
		t.Errorf("minLayoverMinutes = %d, want 90", combiner.minLayoverMinutes)
	}
	if combiner.maxLayoverMinutes != 1440 {
		t.Errorf("maxLayoverMinutes = %d, want 1440", combiner.maxLayoverMinutes)
	}
}

func TestRoundTripCombiner_Combine(t *testing.T) {
	combiner := NewRoundTripCombiner()

	// Create sample flights
	now := time.Now().UTC()

	outbound := &entity.Flight{
		ID:             "out1",
		FlightNumber:   "GA101",
		Provider:       "Garuda",
		AirlineCode:    "GA",
		Stops:          0,
		Price:          decimal.NewFromInt(1500000),
		Currency:       "IDR",
		AvailableSeats: 10,
		CabinClass:     "economy",
		Origin: entity.Location{
			Airport:  "CGK",
			City:     "Jakarta",
			Time:     now.Add(10 * time.Hour),
			Timezone: time.UTC,
		},
		Destination: entity.Location{
			Airport:  "DPS",
			City:     "Denpasar",
			Time:     now.Add(14 * time.Hour),
			Timezone: time.UTC,
		},
	}

	// Return flight: departs 3 hours after outbound arrival (valid 90-1440 min layover)
	returnFlight := &entity.Flight{
		ID:             "ret1",
		FlightNumber:   "GA102",
		Provider:       "Garuda",
		AirlineCode:    "GA",
		Stops:          0,
		Price:          decimal.NewFromInt(1500000),
		Currency:       "IDR",
		AvailableSeats: 10,
		CabinClass:     "economy",
		Origin: entity.Location{
			Airport:  "DPS",
			City:     "Denpasar",
			Time:     now.Add(17 * time.Hour), // 3 hours after outbound arrival
			Timezone: time.UTC,
		},
		Destination: entity.Location{
			Airport:  "CGK",
			City:     "Jakarta",
			Time:     now.Add(21 * time.Hour),
			Timezone: time.UTC,
		},
	}

	itineraries := combiner.Combine([]*entity.Flight{outbound}, []*entity.Flight{returnFlight})

	if len(itineraries) != 1 {
		t.Errorf("Expected 1 itinerary, got %d", len(itineraries))
	}

	if itineraries[0].OutboundFlight.ID != "out1" {
		t.Errorf("OutboundFlight ID = %s, want out1", itineraries[0].OutboundFlight.ID)
	}

	if itineraries[0].ReturnFlight.ID != "ret1" {
		t.Errorf("ReturnFlight ID = %s, want ret1", itineraries[0].ReturnFlight.ID)
	}

	expectedPrice := decimal.NewFromInt(3000000)
	if !itineraries[0].TotalPrice.Equal(expectedPrice) {
		t.Errorf("TotalPrice = %s, want %s", itineraries[0].TotalPrice, expectedPrice)
	}
}

func TestRoundTripCombiner_Combine_TooShortLayover(t *testing.T) {
	combiner := NewRoundTripCombiner()

	now := time.Now().UTC()

	outbound := &entity.Flight{
		ID:             "out1",
		FlightNumber:   "GA101",
		Provider:       "Garuda",
		AirlineCode:    "GA",
		Stops:          0,
		Price:          decimal.NewFromInt(1500000),
		Currency:       "IDR",
		AvailableSeats: 10,
		CabinClass:     "economy",
		Origin: entity.Location{
			Airport:  "CGK",
			City:     "Jakarta",
			Time:     now.Add(10 * time.Hour),
			Timezone: time.UTC,
		},
		Destination: entity.Location{
			Airport:  "DPS",
			City:     "Denpasar",
			Time:     now.Add(14 * time.Hour),
			Timezone: time.UTC,
		},
	}

	// Return flight: departs only 30 minutes after outbound arrival (too short, min is 90 min)
	returnFlight := &entity.Flight{
		ID:             "ret1",
		FlightNumber:   "GA102",
		Provider:       "Garuda",
		AirlineCode:    "GA",
		Stops:          0,
		Price:          decimal.NewFromInt(1500000),
		Currency:       "IDR",
		AvailableSeats: 10,
		CabinClass:     "economy",
		Origin: entity.Location{
			Airport:  "DPS",
			City:     "Denpasar",
			Time:     now.Add(14*time.Hour + 30*time.Minute), // 30 min after outbound arrival
			Timezone: time.UTC,
		},
		Destination: entity.Location{
			Airport:  "CGK",
			City:     "Jakarta",
			Time:     now.Add(18 * time.Hour),
			Timezone: time.UTC,
		},
	}

	itineraries := combiner.Combine([]*entity.Flight{outbound}, []*entity.Flight{returnFlight})

	if len(itineraries) != 0 {
		t.Errorf("Expected 0 itineraries due to short layover, got %d", len(itineraries))
	}
}
