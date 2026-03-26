package airasia_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/provider/airasia"
	"github.com/wisnuaga/flight-api/internal/test_helper"
)

func TestClient_Name(t *testing.T) {
	client := airasia.NewClient(test_helper.GetTestDataPath("airasia_search_response.json"))
	if got := client.Name(); got != "AirAsia" {
		t.Errorf("Name() = %q, want %q", got, "AirAsia")
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
			mockPath: test_helper.GetTestDataPath("airasia_search_response.json"),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:            baseReq,
			expectedLen:      4,
			expectErrStatus:  false,
			checkFirstFlight: true,
			expectedFirst: &entity.Flight{
				ID:             "QZ520",
				Provider:       "AirAsia",
				FlightNumber:   "QZ520",
				Origin:         entity.Location{Airport: "CGK"},
				Destination:    entity.Location{Airport: "DPS"},
				DepartureTime:  mustParseTime("2025-12-15T04:45:00+07:00"),
				ArrivalTime:    mustParseTime("2025-12-15T07:25:00+08:00"),
				Duration:       100 * time.Minute,
				Price:          decimal.NewFromInt(650000),
				Currency:       "IDR",
				CabinClass:     "economy",
				AvailableSeats: 67,
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

			client := airasia.NewClient(tc.mockPath)

			var got []*entity.Flight
			var err error

			// AirAsia has a 10% simulated random timeout flake error. We retry to bypass it in unit tests.
			for i := 0; i < 20; i++ {
				got, err = client.Search(ctx, tc.input)

				if tc.expectErrStatus {
					// We only care if we get a real expected structural error, like file missing (not random simulation)
					if err != nil && !strings.Contains(err.Error(), "upstream") {
						return // test passed
					}
					// If it simulated, retry to get the real error
					time.Sleep(10 * time.Millisecond)
					continue
				}

				if err == nil {
					break // Passed cleanly
				}
			}

			if tc.expectErrStatus {
				t.Fatalf("Search() expected error, got nil")
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
		{"Provider", got.Provider, want.Provider},
		{"FlightNumber", got.FlightNumber, want.FlightNumber},
		{"Origin.Airport", got.Origin.Airport, want.Origin.Airport},
		{"Destination.Airport", got.Destination.Airport, want.Destination.Airport},
		{"DepartureTime", got.DepartureTime.UTC(), want.DepartureTime.UTC()},
		{"ArrivalTime", got.ArrivalTime.UTC(), want.ArrivalTime.UTC()},
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
