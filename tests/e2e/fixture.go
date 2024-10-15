package e2e

import (
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/gorm"
	"log/slog"
	"mandarine/pkg/logging"
	"os"
)

func MustNewFixtures(db *gorm.DB, paths ...string) *testfixtures.Loader {
	sqlDb, err := db.DB()
	if err != nil {
		slog.Error("Postgres connection error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDb),
		testfixtures.Dialect("postgres"),
		testfixtures.Paths(paths...),
	)
	if err != nil {
		slog.Error("Fixtures creation error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	return fixtures
}

func MustLoadFixtures(fixtures *testfixtures.Loader) {
	if err := fixtures.Load(); err != nil {
		slog.Error("Fixtures loading error", logging.ErrorAttr(err))
		os.Exit(1)
	}
}
