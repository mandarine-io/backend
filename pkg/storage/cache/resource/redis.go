package resource

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBIndex  int
}

func MustConnectRedis(cfg *RedisConfig) *redis.Client {
	// Create client
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.DBIndex,
		},
	)

	log.Info().Msgf("connected to redis host %s", addr)

	return redisClient
}
