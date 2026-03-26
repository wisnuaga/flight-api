package cache

import (
	"context"
	"sync"
	"time"
)

type memoryEntry[T any] struct {
	value     T
	expiresAt time.Time
}

// MemoryCache is a thread-safe generic in-memory cache.
type MemoryCache[T any] struct {
	mu      sync.RWMutex
	entries map[string]memoryEntry[T]
}

// NewMemoryCache creates a new in-memory cache instance.
func NewMemoryCache[T any]() *MemoryCache[T] {
	return &MemoryCache[T]{
		entries: make(map[string]memoryEntry[T]),
	}
}

func (c *MemoryCache[T]) Get(ctx context.Context, key string) (T, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	var zero T
	if !exists {
		return zero, false, nil
	}

	if time.Now().After(entry.expiresAt) {
		return zero, false, nil
	}

	return entry.value, true, nil
}

func (c *MemoryCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = memoryEntry[T]{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *MemoryCache[T]) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
	return nil
}
