package registry

import (
	accountSvc "mandarine/internal/api/service/account"
	authSvc "mandarine/internal/api/service/auth"
	resourceSvc "mandarine/internal/api/service/resource"
)

type Services struct {
	Auth     *authSvc.Service
	Account  *accountSvc.Service
	Resource *resourceSvc.Service
}

func newServices(c *Container) *Services {
	return &Services{
		Auth: authSvc.NewService(
			c.Repositories.User,
			c.Repositories.BannedToken,
			c.OauthProviders,
			c.CacheManager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Account: accountSvc.NewService(
			c.Repositories.User,
			c.CacheManager,
			c.SmtpSender,
			c.TemplateEngine,
			c.Config,
		),
		Resource: resourceSvc.NewService(c.S3Client),
	}
}
