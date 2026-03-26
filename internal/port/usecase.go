package port

import (
	"context"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

// FlightUsecase encapsulates the primary flight search use cases.
type FlightUsecase interface {
	Search(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error)
}
