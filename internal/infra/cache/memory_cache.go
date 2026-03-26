package cache

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

type cacheEntry struct {
	flights   []*entity.Flight
	expiresAt time.Time
}

// MemoryCache is a thread-safe in-memory TTL cache implementing port.FlightCache.
type MemoryCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	ttl     time.Duration
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		entries: make(map[string]cacheEntry),
		ttl:     ttl,
	}
}

func (c *MemoryCache) generateKey(providerName string, req *entity.SearchRequest) string {
	raw := fmt.Sprintf("%s_%s_%s_%s_%d",
		providerName,
		req.Origin,
		req.Destination,
		req.DepartureDate.Format("2006-01-02"),
		req.Passengers)

	hash := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", hash)
}

func (c *MemoryCache) Get(providerName string, req *entity.SearchRequest) ([]*entity.Flight, bool) {
	key := c.generateKey(providerName, req)

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.flights, true
}

func (c *MemoryCache) Set(providerName string, req *entity.SearchRequest, flights []*entity.Flight) {
	key := c.generateKey(providerName, req)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		flights:   flights,
		expiresAt: time.Now().Add(c.ttl),
	}
}
