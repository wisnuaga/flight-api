package usecase

import (
	"context"
	"fmt"
	"log/slog"
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
			if err != nil {
				// Inject traceID if present inside context for observability correlation
				traceID, _ := ctx.Value("trace_id").(string)
				slog.Error("flight provider aggregated search failed",
					slog.String("trace_id", traceID),
					slog.String("provider", p.Name()),
					slog.String("error", err.Error()),
				)
			}

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

	bestFlights := make(map[string]*domain.Flight)
	success := 0
	failed := 0

	// O(N) Single-pass deduplication directly from the channel stream.
	// Eradicates the massive 'allFlights' intermediate memory slice entirely!
	for res := range resultCh {
		if res.err != nil {
			failed++
			continue
		}

		success++
		for _, f := range res.flights {
			// Code-share deductive grouping heuristic: origin + dest + 5-minute departure/arrival truncation
			dep := f.DepartureTime.Truncate(5 * time.Minute).Unix()
			arr := f.ArrivalTime.Truncate(5 * time.Minute).Unix()
			key := fmt.Sprintf("%s_%s_%d_%d", f.Origin, f.Destination, dep, arr)

			if existing, ok := bestFlights[key]; !ok || f.Price < existing.Price {
				bestFlights[key] = f
			}
		}
	}

	// Pack the tightly verified unique flights into a single exact-bound slice
	var deduplicated []*domain.Flight
	for _, f := range bestFlights {
		deduplicated = append(deduplicated, f)
	}

	return u.buildSearchResult(deduplicated, req, success, failed), nil
}

func (u *FlightUsecaseImpl) buildSearchResult(flights []*domain.Flight, req *domain.SearchRequest, success, failed int) *domain.SearchResult {
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
