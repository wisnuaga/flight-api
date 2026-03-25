package handler

import (
	"context"

	"github.com/wisnuaga/flight-api/internal/domain"
)

type FlightUsecase interface {
	Search(ctx context.Context, req *domain.SearchRequest) (*domain.SearchResult, error)
}
