package check

import (
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type GormCheck struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type GormOption func(c *GormCheck) error

func WithGormLogger(logger zerolog.Logger) GormOption {
	return func(c *GormCheck) error {
		c.logger = logger
		return nil
	}
}

func NewGormCheck(db *gorm.DB, opts ...GormOption) (*GormCheck, error) {
	check := &GormCheck{
		db:     db,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(check); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return check, nil
}

func (c *GormCheck) Pass() bool {
	c.logger.Debug().Msg("check gorm connection")

	sqlDB, err := c.db.DB()
	if err != nil {
		c.logger.Error().Stack().Err(err).Msg("failed to get sql db")
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		c.logger.Error().Stack().Err(err).Msg("failed to ping sql db")
	}
	return err == nil
}

func (c *GormCheck) Name() string {
	return "database"
}
