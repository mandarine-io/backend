package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/scheduler"
)

func Scheduler(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup scheduler")

		var err error
		c.Infrastructure.Scheduler, err = scheduler.NewSchedulerWithLogger(
			c.Logger.With().Str("component", "scheduler").Logger(),
		)
		if err != nil {
			return err
		}

		c.Infrastructure.Scheduler.Start()

		return nil
	}
}
