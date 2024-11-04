package postgres

import (
	"fmt"
	gormHelper "github.com/mandarine-io/Backend/pkg/storage/database/plugin/gorm"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormConfig struct {
	Address  string
	Username string
	Password string
	DBName   string
}

func MustNewGormDb(cfg *GormConfig) *gorm.DB {
	db, err := gorm.Open(
		postgres.Open(GetDSN(cfg)), &gorm.Config{
			Logger: gormHelper.Logger{},
		},
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to connect to postgres")
	}

	log.Info().Msgf("connected to postgres host %s", cfg.Address)

	return db
}

func CloseGormDb(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func GetDSN(cfg *GormConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Address,
		cfg.DBName,
	)
}
