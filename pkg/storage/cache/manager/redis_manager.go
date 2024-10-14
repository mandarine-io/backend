package manager

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCacheManager struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCacheManager(client *redis.Client, ttl time.Duration) CacheManager {
	return &RedisCacheManager{client: client, ttl: ttl}
}

func (r *RedisCacheManager) Get(ctx context.Context, key string, value interface{}) error {
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		if errors.Is(res.Err(), redis.Nil) {
			return ErrCacheEntryNotFound
		}
		return res.Err()
	}

	bytes, err := res.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, value)
}

func (r *RedisCacheManager) Set(ctx context.Context, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, r.ttl).Err()
}

func (r *RedisCacheManager) SetWithExpiration(
	ctx context.Context, key string, value interface{}, expiration time.Duration,
) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, expiration).Err()
}

func (r *RedisCacheManager) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisCacheManager) Invalidate(ctx context.Context, keyRegex string) error {
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
