package finalizer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/database/gorm/postgres"
)

func GormDatabase(c *di.Container) di.Finalizer {
	return func() error {
		c.Logger.Debug().Msg("tear down database")
		return postgres.CloseDb(c.Infrastructure.DB)
	}
}
