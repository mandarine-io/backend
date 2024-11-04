package registry

import (
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/pkg/storage/cache/db_cacher"
	"github.com/mandarine-io/Backend/pkg/storage/database"
	gormPlugin "github.com/mandarine-io/Backend/pkg/storage/database/plugin/gorm"
	"github.com/mandarine-io/Backend/pkg/storage/database/postgres"
	"github.com/rs/zerolog/log"
)

func setupDatabase(c *Container) {
	log.Debug().Msg("setup database")
	var dsn string
	switch c.Config.Database.Type {
	case config.PostgresDatabaseType:
		if c.Config.Database.Postgres == nil {
			log.Fatal().Msg("postgres config is nil")
		}
		postgresConfig := mapAppPostgresConfigToPostgresGormConfig(&c.Config.Database)
		c.DB = postgres.MustNewGormDb(postgresConfig)
		dsn = postgres.GetDSN(postgresConfig)
	default:
		log.Fatal().Msgf("unknown database type: %s", c.Config.Database.Type)
	}

	// Setup GORM plugins
	err := gormPlugin.UseCachePlugin(c.DB, db_cacher.NewDbCacher(c.CacheManager))
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to use cache plugin")
	}

	// Migrate database
	log.Debug().Msg("migrate database")
	err = database.Migrate(dsn, c.Config.Migrations.Path)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to migrate database")
	}
}

func mapAppPostgresConfigToPostgresGormConfig(cfg *config.DatabaseConfig) *postgres.GormConfig {
	return &postgres.GormConfig{
		Address:  cfg.Postgres.Address,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
	}
}
