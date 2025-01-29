package finalizer

import (
	"github.com/mandarine-io/backend/internal/di"
)

func Cache(c *di.Container) di.Finalizer {
	return func() error {
		c.Logger.Debug().Msg("tear down cache manager")

		if c.Infrastructure.CacheRDB == nil {
			return nil
		}

		return c.Infrastructure.CacheRDB.Close()
	}
}
