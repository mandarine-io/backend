package registry

import (
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler/health"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler/swagger"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler/v0/account"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler/v0/auth"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler/v0/resource"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler/ws"
	"github.com/rs/zerolog/log"
)

type Handlers []handler.ApiHandler

func setupHandlers(c *Container) {
	log.Debug().Msg("setup handlers")
	c.Handlers = Handlers{
		auth.NewHandler(c.SVCs.Auth, c.Config),
		account.NewHandler(c.SVCs.Account),
		resource.NewHandler(c.SVCs.Resource),
		swagger.NewHandler(),
		health.NewHandler(c.SVCs.Health),
		ws.NewHandler(c.SVCs.WS),
	}
}
