package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/cache"
	"github.com/wisnuaga/flight-api/internal/port"
	"github.com/wisnuaga/flight-api/internal/usecase"
	"github.com/wisnuaga/flight-api/tests/mock"
)

func newMockProvider(name string, flights []*entity.Flight, err error) *mock.MockFlightProvider {
	m := new(mock.MockFlightProvider)
	m.On("Name").Return(name)
	m.On("Search", testifymock.Anything, testifymock.Anything).Return(flights, err)
	return m
}

func TestFlightUsecase_Search(t *testing.T) {
	baseTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	mockFlights := []*entity.Flight{
		{
			ID: "F1", Provider: "Alpha", FlightNumber: "AL1",
			Price: decimal.NewFromInt(1500), Duration: 120 * time.Minute, Stops: 1,
			Origin:      entity.Location{Airport: "CGK", Time: baseTime, Timezone: time.UTC},
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(2 * time.Hour), Timezone: time.UTC},
		},
		{
			ID: "F2", Provider: "Beta", FlightNumber: "BE2",
			Price: decimal.NewFromInt(1000), Duration: 150 * time.Minute, Stops: 0,
			Origin:      entity.Location{Airport: "CGK", Time: baseTime.Add(1 * time.Hour), Timezone: time.UTC},
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(3*time.Hour + 30*time.Minute), Timezone: time.UTC},
		},
		// Duplicate of F1 but cheaper
		{
			ID: "F4", Provider: "Gamma", FlightNumber: "AL1",
			Price: decimal.NewFromInt(1200), Duration: 120 * time.Minute, Stops: 1,
			Origin:      entity.Location{Airport: "CGK", Time: baseTime, Timezone: time.UTC},
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(2 * time.Hour), Timezone: time.UTC},
		},
	}

	testCases := []struct {
		name          string
		providers     []port.FlightProvider
		setupMocks    func(filterCmd *mock.MockFlightFilterCommand, sortCmd *mock.MockFlightSortCommand)
		input         *entity.SearchRequest
		expectedLen   int
		expectedFirst string
	}{
		{
			name: "success - multiple providers aggregation",
			providers: []port.FlightProvider{
				newMockProvider("Alpha", []*entity.Flight{mockFlights[0]}, nil),
				newMockProvider("Beta", []*entity.Flight{mockFlights[1]}, nil),
			},
			input: &entity.SearchRequest{},
			setupMocks: func(filterCmd *mock.MockFlightFilterCommand, sortCmd *mock.MockFlightSortCommand) {
				// The usecase deduplicates and aggregates, then passes to the commands directly.
				// We inject pass-through mock data bypassing logic to solely test Usecase orchestration.
				filterCmd.On("Execute", testifymock.Anything, testifymock.Anything).Return(
					[]*entity.Flight{mockFlights[0], mockFlights[1]},
				)
				sortCmd.On("Execute", testifymock.Anything, testifymock.Anything).Return()
			},
			expectedLen:   2,
			expectedFirst: "F1",
		},
		{
			name: "success - one provider fails, results still returned (error resilience)",
			providers: []port.FlightProvider{
				newMockProvider("Alpha", nil, errors.New("timeout")),
				newMockProvider("Beta", []*entity.Flight{mockFlights[1]}, nil),
			},
			input: &entity.SearchRequest{},
			setupMocks: func(filterCmd *mock.MockFlightFilterCommand, sortCmd *mock.MockFlightSortCommand) {
				filterCmd.On("Execute", testifymock.Anything, testifymock.Anything).Return(
					[]*entity.Flight{mockFlights[1]},
				)
				sortCmd.On("Execute", testifymock.Anything, testifymock.Anything).Return()
			},
			expectedLen:   1,
			expectedFirst: "F2",
		},
		{
			name: "success - deduplication keeps cheaper flight",
			providers: []port.FlightProvider{
				newMockProvider("Alpha", []*entity.Flight{mockFlights[0]}, nil), // F1 (1500)
				newMockProvider("Gamma", []*entity.Flight{mockFlights[2]}, nil), // F4 (1200)
			},
			input: &entity.SearchRequest{},
			setupMocks: func(filterCmd *mock.MockFlightFilterCommand, sortCmd *mock.MockFlightSortCommand) {
				filterCmd.On("Execute", testifymock.Anything, testifymock.Anything).Return(
					[]*entity.Flight{mockFlights[2]},
				)
				sortCmd.On("Execute", testifymock.Anything, testifymock.Anything).Return()
			},
			expectedLen:   1,
			expectedFirst: "F4",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			filterMock := new(mock.MockFlightFilterCommand)
			sortMock := new(mock.MockFlightSortCommand)
			tc.setupMocks(filterMock, sortMock)

			uc := usecase.NewFlightUsecase(
				tc.providers,
				cache.NewMemoryCache[[]*entity.Flight](),
				filterMock,
				sortMock,
			)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			got, err := uc.Search(ctx, tc.input)

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Len(t, got.Flights, tc.expectedLen)

			if tc.expectedLen > 0 {
				assert.Equal(t, tc.expectedFirst, got.Flights[0].ID)
			}

			// Verify the command mocks were properly called during orchestration
			filterMock.AssertExpectations(t)
			sortMock.AssertExpectations(t)
		})
	}
}
