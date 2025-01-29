package initializer

import (
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/cache/redis"
	redis2 "github.com/mandarine-io/backend/internal/infrastructure/pubsub/redis"
)

func PubSub(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup pub/sub agent")

		var err error
		redisConfig := toRedisPubSubConfig(c.Config.PubSub)
		c.Infrastructure.PubSubRDB, err = redis.NewClient(redisConfig)
		if err != nil {
			return err
		}

		c.Logger.Info().Msgf("connect to pubsub redis %s", redisConfig.Address)

		c.Infrastructure.PubSubAgent, err = redis2.NewAgent(
			c.Infrastructure.PubSubRDB,
			redis2.WithLogger(c.Logger.With().Str("component", "redis-pubsub").Logger()),
		)

		return err
	}
}

func toRedisPubSubConfig(cfg config.RedisPubSubConfig) redis.Config {
	return redis.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DBIndex:  cfg.DBIndex,
	}
}
