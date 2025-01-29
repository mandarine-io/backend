package postgres

import (
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/database/gorm/plugin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Address  string
	Username string
	Password string
	DBName   string
}

func NewDb(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(GetDSN(cfg)), &gorm.Config{
			Logger:          plugin.Logger{},
			PrepareStmt:     true,
			CreateBatchSize: 100,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}

func CloseDb(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func GetDSN(cfg Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Address,
		cfg.DBName,
	)
}
