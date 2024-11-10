package registry

import (
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/mandarine-io/Backend/pkg/oauth/factory"
	"github.com/rs/zerolog/log"
)

func setupOAuthClients(c *Container) {
	log.Debug().Msg("setup oauth providers")
	c.OauthProviders = map[string]oauth.Provider{}
	for k, v := range c.Config.OAuthClients {
		log.Debug().Msgf("setup oauth provider: %s", k)
		provider := factory.NewProviderByKey(k, v.ClientID, v.ClientSecret)
		if provider == nil {
			log.Warn().Msgf("unknown oauth provider: %s", k)
			continue
		}
		c.OauthProviders[k] = provider
	}
}
