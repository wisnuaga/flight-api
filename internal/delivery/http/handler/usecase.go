package handler

import (
	"context"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

type FlightUsecase interface {
	Search(ctx context.Context, req *entity.SearchRequest) (*entity.SearchResult, error)
}
