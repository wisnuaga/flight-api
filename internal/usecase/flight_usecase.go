package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain"
	"github.com/wisnuaga/flight-api/internal/repository/provider"
)

type searchResult struct {
	flights []*domain.Flight
	err     error
}

type FlightUsecaseImpl struct {
	providers []provider.FlightProvider
}

func NewFlightUsecase(registry *provider.Registry) *FlightUsecaseImpl {
	return &FlightUsecaseImpl{providers: registry.GetProviders()}
}

func (u *FlightUsecaseImpl) Search(ctx context.Context, req *domain.SearchRequest) (*domain.SearchResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	resultCh := make(chan *searchResult, len(u.providers))

	for _, p := range u.providers {
		wg.Add(1)

		go func(p provider.FlightProvider) {
			defer wg.Done()

			flights, err := p.Search(ctx, req)

			resultCh <- &searchResult{
				flights: flights,
				err:     err,
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var allFlights []*domain.Flight
	success := 0
	failed := 0

	for res := range resultCh {
		if res.err != nil {
			failed++
			continue
		}

		success++
		allFlights = append(allFlights, res.flights...)
	}

	return u.buildSearchResult(allFlights, req, success, failed), nil
}

func (u *FlightUsecaseImpl) buildSearchResult(flights []*domain.Flight, req *domain.SearchRequest, success, failed int) *domain.SearchResult {
	flights = u.deduplicateFlights(flights)

	predicates := BuildFilterPredicates(&req.Filter)
	flights = ApplyFilters(flights, predicates)

	ApplySorting(flights, req.Sort)

	return &domain.SearchResult{
		Flights: flights,
		Meta: &domain.SearchMeta{
			TotalFlights: len(flights),
			Providers:    len(u.providers),
			SuccessCount: success,
			FailedCount:  failed,
		},
	}
}

func (u *FlightUsecaseImpl) deduplicateFlights(flights []*domain.Flight) []*domain.Flight {
	bestFlights := make(map[string]*domain.Flight)

	for _, f := range flights {
		key := f.FlightNumber + f.DepartureTime.String()

		existing, ok := bestFlights[key]
		// Retain the flight if it's new, or if its price is strictly cheaper than the previously seen one
		if !ok || f.Price < existing.Price {
			bestFlights[key] = f
		}
	}

	var result []*domain.Flight
	for _, f := range bestFlights {
		result = append(result, f)
	}

	return result
}
