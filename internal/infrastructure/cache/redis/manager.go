package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type Option func(*manager) error

func WithTTL(ttl time.Duration) Option {
	return func(m *manager) error {
		m.ttl = ttl
		return nil
	}
}

func WithLogger(logger zerolog.Logger) Option {
	return func(m *manager) error {
		m.logger = logger
		return nil
	}
}

type manager struct {
	client redis.UniversalClient
	logger zerolog.Logger
	ttl    time.Duration
}

func NewManager(client redis.UniversalClient, opts ...Option) (cache.Manager, error) {
	m := &manager{
		client: client,
		ttl:    5 * time.Minute,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return m, nil
}

func (m *manager) Get(ctx context.Context, key string, value any) error {
	m.logger.Debug().Msgf("get from cache %s", key)

	res := m.client.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		return cache.ErrCacheEntryNotFound
	}
	if res.Err() != nil {
		return res.Err()
	}

	bytes, err := res.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, value)
}

func (m *manager) Set(ctx context.Context, key string, value any) error {
	return m.SetWithExpiration(ctx, key, value, m.ttl)
}

func (m *manager) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	m.logger.Debug().Msgf("set to cache %s with expiration %s", key, expiration)

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return m.client.Set(ctx, key, jsonValue, expiration).Err()
}

func (m *manager) Delete(ctx context.Context, keys ...string) error {
	m.logger.Debug().Msgf("delete from cache %s", strings.Join(keys, ","))

	for _, key := range keys {
		if err := m.client.Del(ctx, key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) Invalidate(ctx context.Context, keyRegex string) error {
	m.logger.Debug().Msgf("invalidate cache by regex %s", keyRegex)

	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = m.client.Scan(ctx, cursor, keyRegex, 0).Result()
		if err != nil {
			return err
		}
		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if err := m.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return nil
}
