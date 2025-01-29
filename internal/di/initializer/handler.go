package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/transport/http/handler/health"
	"github.com/mandarine-io/backend/internal/transport/http/handler/metrics"
	"github.com/mandarine-io/backend/internal/transport/http/handler/swagger"
	"github.com/mandarine-io/backend/internal/transport/http/handler/v0/account"
	"github.com/mandarine-io/backend/internal/transport/http/handler/v0/auth"
	"github.com/mandarine-io/backend/internal/transport/http/handler/v0/geocoding"
	master_profile "github.com/mandarine-io/backend/internal/transport/http/handler/v0/master/profile"
	master_service "github.com/mandarine-io/backend/internal/transport/http/handler/v0/master/service"
	"github.com/mandarine-io/backend/internal/transport/http/handler/v0/resource"
	"github.com/mandarine-io/backend/internal/transport/http/handler/v0/ws"
	"github.com/rs/zerolog/log"
)

func Handlers(c *di.Container) di.Initializer {
	return func() error {
		log.Debug().Msg("setup handlers")

		c.Handlers = di.Handlers{
			account.NewHandler(
				c.DomainSVCs.Account,
				account.WithLogger(c.Logger.With().Str("handler", "account").Logger()),
			),
			auth.NewHandler(
				c.DomainSVCs.Auth,
				c.Config,
				auth.WithLogger(c.Logger.With().Str("handler", "auth").Logger()),
			),
			health.NewHandler(
				c.DomainSVCs.Health,
				health.WithLogger(c.Logger.With().Str("handler", "health").Logger()),
			),
			geocoding.NewHandler(
				c.DomainSVCs.Geocoding,
				geocoding.WithLogger(c.Logger.With().Str("handler", "geocoding").Logger()),
			),
			master_profile.NewHandler(
				c.DomainSVCs.MasterProfile,
				master_profile.WithLogger(c.Logger.With().Str("handler", "master_profile").Logger()),
			),
			master_service.NewHandler(
				c.DomainSVCs.MasterService,
				master_service.WithLogger(c.Logger.With().Str("handler", "master_service").Logger()),
			),
			metrics.NewHandler(
				metrics.WithLogger(c.Logger.With().Str("handler", "metrics").Logger()),
			),
			resource.NewHandler(
				c.DomainSVCs.Resource,
				resource.WithLogger(c.Logger.With().Str("handler", "resource").Logger()),
			),
			swagger.NewHandler(
				swagger.WithLogger(c.Logger.With().Str("handler", "swagger").Logger()),
			),
			ws.NewHandler(
				c.DomainSVCs.Websocket,
				ws.WithLogger(c.Logger.With().Str("handler", "websocket").Logger()),
			),
		}

		return nil
	}
}
