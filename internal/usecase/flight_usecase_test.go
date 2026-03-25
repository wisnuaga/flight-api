package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain"
	"github.com/wisnuaga/flight-api/internal/repository/provider"
)

type mockProvider struct {
	name    string
	flights []*domain.Flight
	err     error
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.Flight, error) {
	// Simulate small latency
	time.Sleep(10 * time.Millisecond)
	return m.flights, m.err
}

func ptr[T any](v T) *T {
	return &v
}

func TestFlightUsecase_Search(t *testing.T) {
	baseTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	mockFlights := []*domain.Flight{
		{ID: "F1", Provider: "Alpha", Price: 1500, Duration: 120 * time.Minute, Stops: 1, DepartureTime: baseTime, ArrivalTime: baseTime.Add(2 * time.Hour)},
		{ID: "F2", Provider: "Beta", Price: 1000, Duration: 150 * time.Minute, Stops: 0, DepartureTime: baseTime.Add(1 * time.Hour), ArrivalTime: baseTime.Add(3*time.Hour + 30*time.Minute)},
		{ID: "F3", Provider: "Alpha", Price: 2000, Duration: 90 * time.Minute, Stops: 0, DepartureTime: baseTime.Add(2 * time.Hour), ArrivalTime: baseTime.Add(3*time.Hour + 30*time.Minute)},
	}

	testCases := []struct {
		name          string
		providers     []provider.FlightProvider
		input         *domain.SearchRequest
		expectedLen   int
		expectedFirst string // ID of the expected first flight
	}{
		{
			name: "success - multiple providers, default sort asc price",
			providers: []provider.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*domain.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*domain.Flight{mockFlights[1]}, err: nil},
			},
			input:         &domain.SearchRequest{},
			expectedLen:   3,
			expectedFirst: "F2", // Price 1000 is lowest
		},
		{
			name: "success - filter by max price 1500",
			providers: []provider.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*domain.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*domain.Flight{mockFlights[1]}, err: nil},
			},
			input: &domain.SearchRequest{
				Filter: domain.SearchFilter{
					MaxPrice: ptr(float64(1500)),
				},
			},
			expectedLen:   2,
			expectedFirst: "F2",
		},
		{
			name: "success - sort by duration descending",
			providers: []provider.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*domain.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*domain.Flight{mockFlights[1]}, err: nil},
			},
			input: &domain.SearchRequest{
				Sort: domain.SearchSort{
					Field: domain.SortByDuration,
					Order: domain.SortDesc,
				},
			},
			expectedLen:   3,
			expectedFirst: "F2", // Duration 150 is the longest
		},
		{
			name: "success - one provider fails, results still returned",
			providers: []provider.FlightProvider{
				&mockProvider{name: "Alpha", flights: nil, err: errors.New("timeout")},
				&mockProvider{name: "Beta", flights: []*domain.Flight{mockFlights[1]}, err: nil},
			},
			input:         &domain.SearchRequest{},
			expectedLen:   1,
			expectedFirst: "F2",
		},
		{
			name: "success - filter by max stops",
			providers: []provider.FlightProvider{
				&mockProvider{name: "Alpha", flights: []*domain.Flight{mockFlights[0], mockFlights[2]}, err: nil},
				&mockProvider{name: "Beta", flights: []*domain.Flight{mockFlights[1]}, err: nil},
			},
			input: &domain.SearchRequest{
				Filter: domain.SearchFilter{
					MaxStops: ptr(0),
				},
			},
			expectedLen:   2,
			expectedFirst: "F2", // Filtered out F1 which has 1 stop
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			usecase := &FlightUsecaseImpl{providers: tc.providers}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			got, err := usecase.Search(ctx, tc.input)
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
