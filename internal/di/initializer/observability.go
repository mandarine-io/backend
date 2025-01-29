package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/observability"
)

func Metrics(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup metrics adapter")

		c.Metrics = observability.NewMetricAdapter(
			observability.WithLogger(c.Logger.With().Str("observability", "metrics_adapter").Logger()),
		)

		return nil
	}
}
