package finalizer

import (
	"github.com/mandarine-io/backend/internal/di"
)

func Scheduler(c *di.Container) di.Finalizer {
	return func() error {
		c.Logger.Debug().Msg("tear down scheduler")

		return c.Infrastructure.Scheduler.Shutdown()
	}
}
