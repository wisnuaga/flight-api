package provider

import (
	"context"

	"github.com/wisnuaga/flight-api/internal/domain"
)

type FlightProvider interface {
	Name() string
	Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.Flight, error)
}
