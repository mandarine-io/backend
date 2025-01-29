package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCacheEntryNotFound = errors.New("cache entry not found")
)

type Manager interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any) error
	SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Invalidate(ctx context.Context, keyRegex string) error
}
