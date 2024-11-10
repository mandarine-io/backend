package redis

import (
	"context"
	"encoding/json"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type manager struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func NewManager(client redis.UniversalClient, ttl time.Duration) cache.Manager {
	return &manager{client: client, ttl: ttl}
}

func (r *manager) Get(ctx context.Context, key string, value interface{}) error {
	log.Debug().Msgf("get from cache %s", key)

	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		if errors.Is(res.Err(), redis.Nil) {
			return cache.ErrCacheEntryNotFound
		}
		return res.Err()
	}

	bytes, err := res.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, value)
}

func (r *manager) Set(ctx context.Context, key string, value interface{}) error {
	return r.SetWithExpiration(ctx, key, value, r.ttl)
}

func (r *manager) SetWithExpiration(
	ctx context.Context, key string, value interface{}, expiration time.Duration,
) error {
	log.Debug().Msgf("set to cache %s with expiration %s", key, expiration)

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, expiration).Err()
}

func (r *manager) Delete(ctx context.Context, keys ...string) error {
	log.Debug().Msgf("delete from cache %s", strings.Join(keys, ","))

	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *manager) Invalidate(ctx context.Context, keyRegex string) error {
	log.Debug().Msgf("invalidate cache by regex %s", keyRegex)

	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = r.client.Scan(ctx, cursor, keyRegex, 0).Result()
		if err != nil {
			return err
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if err := r.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return nil
}
