package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/port"
)

type mockProvider struct {
	name    string
	flights []*entity.Flight
	err     error
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	// Simulate small latency
	time.Sleep(10 * time.Millisecond)
	return m.flights, m.err
}

type mockCache struct{}

func (m *mockCache) Get(provider string, req *entity.SearchRequest) ([]*entity.Flight, bool) {
	return nil, false
}

func (m *mockCache) Set(provider string, req *entity.SearchRequest, flights []*entity.Flight) {}

func ptr[T any](v T) *T {
	return &v
}

func TestFlightUsecase_Search(t *testing.T) {
	baseTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	mockFlights := []*entity.Flight{
		{ID: "F1", Provider: "Alpha", FlightNumber: "AL1", Price: 1500, Duration: 120 * time.Minute, Stops: 1, DepartureTime: baseTime, ArrivalTime: baseTime.Add(2 * time.Hour)},
		{ID: "F2", Provider: "Beta", FlightNumber: "BE2", Price: 1000, Duration: 150 * time.Minute, Stops: 0, DepartureTime: baseTime.Add(1 * time.Hour), ArrivalTime: baseTime.Add(3*time.Hour + 30*time.Minute)},
		{ID: "F3", Provider: "Alpha", FlightNumber: "AL3", Price: 2000, Duration: 90 * time.Minute, Stops: 0, DepartureTime: baseTime.Add(2 * time.Hour), ArrivalTime: baseTime.Add(3*time.Hour + 30*time.Minute)},
		// Duplicate of F1 but cheaper
		{ID: "F4", Provider: "Gamma", FlightNumber: "AL1", Price: 1200, Duration: 120 * time.Minute, Stops: 1, DepartureTime: baseTime, ArrivalTime: baseTime.Add(2 * time.Hour)},
	}

	testCases := []struct {
		name          string
		providers     []port.FlightProvider
		input         *entity.SearchRequest
		expectedLen   int
		expectedFirst string // ID of the expected first flight
	}{
		{
			name: "success - multiple providers, default sort asc price",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*entity.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*entity.Flight{mockFlights[1]}, err: nil},
			},
			input:         &entity.SearchRequest{},
			expectedLen:   3,
			expectedFirst: "F2", // Price 1000 is lowest
		},
		{
			name: "success - filter by max price 1500",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*entity.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*entity.Flight{mockFlights[1]}, err: nil},
			},
			input: &entity.SearchRequest{
				Filter: entity.SearchFilter{
					MaxPrice: ptr(float64(1500)),
				},
			},
			expectedLen:   2,
			expectedFirst: "F2",
		},
		{
			name: "success - sort by duration descending",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*entity.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*entity.Flight{mockFlights[1]}, err: nil},
			},
			input: &entity.SearchRequest{
				Sort: entity.SearchSort{
					Field: entity.SortByDuration,
					Order: entity.SortDesc,
				},
			},
			expectedLen:   3,
			expectedFirst: "F2", // Duration 150 is the longest
		},
		{
			name: "success - one provider fails, results still returned",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: nil, err: errors.New("timeout")},
				&mockProvider{name: "Beta", flights: []*entity.Flight{mockFlights[1]}, err: nil},
			},
			input:         &entity.SearchRequest{},
			expectedLen:   1,
			expectedFirst: "F2",
		},
		{
			name: "success - filter by max stops",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*entity.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*entity.Flight{mockFlights[1]}, err: nil},
			},
			input: &entity.SearchRequest{
				Filter: entity.SearchFilter{
					MaxStops: ptr(0),
				},
			},
			expectedLen:   2,
			expectedFirst: "F2", // Filtered out F1 which has 1 stop
		},
		{
			name: "success - deduplication keeps cheaper flight",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*entity.Flight{mockFlights[0]}, err: nil}, // F1 (1500)
				&mockProvider{name: "Gamma", flights: []*entity.Flight{mockFlights[3]}, err: nil}, // F4 (1200) - same route/time
			},
			input:         &entity.SearchRequest{},
			expectedLen:   1,    // Should deduplicate F1 and F4 into 1
			expectedFirst: "F4", // Should keep the cheaper one
		},
		{
			name: "success - best value ranking",
			providers: []port.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*entity.Flight{mockFlights[2]}, err: nil}, // F3 (price 2000, dur 90)
				&mockProvider{name: "Beta", flights: []*entity.Flight{mockFlights[1]}, err: nil},  // F2 (price 1000, dur 150)
				&mockProvider{name: "Gamma", flights: []*entity.Flight{mockFlights[3]}, err: nil}, // F4 (price 1200, dur 120)
			},
			input: &entity.SearchRequest{
				Sort: entity.SearchSort{
					Field:          entity.SortByBestValue,
					Order:          entity.SortAsc,
					PriceWeight:    1.0,
					DurationWeight: 1.0,
				},
			},
			expectedLen:   3,
			expectedFirst: "F4",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			uc := &FlightUsecaseImpl{
				providers: tc.providers,
				cache:     &mockCache{},
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			got, err := uc.Search(ctx, tc.input)
			if err != nil {
				t.Fatalf("Search() unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("Search() returned nil result")
			}

			if len(got.Flights) != tc.expectedLen {
				t.Errorf("Search() returned %d flights, want %d", len(got.Flights), tc.expectedLen)
			}

			if tc.expectedLen > 0 && got.Flights[0].ID != tc.expectedFirst {
				t.Errorf("Search() first flight ID = %s, want %s", got.Flights[0].ID, tc.expectedFirst)
			}
		})
	}
}
