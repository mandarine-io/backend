package db_cacher

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-gorm/caches/v4"
	manager2 "github.com/mandarine-io/Backend/pkg/storage/cache/manager"
	"github.com/rs/zerolog/log"
)

type dbCacher struct {
	manager manager2.CacheManager
}

func NewDbCacher(manager manager2.CacheManager) caches.Cacher {
	return &dbCacher{manager: manager}
}

func (c *dbCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
	log.Debug().Msgf("get from DB cache %s", key)

	err := c.manager.Get(ctx, key, q)
	if errors.Is(err, manager2.ErrCacheEntryNotFound) {
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
