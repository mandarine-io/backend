package registry

import (
	"github.com/mandarine-io/Backend/internal/api/service/account"
	"github.com/mandarine-io/Backend/internal/api/service/auth"
	"github.com/mandarine-io/Backend/internal/api/service/health"
	"github.com/mandarine-io/Backend/internal/api/service/resource"
	"github.com/mandarine-io/Backend/internal/api/service/ws"
	"github.com/rs/zerolog/log"
)

type Services struct {
	Auth     *auth.Service
	Account  *account.Service
	Health   *health.Service
	Resource *resource.Service
	WS       *ws.Service
}

func setupServices(c *Container) {
	log.Debug().Msg("setup services")
	c.SVCs = &Services{
		Auth: auth.NewService(
			c.Repos.User,
			c.Repos.BannedToken,
			c.OauthProviders,
			c.CacheManager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Account: account.NewService(
			c.Repos.User,
			c.CacheManager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Health:   health.NewService(c.DB, c.CacheRDB, c.PubSubRDB, c.S3, c.SmtpSender),
		Resource: resource.NewService(c.S3Client),
		WS:       ws.NewService(c.WebsocketPool),
	}
}
