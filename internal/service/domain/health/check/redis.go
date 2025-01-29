package check

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisCheck struct {
	rdb    redis.UniversalClient
	logger zerolog.Logger
}

type RedisOption func(c *RedisCheck) error

func WithRedisLogger(logger zerolog.Logger) RedisOption {
	return func(c *RedisCheck) error {
		c.logger = logger
		return nil
	}
}

func NewRedisCheck(rdb redis.UniversalClient, opts ...RedisOption) (*RedisCheck, error) {
	check := &RedisCheck{
		rdb:    rdb,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(check); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return check, nil
}

func (c *RedisCheck) Pass() bool {
	c.logger.Debug().Msg("check redis connection")

	err := c.rdb.Ping(context.Background()).Err()
	if err != nil {
		c.logger.Error().Stack().Err(err).Msg("failed to ping redis")
	}

	return err == nil
}

func (c *RedisCheck) Name() string {
	return "redis"
}
