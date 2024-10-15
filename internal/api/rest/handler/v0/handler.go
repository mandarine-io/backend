package v0

import (
	"mandarine/internal/api/config"
	"mandarine/internal/api/rest/handler/v0/account"
	"mandarine/internal/api/rest/handler/v0/auth"
	"mandarine/internal/api/service"
)

type Handlers struct {
	Register      auth.RegisterHandler
	Login         auth.LoginHandler
	SocialLogin   auth.SocialLoginHandler
	ResetPassword auth.ResetPasswordHandler
	Logout        auth.LogoutHandler

	Account account.AccountHandler
}

func NewHandlers(services *service.Services, cfg *config.Config) *Handlers {
	return &Handlers{
		Login:         auth.NewLoginHandler(services.Login, cfg),
		SocialLogin:   auth.NewSocialLoginHandler(services.SocialLogins, cfg),
		Register:      auth.NewRegisterHandler(services.Register),
		ResetPassword: auth.NewResetPasswordHandler(services.ResetPassword),
		Logout:        auth.NewLogoutHandler(services.Logout),

		Account: account.NewAccountHandler(services.Account),
	}
}
