package check

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCheck struct {
	rdb redis.UniversalClient
}

func NewRedisCheck(rdb redis.UniversalClient) *RedisCheck {
	return &RedisCheck{rdb: rdb}
}

func (r *RedisCheck) Pass() bool {
	log.Debug().Msg("check redis connection")
	err := r.rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to ping redis")
	}
	return err == nil
}

func (r *RedisCheck) Name() string {
	return "redis"
}
