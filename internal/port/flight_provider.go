package port

import (
	"context"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

// FlightProvider is the interface that all airline provider adapters must implement.
type FlightProvider interface {
	Name() string
	Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error)
}
