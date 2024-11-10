package registry

import (
	"github.com/mandarine-io/Backend/internal/transport/http/handler"
	"github.com/mandarine-io/Backend/internal/transport/http/handler/health"
	"github.com/mandarine-io/Backend/internal/transport/http/handler/swagger"
	"github.com/mandarine-io/Backend/internal/transport/http/handler/v0/account"
	"github.com/mandarine-io/Backend/internal/transport/http/handler/v0/auth"
	masterprofile "github.com/mandarine-io/Backend/internal/transport/http/handler/v0/master/profile"
	"github.com/mandarine-io/Backend/internal/transport/http/handler/v0/resource"
	"github.com/mandarine-io/Backend/internal/transport/http/handler/v0/ws"
	"github.com/rs/zerolog/log"
)

type Handlers []handler.ApiHandler

func setupHandlers(c *Container) {
	log.Debug().Msg("setup handlers")
	c.Handlers = Handlers{
		account.NewHandler(c.SVCs.Account),
		auth.NewHandler(c.SVCs.Auth, c.Config),
		health.NewHandler(c.SVCs.Health),
		masterprofile.NewHandler(c.SVCs.MasterProfile),
		resource.NewHandler(c.SVCs.Resource),
		swagger.NewHandler(),
		ws.NewHandler(c.SVCs.Websocket),
	}
}
