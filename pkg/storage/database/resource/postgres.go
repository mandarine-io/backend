package resource

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"mandarine/pkg/logging"
	"os"
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
		slog.Error("Postgres client creation error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	slog.Info("Connected to Postgres host " + addr)

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
