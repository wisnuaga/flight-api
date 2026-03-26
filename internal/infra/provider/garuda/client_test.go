package garuda_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/provider/garuda"
)

func TestClient_Name(t *testing.T) {
	client := garuda.NewClient("../../../../../tests/garuda_ok.json")
	if got := client.Name(); got != "Garuda" {
		t.Errorf("Name() = %q, want %q", got, "Garuda")
	}
}

func ptrString(s string) *string {
	return &s
}

func TestClient_Search(t *testing.T) {
	baseReq := &entity.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
		Passengers:    1,
		Filter: entity.SearchFilter{
			CabinClass: ptrString("economy"),
		},
	}

	testCases := []struct {
		name             string
		mockPath         string
		ctx              func() (context.Context, context.CancelFunc)
		input            *entity.SearchRequest
		expectedLen      int
		expectedErr      error
		checkFirstFlight bool
		expectedFirst    *entity.Flight
	}{
		{
			name:     "success - returns mapped flights from valid mock file",
			mockPath: "../../../../../tests/garuda_ok.json",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:            baseReq,
			expectedLen:      3,
			expectedErr:      nil,
			checkFirstFlight: true,
			expectedFirst: &entity.Flight{
				ID:             "GA400",
				Provider:       "Garuda",
				FlightNumber:   "GA400",
				Origin:         "CGK",
				Destination:    "DPS",
				DepartureTime:  mustParseTime("2025-12-15T06:00:00+07:00"),
				ArrivalTime:    mustParseTime("2025-12-15T08:50:00+08:00"),
				Duration:       110 * time.Minute,
				Price:          1250000,
				Currency:       "IDR",
				CabinClass:     "economy",
				AvailableSeats: 28,
			},
		},
		{
			name:     "error - mock file does not exist",
			mockPath: "../../../../../tests/nonexistent.json",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:       baseReq,
			expectedLen: 0,
			expectedErr: errors.New("file not found"),
		},
		{
			name:     "error - context cancelled before response",
			mockPath: "../../../../../tests/garuda_ok.json",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Millisecond)
			},
			input:       baseReq,
			expectedLen: 0,
			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := tc.ctx()
			defer cancel()

			client := garuda.NewClient(tc.mockPath)
			got, err := client.Search(ctx, tc.input)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Search() error = nil, want non-nil (%v)", tc.expectedErr)
				}
				if errors.Is(tc.expectedErr, context.DeadlineExceeded) ||
					errors.Is(tc.expectedErr, context.Canceled) {
					if !errors.Is(err, tc.expectedErr) {
						t.Errorf("Search() error = %v, want %v", err, tc.expectedErr)
					}
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
