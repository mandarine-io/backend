package memory_test

import (
	"context"
	"github.com/mandarine-io/Backend/pkg/storage/cache/manager"
	"github.com/mandarine-io/Backend/pkg/storage/cache/manager/memory"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CacheManager_Memory(t *testing.T) {
	ttl := time.Minute
	cacheManager := memory.NewCacheManager(ttl)
	ctx := context.Background()

	t.Run(
		"set and get value", func(t *testing.T) {
			key := "test-key"
			value := "test-value"

			err := cacheManager.Set(ctx, key, value)
			assert.NoError(t, err)

			var result string
			err = cacheManager.Get(ctx, key, &result)
			assert.NoError(t, err)
			assert.Equal(t, value, result)
		},
	)

	t.Run(
		"get non-existent key", func(t *testing.T) {
			var result string
			err := cacheManager.Get(ctx, "non-existent-key", &result)
			assert.ErrorIs(t, err, manager.ErrCacheEntryNotFound)
		},
	)

	t.Run(
		"delete key", func(t *testing.T) {
			key := "test-delete-key"
			value := "value-to-delete"

			err := cacheManager.Set(ctx, key, value)
			assert.NoError(t, err)

			err = cacheManager.Delete(ctx, key)
			assert.NoError(t, err)

			var result string
			err = cacheManager.Get(ctx, key, &result)
			assert.ErrorIs(t, err, manager.ErrCacheEntryNotFound)
		},
	)

	t.Run(
		"set with expiration and auto-clean", func(t *testing.T) {
			key := "test-expire-key"
			value := "value-to-expire"

			err := cacheManager.SetWithExpiration(ctx, key, value, 1*time.Second)
			assert.NoError(t, err)

			time.Sleep(2 * time.Second) // Дождитесь истечения срока действия записи

			var result string
			err = cacheManager.Get(ctx, key, &result)
			assert.ErrorIs(t, err, manager.ErrCacheEntryNotFound)
		},
	)

	t.Run(
		"invalidate keys by regex", func(t *testing.T) {
			_ = cacheManager.Set(ctx, "prefix-key1", "value1")
			_ = cacheManager.Set(ctx, "prefix-key2", "value2")
			_ = cacheManager.Set(ctx, "other-key", "value3")

			err := cacheManager.Invalidate(ctx, "^prefix-.*")
			assert.NoError(t, err)

			var result string
			err = cacheManager.Get(ctx, "prefix-key1", &result)
			assert.ErrorIs(t, err, manager.ErrCacheEntryNotFound)

			err = cacheManager.Get(ctx, "prefix-key2", &result)
			assert.ErrorIs(t, err, manager.ErrCacheEntryNotFound)

			err = cacheManager.Get(ctx, "other-key", &result)
			assert.NoError(t, err)
			assert.Equal(t, "value3", result)
		},
	)
}
