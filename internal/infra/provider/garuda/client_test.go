package garuda_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/provider/garuda"
	"github.com/wisnuaga/flight-api/internal/test_helper"
)

func TestClient_Name(t *testing.T) {
	client := garuda.NewClient(test_helper.GetTestDataPath("garuda_search_response.json"))
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
			mockPath: test_helper.GetTestDataPath("garuda_search_response.json"),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:            baseReq,
			expectedLen:      3,
			expectedErr:      nil,
			checkFirstFlight: true,
			expectedFirst: &entity.Flight{
				ID:           "GA400",
				Provider:     "Garuda",
				FlightNumber: "GA400",
				Origin: entity.Location{
					Airport:  "CGK",
					Time:     mustParseTime("2025-12-15T06:00:00+07:00").UTC(),
					Timezone: time.UTC, // Will be set to Asia/Jakarta by mapper
				},
				Destination: entity.Location{
					Airport:  "DPS",
					Time:     mustParseTime("2025-12-15T08:50:00+08:00").UTC(),
					Timezone: time.UTC, // Will be set to Asia/Jakarta by mapper
				},
				Duration:       110 * time.Minute,
				Price:          decimal.NewFromInt(1250000),
				Currency:       "IDR",
				CabinClass:     "economy",
				AvailableSeats: 28,
			},
		},
		{
			name:     "error - mock file does not exist",
			mockPath: test_helper.GetTestDataPath("nonexistent.json"),
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			input:       baseReq,
			expectedLen: 0,
			expectedErr: errors.New("file not found"),
		},
		{
			name:     "error - context cancelled before response",
			mockPath: test_helper.GetTestDataPath("garuda_search_response.json"),
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
