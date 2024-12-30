package check

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GormCheck struct {
	db *gorm.DB
}

func NewGormCheck(db *gorm.DB) *GormCheck {
	return &GormCheck{db: db}
}

func (r *GormCheck) Pass() bool {
	log.Debug().Msg("check gorm connection")

	sqlDB, err := r.db.DB()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get sql db")
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to ping sql db")
	}
	return err == nil
}

func (r *GormCheck) Name() string {
	return "database"
}
