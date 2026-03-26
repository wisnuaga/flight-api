package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/port"
)

type searchResult struct {
	flights  []*entity.Flight
	err      error
	cacheHit bool
}

type FlightUsecaseImpl struct {
	providers []port.FlightProvider
	cache     port.Cache[[]*entity.Flight]
	filterCmd port.FlightFilterCommand
	sortCmd   port.FlightSortCommand
}

func NewFlightUsecase(
	providers []port.FlightProvider,
	c port.Cache[[]*entity.Flight],
	filterCmd port.FlightFilterCommand,
	sortCmd port.FlightSortCommand,
) *FlightUsecaseImpl {
	return &FlightUsecaseImpl{
		providers: providers,
		cache:     c,
		filterCmd: filterCmd,
		sortCmd:   sortCmd,
	}
}

func GenerateFlightSearchKey(provider string, req *entity.SearchRequest) string {
	raw := fmt.Sprintf("%s_%s_%s_%s_%d",
		provider,
		req.Origin,
		req.Destination,
		req.DepartureDate.Format("2006-01-02"),
		req.Passengers)
	hash := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("flight:search:v1:%x", hash)
}

func (u *FlightUsecaseImpl) Search(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	resultCh := make(chan *searchResult, len(u.providers))

	for _, p := range u.providers {
		wg.Add(1)

		go func(p port.FlightProvider) {
			defer wg.Done()

			cacheKey := GenerateFlightSearchKey(p.Name(), req)

			// Provider isolation caching check BEFORE executing network request
			if cachedFlights, ok, err := u.cache.Get(ctx, cacheKey); err == nil && ok {
				select {
				case <-ctx.Done():
				case resultCh <- &searchResult{flights: cachedFlights, err: nil, cacheHit: true}:
				}
				return
			}

			// Exponential backoff retry (up to 3 times) for flaky upstream APIs
			flights, err := u.retryProviderSearch(ctx, p, req)
			if err == nil {
				_ = u.cache.Set(ctx, cacheKey, flights, 5*time.Minute)
			}

			select {
			case <-ctx.Done():
			case resultCh <- &searchResult{flights: flights, err: err, cacheHit: false}:
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	bestFlights := make(map[string]*entity.Flight)
	success := 0
	failed := 0
	cacheHits := 0

	// O(N) Single-pass deduplication directly from the channel stream
	for res := range resultCh {
		if res.err != nil {
			failed++
			continue
		}
		if res.cacheHit {
			cacheHits++
		}

		success++
		for _, f := range res.flights {
			// Code-share deduplication heuristic: origin + dest + 5-min rounded departure/arrival
			dep := f.DepartureTime.Truncate(5 * time.Minute).Unix()
			arr := f.ArrivalTime.Truncate(5 * time.Minute).Unix()
			key := fmt.Sprintf("%s_%s_%d_%d", f.Origin.Airport, f.Destination.Airport, dep, arr)

			// Simple de-duplication: keep the cheaper flight
			if existing, ok := bestFlights[key]; !ok || f.Price.LessThan(existing.Price) {
				bestFlights[key] = f
			}
		}
	}

	var deduplicated []*entity.Flight
	for _, f := range bestFlights {
		deduplicated = append(deduplicated, f)
	}

	result := u.buildSearchResult(deduplicated, req, success, failed)
	result.Meta.CacheHit = cacheHits > 0
	result.Meta.SearchTimeMs = int(time.Since(start).Milliseconds())

	return result, nil
}

func (u *FlightUsecaseImpl) buildSearchResult(flights []*entity.Flight, req *entity.SearchRequest, success, failed int) *entity.SearchResult {
	flights = u.filterCmd.Execute(flights, &req.Filter)
	u.sortCmd.Execute(flights, req.Sort)

	return &entity.SearchResult{
		Flights: flights,
		Meta: &entity.SearchMeta{
			TotalFlights: len(flights),
			Providers:    len(u.providers),
			SuccessCount: success,
			FailedCount:  failed,
		},
	}
}

// retryProviderSearch handles retries with exponential backoff
func (u *FlightUsecaseImpl) retryProviderSearch(ctx context.Context, p port.FlightProvider, req *entity.SearchRequest) ([]*entity.Flight, error) {
	var flights []*entity.Flight
	var err error

	const maxAttempts = 3
	baseBackoff := 50 * time.Millisecond

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		flights, err = p.Search(ctx, req)
		if err == nil {
			return flights, nil
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		traceID, _ := ctx.Value("trace_id").(string)
		slog.Warn("provider search failed, retrying",
			slog.String("trace_id", traceID),
			slog.String("provider", p.Name()),
			slog.Int("attempt", attempt),
			slog.String("error", err.Error()),
		)

		// exponential backoff
		select {
		case <-time.After(time.Duration(1<<uint(attempt-1)) * baseBackoff):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	traceID, _ := ctx.Value("trace_id").(string)
	slog.Error("provider search failed completely",
		slog.String("trace_id", traceID),
		slog.String("provider", p.Name()),
		slog.String("error", err.Error()),
	)

	return nil, err
}
