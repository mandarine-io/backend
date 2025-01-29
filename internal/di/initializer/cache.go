package initializer

import (
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/di"
	redis2 "github.com/mandarine-io/backend/internal/infrastructure/cache/redis"
	"time"
)

func Cache(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup cache manager")

		var err error
		redisConfig := toRedisCacheConfig(c.Config.Cache)
		c.Infrastructure.CacheRDB, err = redis2.NewClient(redisConfig)
		if err != nil {
			return err
		}

		c.Logger.Info().Msgf("connect to cache redis %s", redisConfig.Address)

		c.Infrastructure.CacheManager, err = redis2.NewManager(
			c.Infrastructure.CacheRDB,
			redis2.WithTTL(time.Duration(c.Config.Cache.TTL)*time.Second),
			redis2.WithLogger(c.Logger.With().Str("component", "redis-cache-manager").Logger()),
		)

		return err
	}
}

func toRedisCacheConfig(cfg config.RedisCacheConfig) redis2.Config {
	return redis2.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DBIndex:  cfg.DBIndex,
	}
}
