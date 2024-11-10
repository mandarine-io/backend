package registry

import (
	"github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/pkg/pubsub/memory"
	redis3 "github.com/mandarine-io/Backend/pkg/pubsub/redis"
	redis2 "github.com/mandarine-io/Backend/pkg/storage/cache/redis"
	"github.com/rs/zerolog/log"
)

func setupPubSub(c *Container) {
	log.Debug().Msg("setup pub/sub")
	switch c.Config.PubSub.Type {
	case config.MemoryPubSubType:
		c.PubSub.Agent = memory.NewAgent()
	case config.RedisPubSubType:
		if c.Config.PubSub.Redis == nil {
			log.Fatal().Msg("redis config is nil")
		}

		redisConfig := mapAppPubSubRedisConfigToRedisConfig(&c.Config.PubSub)
		c.PubSub.RDB = redis2.MustNewClient(redisConfig)
		c.PubSub.Agent = redis3.NewAgent(c.PubSub.RDB)
	default:
		log.Fatal().Msgf("unknown pub/sub type: %s", c.Config.PubSub.Type)
	}
}

func mapAppPubSubRedisConfigToRedisConfig(cfg *config.PubSubConfig) *redis2.Config {
	return &redis2.Config{
		Address:  cfg.Redis.Address,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DBIndex:  cfg.Redis.DBIndex,
	}
}
