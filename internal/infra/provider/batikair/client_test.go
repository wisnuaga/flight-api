package batikair_test

import (
	"context"
	"testing"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/provider/batikair"
)

func TestClient_Name(t *testing.T) {
	client := batikair.NewClient("../../../../../tests/factory/batik_air_ok.json")
	if got := client.Name(); got != "Batik Air" {
		t.Errorf("Name() = %q, want %q", got, "Batik Air")
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
			mockPath: "../../../../../tests/factory/batik_air_ok.json",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:            baseReq,
			expectedLen:      3,
			expectErrStatus:  false,
			checkFirstFlight: true,
			expectedFirst: &entity.Flight{
				ID:             "ID6514",
				Provider:       "Batik Air",
				FlightNumber:   "ID6514",
				Origin:         "CGK",
				Destination:    "DPS",
				DepartureTime:  mustParseTime("2025-12-15T07:15:00+07:00"),
				ArrivalTime:    mustParseTime("2025-12-15T10:00:00+08:00"),
				Duration:       105 * time.Minute,
				Price:          1100000,
				Currency:       "IDR",
				CabinClass:     "Y",
				AvailableSeats: 32,
			},
		},
		{
			name:     "error - mock file does not exist",
			mockPath: "../../../../../tests/factory/nonexistent.json",
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

			client := batikair.NewClient(tc.mockPath)
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
		{"Provider", got.Provider, want.Provider},
		{"FlightNumber", got.FlightNumber, want.FlightNumber},
		{"Origin", got.Origin, want.Origin},
		{"Destination", got.Destination, want.Destination},
		{"DepartureTime", got.DepartureTime.UTC(), want.DepartureTime.UTC()},
		{"ArrivalTime", got.ArrivalTime.UTC(), want.ArrivalTime.UTC()},
		{"Duration", got.Duration, want.Duration},
		{"Price", got.Price, want.Price},
		{"Currency", got.Currency, want.Currency},
		{"CabinClass", got.CabinClass, want.CabinClass},
		{"AvailableSeats", got.AvailableSeats, want.AvailableSeats},
	}
	for _, f := range fields {
		if f.got != f.exp {
			t.Errorf("Flight.%s = %v, want %v", f.name, f.got, f.exp)
		}
	}
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
