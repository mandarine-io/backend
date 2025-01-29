package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/locale/local"
)

func Locale(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup locale")

		var err error
		c.Infrastructure.LocaleBundle, err = local.NewBundle(
			c.Config.Locale.Path,
			local.WithDefaultLang(c.Config.Locale.Language),
			local.WithLogger(c.Logger.With().Str("component", "locale-bundle").Logger()),
		)

		return err
	}
}
