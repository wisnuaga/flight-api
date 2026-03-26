package cache_test

import (
	"testing"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/infra/cache"
)

func TestMemoryCache_GetAndSet(t *testing.T) {
	c := cache.NewMemoryCache(50 * time.Millisecond)

	req := &entity.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: time.Now(),
		Passengers:    1,
	}

	result := []*entity.Flight{
		{ID: "F1"},
		{ID: "F2"},
	}

	// Test initial empty
	if _, ok := c.Get("Alpha", req); ok {
		t.Error("Expected cache to be empty")
	}

	// Test Set and Hit
	c.Set("Alpha", req, result)
	if cached, ok := c.Get("Alpha", req); !ok {
		t.Error("Expected cache hit")
	} else if len(cached) != 2 {
		t.Errorf("Expected 2 flights, got %d", len(cached))
	}

	// Test TTL Expiration
	time.Sleep(60 * time.Millisecond)
	if _, ok := c.Get("Alpha", req); ok {
		t.Error("Expected cache to expire after TTL")
	}
}
