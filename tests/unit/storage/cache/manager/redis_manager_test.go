package manager_test

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redismock/v9"
	"mandarine/pkg/storage/cache/manager"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CacheManager_Redis(t *testing.T) {
	ctx := context.Background()
	client, mock := redismock.NewClientMock()
	ttl := time.Minute
	cacheManager := manager.NewRedisCacheManager(client, ttl)

	t.Run(
		"set and get value", func(t *testing.T) {
			key := "test-key"
			value := "test-value"

			jsonValue, _ := json.Marshal(value)

			mock.ExpectSet(key, jsonValue, 60*time.Second).SetVal("OK")
			mock.ExpectGet(key).SetVal(string(jsonValue))

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
			key := "non-existent-key"

			mock.ExpectGet(key).RedisNil()

			var result string
			err := cacheManager.Get(ctx, key, &result)
			assert.ErrorIs(t, err, manager.ErrCacheEntryNotFound)
		},
	)

	t.Run(
		"delete key", func(t *testing.T) {
			key := "test-delete-key"

			mock.ExpectDel(key).SetVal(1)

			err := cacheManager.Delete(ctx, key)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"set with expiration and get value", func(t *testing.T) {
			key := "test-expire-key"
			value := "value-to-expire"
			expiration := 60 * time.Second

			jsonValue, _ := json.Marshal(value)

			mock.ExpectSet(key, jsonValue, expiration).SetVal("OK")
			mock.ExpectGet(key).SetVal(string(jsonValue))

			err := cacheManager.SetWithExpiration(ctx, key, value, expiration)
			assert.NoError(t, err)

			var result string
			err = cacheManager.Get(ctx, key, &result)
			assert.NoError(t, err)
			assert.Equal(t, value, result)
		},
	)

	t.Run(
		"invalidate keys by regex", func(t *testing.T) {
			key1 := "prefix-key1"
			key2 := "prefix-key2"

			mock.ExpectScan(0, "prefix-*", 0).SetVal([]string{key1, key2}, 0)
			mock.ExpectDel(key1, key2).SetVal(2)

			err := cacheManager.Invalidate(ctx, "prefix-*")
			assert.NoError(t, err)
		},
	)
}
