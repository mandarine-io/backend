package service

import (
	"mandarine/internal/api/config"
	"mandarine/internal/api/persistence/repo"
	"mandarine/internal/api/service/account"
	"mandarine/internal/api/service/auth"
	"mandarine/pkg/oauth"
	"mandarine/pkg/oauth/google"
	"mandarine/pkg/oauth/mailru"
	"mandarine/pkg/oauth/yandex"
	"mandarine/pkg/smtp"
	"mandarine/pkg/storage/cache/manager"
	"mandarine/pkg/template"
)

type Services struct {
	Login         *auth.LoginService
	Register      *auth.RegisterService
	ResetPassword *auth.ResetPasswordService
	SocialLogins  map[string]*auth.SocialLoginService
	Logout        *auth.LogoutService

	Account *account.AccountService
}

func NewServices(
	repos *repo.Repositories,
	cacheManager manager.CacheManager,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	oauthProviders map[string]oauth.Provider,
	cfg *config.Config,
) *Services {
	return &Services{
		Login:         auth.NewLoginService(repos.User, cfg),
		Register:      auth.NewRegisterService(repos.User, cacheManager, smtpSender, templateEngine, cfg),
		ResetPassword: auth.NewResetPasswordService(repos.User, cacheManager, smtpSender, templateEngine, cfg),
		SocialLogins: map[string]*auth.SocialLoginService{
			google.ProviderKey: auth.NewSocialLoginService(repos.User, oauthProviders[google.ProviderKey], google.ProviderKey, cfg),
			yandex.ProviderKey: auth.NewSocialLoginService(repos.User, oauthProviders[yandex.ProviderKey], yandex.ProviderKey, cfg),
			mailru.ProviderKey: auth.NewSocialLoginService(repos.User, oauthProviders[mailru.ProviderKey], mailru.ProviderKey, cfg),
		},
		Logout: auth.NewLogoutService(repos.BannedToken, cfg),

		Account: account.NewAccountService(repos.User, cacheManager, smtpSender, templateEngine, cfg),
	}
}
