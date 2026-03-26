package port

import (
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

// FlightCache is the interface for provider-level search result caching.
type FlightCache interface {
	Get(provider string, req *entity.SearchRequest) ([]*entity.Flight, bool)
	Set(provider string, req *entity.SearchRequest, flights []*entity.Flight)
}
