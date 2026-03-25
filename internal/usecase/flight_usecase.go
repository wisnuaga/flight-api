package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain"
	"github.com/wisnuaga/flight-api/internal/repository/provider"
)

type FlightUsecase interface {
	Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.Flight, error)
}

type FlightUsecaseImpl struct {
	providers []provider.FlightProvider
}

func NewFlightUsecase(registry *provider.Registry) *FlightUsecaseImpl {
	return &FlightUsecaseImpl{providers: registry.GetProviders()}
}

func (u *FlightUsecaseImpl) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.Flight, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	resultCh := make(chan []*domain.Flight, len(u.providers))

	for _, p := range u.providers {
		wg.Add(1)

		go func(p provider.FlightProvider) {
			defer wg.Done()

			flights, err := p.Search(ctx, req)
			if err != nil {
				return
			}

			resultCh <- flights
		}(p)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var allFlights []*domain.Flight
	for flights := range resultCh {
		allFlights = append(allFlights, flights...)
	}

	return allFlights, nil
}
