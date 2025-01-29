package database

import (
	"errors"
	"fmt"
	goMigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type migrateLogger struct {
	logger zerolog.Logger
}

func (l *migrateLogger) Printf(format string, v ...any) {
	if len(format) > 0 && format[len(format)-1] == '\n' {
		format = format[0 : len(format)-1]
	}
	l.logger.Info().Msgf(format, v...)
}

func (l *migrateLogger) Verbose() bool {
	return log.Logger.GetLevel() == zerolog.DebugLevel
}

func Migrate(dsn string, migrationDir string) error {
	logger := &migrateLogger{
		log.With().Str("component", "db-migrator").Logger(),
	}
	logger.Printf("migrating database")

	sourceURL := fmt.Sprintf("file://%s", migrationDir)

	migrate, err := goMigrate.New(sourceURL, dsn)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = migrate.Close()
	}()

	migrate.Log = logger
	if err = migrate.Up(); err != nil && !errors.Is(err, goMigrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, goMigrate.ErrNoChange) {
		logger.Printf("migrations are already installed")
	} else {
		logger.Printf("migrations installed successfully")
	}

	return nil
}
