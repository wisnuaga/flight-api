package port

import (
	"context"
	"time"
)

// Cache is a generic interface for Key-Value caching.
type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, bool, error)
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
