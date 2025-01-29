package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Address  string
	Username string
	Password string
	DBIndex  int
}

func NewClient(cfg Config) (redis.UniversalClient, error) {
	rdb := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    []string{cfg.Address},
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.DBIndex,
		},
	)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}
