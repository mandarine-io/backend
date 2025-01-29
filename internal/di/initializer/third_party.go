package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	geocodingfactory "github.com/mandarine-io/backend/third_party/geocoding/factory"
	oauthfactory "github.com/mandarine-io/backend/third_party/oauth/factory"
)

func ThirdParty(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup third party")

		var err error
		for _, p := range c.Config.OAuthProviders {
			c.ThirdParties.OAuth[p.Name], err = oauthfactory.NewProviderByKey(p.Name, p.ClientID, p.ClientSecret)
			if err != nil {
				c.Logger.Warn().Err(err).Msgf("failed to create OAuth provider by key %s", p.Name)
			}
		}

		for _, p := range c.Config.GeocodingProviders {
			c.ThirdParties.Geocoding[p.Name], err = geocodingfactory.NewProviderByKey(p.Name, p.APIKey)
			if err != nil {
				c.Logger.Warn().Err(err).Msgf("failed to create geocoding provider by key %s", p.Name)
			}
		}

		return nil
	}
}
