package dbcacher

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-gorm/caches/v4"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/rs/zerolog"
)

type dbCacher struct {
	manager cache.Manager
	logger  zerolog.Logger
}

type Option func(*dbCacher) error

func WithLogger(logger zerolog.Logger) Option {
	return func(c *dbCacher) error {
		c.logger = logger
		return nil
	}
}

func NewDbCacher(manager cache.Manager, opts ...Option) (caches.Cacher, error) {
	c := &dbCacher{
		manager: manager,
		logger:  zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return c, nil
}

func (c *dbCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	c.logger.Debug().Msgf("get from DB cache %s", key)

	err := c.manager.Get(ctx, key, q)
	if errors.Is(err, cache.ErrCacheEntryNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (c *dbCacher) Store(ctx context.Context, key string, val *caches.Query[any]) error {
	c.logger.Debug().Msgf("store in DB cache %s", key)
	return c.manager.Set(ctx, key, *val)
}

func (c *dbCacher) Invalidate(ctx context.Context) error {
	c.logger.Debug().Msg("invalidate DB cache")
	return c.manager.Invalidate(ctx, fmt.Sprintf("%s*", caches.IdentifierPrefix))
}
