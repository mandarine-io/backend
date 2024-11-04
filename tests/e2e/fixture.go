package e2e

import (
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func MustNewFixtures(db *gorm.DB, paths ...string) *testfixtures.Loader {
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get sql db")
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDb),
		testfixtures.Dialect("postgres"),
		testfixtures.Paths(paths...),
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to create fixtures")
	}

	return fixtures
}

func MustLoadFixtures(fixtures *testfixtures.Loader) {
	if err := fixtures.Load(); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load fixtures")
	}
}
