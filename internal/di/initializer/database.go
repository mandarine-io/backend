package initializer

import (
	"github.com/go-gorm/caches/v4"
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/cache/dbcacher"
	"github.com/mandarine-io/backend/internal/infrastructure/database"
	"github.com/mandarine-io/backend/internal/infrastructure/database/gorm/plugin"
	"github.com/mandarine-io/backend/internal/infrastructure/database/gorm/postgres"
	"github.com/rs/zerolog/log"
)

func GormDatabase(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup database")

		var err error
		postgresConfig := toPostgresGormConfig(c.Config.Database)
		c.Infrastructure.DB, err = postgres.NewDb(postgresConfig)
		if err != nil {
			return err
		}

		c.Logger.Info().Msgf("connect to postgres %s", postgresConfig.Address)

		// Setup GORM plugins
		if c.Infrastructure.CacheManager != nil {
			var dbCacher caches.Cacher
			dbCacher, err = dbcacher.NewDbCacher(c.Infrastructure.CacheManager)
			if err != nil {
				return err
			}

			err = plugin.UseCachePlugin(c.Infrastructure.DB, dbCacher)
			if err != nil {
				c.Logger.Warn().Err(err).Msg("failed to use cache plugin")
			}
		}

		// Migrate database
		log.Debug().Msg("migrate database")
		dsn := postgres.GetDSN(postgresConfig)
		return database.Migrate(dsn, c.Config.Migrations.Path)
	}
}

func toPostgresGormConfig(cfg config.PostgresDatabaseConfig) postgres.Config {
	return postgres.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DBName:   cfg.DBName,
	}
}
