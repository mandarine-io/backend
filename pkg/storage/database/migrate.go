package database

import (
	"errors"
	"fmt"
	"log/slog"
	"mandarine/pkg/logging"

	goMigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...))
}

func (l *migrateLogger) Verbose() bool {
	return false
}

func Migrate(dsn string, migrationDir string) error {
	migrate, err := goMigrate.New(fmt.Sprintf("file://%s", migrationDir), dsn)
	if err != nil {
		return err
	}
	defer func(migrate *goMigrate.Migrate) {
		sourceErr, dbErr := migrate.Close()
		if sourceErr != nil {
			slog.Error("Migrate close error", logging.ErrorAttr(sourceErr))
		}
		if dbErr != nil {
			slog.Error("Migrate close error", logging.ErrorAttr(dbErr))
		}
	}(migrate)

	migrate.Log = &migrateLogger{}

	if err = migrate.Up(); err != nil && !errors.Is(err, goMigrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, goMigrate.ErrNoChange) {
		slog.Info("Migrations are already installed")
	} else {
		slog.Info("Migrations installed successfully")
	}

	return nil
}
