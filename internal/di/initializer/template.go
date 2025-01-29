package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/template/local"
)

func Template(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup template engine")

		var err error
		c.Infrastructure.TemplateEngine, err = local.NewEngine(
			c.Config.Template.Path,
			local.WithLogger(c.Logger.With().Str("component", "template-engine").Logger()),
		)

		return err
	}
}
