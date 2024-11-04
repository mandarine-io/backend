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

type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (l *migrateLogger) Verbose() bool {
	return log.Logger.GetLevel() == zerolog.DebugLevel
}

func Migrate(dsn string, migrationDir string) error {
	log.Info().Msg("migrating database")

	sourceUrl := fmt.Sprintf("file://%s", migrationDir)
	migrate, err := goMigrate.New(sourceUrl, dsn)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, dbErr := migrate.Close()
		if sourceErr != nil {
			log.Warn().Stack().Err(sourceErr).Msg("failed to close source")
		}
		if dbErr != nil {
			log.Warn().Stack().Err(dbErr).Msg("failed to close database")
		}
	}()

	migrate.Log = &migrateLogger{}

	if err = migrate.Up(); err != nil && !errors.Is(err, goMigrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, goMigrate.ErrNoChange) {
		log.Info().Msg("migrations are already installed")
	} else {
		log.Info().Msg("migrations installed successfully")
	}

	return nil
}
