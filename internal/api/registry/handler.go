package registry

import (
	"mandarine/internal/api/rest/handler"
	"mandarine/internal/api/rest/handler/health"
	"mandarine/internal/api/rest/handler/swagger"
	"mandarine/internal/api/rest/handler/v0/account"
	"mandarine/internal/api/rest/handler/v0/auth"
	"mandarine/internal/api/rest/handler/v0/resource"
)

type Handlers []handler.ApiHandler

func newHandlers(c *Container) Handlers {
	return Handlers{
		auth.NewHandler(c.Services.Auth, c.Config),
		account.NewHandler(c.Services.Account),
		resource.NewHandler(c.Services.Resource),
		swagger.NewHandler(),
		health.NewHandler(c.DB, c.RedisClient, c.S3, c.SmtpSender),
	}
}
