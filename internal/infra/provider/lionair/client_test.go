package lionair_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/provider/lionair"
	"github.com/wisnuaga/flight-api/internal/test_helper"
)

func TestClient_Name(t *testing.T) {
	client := lionair.NewClient(test_helper.GetTestDataPath("lion_air_search_response.json"))
	if got := client.Name(); got != "Lion Air" {
		t.Errorf("Name() = %q, want %q", got, "Lion Air")
	}
}

func TestClient_Search(t *testing.T) {
	baseReq := &entity.SearchRequest{}

	testCases := []struct {
		name             string
		mockPath         string
		ctx              func() (context.Context, context.CancelFunc)
		input            *entity.SearchRequest
		expectedLen      int
		expectErrStatus  bool
		checkFirstFlight bool
		expectedFirst    *entity.Flight
	}{
		{
			name:     "success - returns mapped flights from valid mock file",
			mockPath: test_helper.GetTestDataPath("lion_air_search_response.json"),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:            baseReq,
			expectedLen:      3,
			expectErrStatus:  false,
			checkFirstFlight: true,
			expectedFirst: &entity.Flight{
				ID:           "JT740_LionAir",
				Airline:      entity.AirlineLionAir,
				FlightNumber: "JT740",
				AirlineCode:  "JT",
				Origin: entity.Location{
					Airport: "CGK",
					Time:    mustParseTime("2025-12-15T05:30:00+07:00").UTC(),
				},
				Destination: entity.Location{
					Airport: "DPS",
					Time:    mustParseTime("2025-12-15T08:15:00+08:00").UTC(),
				},
				Duration:       105 * time.Minute,
				Price:          decimal.NewFromInt(950000),
				Currency:       "IDR",
				CabinClass:     "ECONOMY",
				AvailableSeats: 45,
			},
		},
		{
			name:     "error - mock file does not exist",
			mockPath: test_helper.GetTestDataPath("nonexistent.json"),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:           baseReq,
			expectedLen:     0,
			expectErrStatus: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := tc.ctx()
			defer cancel()

			client := lionair.NewClient(tc.mockPath)
			got, err := client.Search(ctx, tc.input)

			if tc.expectErrStatus {
				if err == nil {
					t.Fatalf("Search() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Search() unexpected error: %v", err)
			}

			if len(got) != tc.expectedLen {
				t.Errorf("Search() returned %d flights, want %d", len(got), tc.expectedLen)
			}

			if tc.checkFirstFlight && len(got) > 0 {
				assertFlight(t, got[0], tc.expectedFirst)
			}
		})
	}
}

func assertFlight(t *testing.T, got, want *entity.Flight) {
	t.Helper()
	fields := []struct {
		name     string
		got, exp interface{}
	}{
		{"ID", got.ID, want.ID},
		{"Airline", got.Airline, want.Airline},
		{"FlightNumber", got.FlightNumber, want.FlightNumber},
		{"Origin.Airport", got.Origin.Airport, want.Origin.Airport},
		{"Destination.Airport", got.Destination.Airport, want.Destination.Airport},
		{"Origin.Time", got.Origin.Time.UTC(), want.Origin.Time.UTC()},
		{"Destination.Time", got.Destination.Time.UTC(), want.Destination.Time.UTC()},
		{"Duration", got.Duration, want.Duration},
		{"Currency", got.Currency, want.Currency},
		{"CabinClass", got.CabinClass, want.CabinClass},
		{"AvailableSeats", got.AvailableSeats, want.AvailableSeats},
	}
	for _, f := range fields {
		if f.got != f.exp {
			t.Errorf("Flight.%s = %v, want %v", f.name, f.got, f.exp)
		}
	}
	// Compare price separately using decimal comparison
	if !got.Price.Equal(want.Price) {
		t.Errorf("Flight.Price = %v, want %v", got.Price, want.Price)
	}
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
