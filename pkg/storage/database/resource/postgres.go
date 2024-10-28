package resource

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

func MustConnectPostgres(cfg *PostgresConfig) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(GetDSN(cfg)), &gorm.Config{
			Logger: dbLogger{},
		},
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to connect to postgres")
	}

	log.Info().Msgf("connected to postgres host %s:%d", cfg.Host, cfg.Port)

	return db
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func GetDSN(cfg *PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port,
		cfg.DBName,
	)
}
