package db_cacher

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-gorm/caches/v4"
	manager2 "mandarine/pkg/storage/cache/manager"
)

type dbCacher struct {
	manager manager2.CacheManager
}

func NewDbCacher(manager manager2.CacheManager) caches.Cacher {
	return &dbCacher{manager: manager}
}

func (c *dbCacher) Get(ctx context.Context, key string, q *caches.Query[any]) (*caches.Query[any], error) {
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
	return c.manager.Set(ctx, key, *val)
}

func (c *dbCacher) Invalidate(ctx context.Context) error {
	return c.manager.Invalidate(ctx, fmt.Sprintf("%s*", caches.IdentifierPrefix))
}
