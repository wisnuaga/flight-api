package usecase

import (
	"context"
	"sort"
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

	sort.Slice(allFlights, func(i, j int) bool {
		return allFlights[i].Price < allFlights[j].Price
	})

	return u.buildSearchResult(allFlights, req), nil
}

func (u *FlightUsecaseImpl) buildSearchResult(flights []*domain.Flight, req *domain.SearchRequest) *domain.SearchResult {
	flights = u.deduplicateFlights(flights)
	flights = u.filterFlights(flights, req)

	return &domain.SearchResult{
		Flights: flights,
		Meta: &domain.SearchMeta{
			TotalFlights: len(flights),
			Providers:    len(u.providers),
			SuccessCount: len(flights),
			FailedCount:  len(u.providers) - len(flights),
		},
	}
}

func (u *FlightUsecaseImpl) filterFlights(flights []*domain.Flight, req *domain.SearchRequest) []*domain.Flight {
	var result []*domain.Flight

	for _, f := range flights {
		if req.CabinClass != "" && f.CabinClass != req.CabinClass {
			continue
		}

		result = append(result, f)
	}

	return result
}

func (u *FlightUsecaseImpl) deduplicateFlights(flights []*domain.Flight) []*domain.Flight {
	seen := make(map[string]bool)
	var result []*domain.Flight

	for _, f := range flights {
		key := f.FlightNumber + f.DepartureTime.String()

		if seen[key] {
			continue
		}

		seen[key] = true
		result = append(result, f)
	}

	return result
}
