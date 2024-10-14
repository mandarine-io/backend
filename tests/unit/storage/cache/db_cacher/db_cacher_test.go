package db_cacher_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-gorm/caches/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mandarine/pkg/storage/cache/db_cacher"
	"mandarine/pkg/storage/cache/manager"
	mock2 "mandarine/pkg/storage/cache/manager/mock"
	"testing"
)

func Test_DbCacher_Get(t *testing.T) {
	ctx := context.TODO()
	cacheManagerMock := new(mock2.CacheManagerMock)
	cacher := db_cacher.NewDbCacher(cacheManagerMock)

	t.Run(
		"Cache hit", func(t *testing.T) {
			key := "test-key"
			expectedQuery := &caches.Query[any]{ /* инициализация соответствующих полей */ }

			cacheManagerMock.On("Get", ctx, key, mock.Anything).Once().Run(
				func(args mock.Arguments) {
					query := args.Get(2).(*caches.Query[any])
					*query = *expectedQuery
				},
			).Return(nil)

			query := &caches.Query[any]{}
			result, err := cacher.Get(ctx, key, query)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedQuery, result)
		},
	)

	t.Run(
		"Cache miss", func(t *testing.T) {
			key := "missing-key"
			cacheManagerMock.On("Get", ctx, key, mock.Anything).Once().Return(manager.ErrCacheEntryNotFound)

			query := &caches.Query[any]{}
			result, err := cacher.Get(ctx, key, query)
			assert.NoError(t, err)
			assert.Nil(t, result)
		},
	)

	t.Run(
		"Error", func(t *testing.T) {
			key := "error-key"
			expectedErr := errors.New("cache error")
			cacheManagerMock.On("Get", ctx, key, mock.Anything).Once().Return(expectedErr)

			query := &caches.Query[any]{}
			result, err := cacher.Get(ctx, key, query)
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
			assert.Nil(t, result)
		},
	)
}

func Test_DbCacher_Store(t *testing.T) {
	ctx := context.TODO()
	cacheManagerMock := new(mock2.CacheManagerMock)
	cacher := db_cacher.NewDbCacher(cacheManagerMock)

	t.Run(
		"Success", func(t *testing.T) {
			key := "store-key"
			query := &caches.Query[any]{ /* инициализация соответствующих полей */ }

			cacheManagerMock.On("Set", ctx, key, *query).Once().Return(nil)

			err := cacher.Store(ctx, key, query)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"Error", func(t *testing.T) {
			key := "store-key"
			query := &caches.Query[any]{ /* инициализация соответствующих полей */ }
			expectedErr := errors.New("internal error")

			cacheManagerMock.On("Set", ctx, key, *query).Once().Return(expectedErr)

			err := cacher.Store(ctx, key, query)
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)
}

func Test_DbCacher_Invalidate(t *testing.T) {
	ctx := context.TODO()
	cacheManagerMock := new(mock2.CacheManagerMock)
	cacher := db_cacher.NewDbCacher(cacheManagerMock)

	t.Run(
		"Success", func(t *testing.T) {
			cacheManagerMock.On("Invalidate", ctx, fmt.Sprintf("%s*", caches.IdentifierPrefix)).Once().Return(nil)

			err := cacher.Invalidate(ctx)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"Failure", func(t *testing.T) {
			expectedErr := errors.New("invalidate error")
			cacheManagerMock.On(
				"Invalidate", ctx, fmt.Sprintf("%s*", caches.IdentifierPrefix),
			).Once().Return(expectedErr)

			err := cacher.Invalidate(ctx)
			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)
}
