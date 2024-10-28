package registry

import (
	"github.com/mandarine-io/Backend/internal/api/rest/handler"
	"github.com/mandarine-io/Backend/internal/api/rest/handler/health"
	"github.com/mandarine-io/Backend/internal/api/rest/handler/swagger"
	"github.com/mandarine-io/Backend/internal/api/rest/handler/v0/account"
	"github.com/mandarine-io/Backend/internal/api/rest/handler/v0/auth"
	"github.com/mandarine-io/Backend/internal/api/rest/handler/v0/resource"
	"github.com/mandarine-io/Backend/internal/api/rest/handler/ws"
	"github.com/rs/zerolog/log"
)

type Handlers []handler.ApiHandler

func newHandlers(c *Container) Handlers {
	log.Debug().Msg("setup handlers")
	return Handlers{
		auth.NewHandler(c.Services.Auth, c.Config),
		account.NewHandler(c.Services.Account),
		resource.NewHandler(c.Services.Resource),
		swagger.NewHandler(),
		health.NewHandler(c.DB, c.RedisClient, c.S3, c.SmtpSender),
		ws.NewHandler(c.Services.WS),
	}
}
