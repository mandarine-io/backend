package db_cacher

import (
	"context"
	"fmt"
	"github.com/go-gorm/caches/v4"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type dbCacher struct {
	manager cache.Manager
}

func NewDbCacher(manager cache.Manager) caches.Cacher {
	return &dbCacher{manager: manager}
}

func (c *dbCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	log.Debug().Msgf("get from DB cache %s", key)

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
	log.Debug().Msgf("store in DB cache %s", key)
	return c.manager.Set(ctx, key, *val)
}

func (c *dbCacher) Invalidate(ctx context.Context) error {
	log.Debug().Msg("invalidate DB cache")
	return c.manager.Invalidate(ctx, fmt.Sprintf("%s*", caches.IdentifierPrefix))
}
