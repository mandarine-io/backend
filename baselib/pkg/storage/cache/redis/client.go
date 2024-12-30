package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address  string
	Username string
	Password string
	DBIndex  int
}

func MustNewClient(cfg *Config) *redis.Client {
	// Create client
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     cfg.Address,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.DBIndex,
		},
	)

	log.Info().Msgf("connected to redis host %s", cfg.Address)

	return redisClient
}
