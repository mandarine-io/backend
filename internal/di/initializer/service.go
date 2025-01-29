package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/service/domain/account"
	"github.com/mandarine-io/backend/internal/service/domain/auth"
	"github.com/mandarine-io/backend/internal/service/domain/geocoding"
	"github.com/mandarine-io/backend/internal/service/domain/health"
	masterprofile "github.com/mandarine-io/backend/internal/service/domain/master/profile"
	masterservice "github.com/mandarine-io/backend/internal/service/domain/master/service"
	"github.com/mandarine-io/backend/internal/service/domain/resource"
	"github.com/mandarine-io/backend/internal/service/domain/ws"
	"github.com/mandarine-io/backend/internal/service/infrastructure/jwt"
	"github.com/mandarine-io/backend/internal/service/infrastructure/otp"
	geocoding2 "github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/mandarine-io/backend/third_party/geocoding/factory"
	"github.com/rs/zerolog/log"
)

func Services(c *di.Container) di.Initializer {
	return func() error {
		log.Debug().Msg("setup infrastructure services")

		c.InfrastructureSVCs = di.InfrastructureServices{
			JWT: jwt.NewService(
				c.Infrastructure.CacheManager,
				c.Config.Security.JWT,
				jwt.WithLogger(c.Logger.With().Str("infra-service", "jwt").Logger()),
			),
			OTP: otp.NewService(
				c.Infrastructure.CacheManager,
				c.Config.Security.OTP,
				otp.WithLogger(c.Logger.With().Str("infra-service", "otp").Logger()),
			),
		}

		log.Debug().Msg("setup domain services")

		geocodingProviders := make([]geocoding2.Provider, 0)
		for _, provider := range c.ThirdParties.Geocoding {
			geocodingProviders = append(geocodingProviders, provider)
		}

		c.DomainSVCs = di.DomainServices{
			Account: account.NewService(
				c.Config,
				c.Repos.User,
				c.Infrastructure.SMTPSender,
				c.Infrastructure.TemplateEngine,
				c.InfrastructureSVCs.OTP,
				account.WithLogger(c.Logger.With().Str("domain-service", "account").Logger()),
			),
			Auth: auth.NewService(
				c.Config,
				c.Infrastructure.SMTPSender,
				c.Infrastructure.TemplateEngine,
				c.Repos.User,
				c.InfrastructureSVCs.JWT,
				c.InfrastructureSVCs.OTP,
				c.ThirdParties.OAuth,
				auth.WithLogger(c.Logger.With().Str("domain-service", "auth").Logger()),
			),
			Health: health.NewService(
				c.Infrastructure.DB,
				c.Infrastructure.CacheRDB,
				c.Infrastructure.PubSubRDB,
				c.Infrastructure.MinioClient,
				c.Infrastructure.SMTPDialer,
				health.WithLogger(c.Logger.With().Str("domain-service", "health").Logger()),
			),
			Geocoding: geocoding.NewService(
				c.Infrastructure.CacheManager,
				factory.NewProviderChained(geocodingProviders...),
				geocoding.WithLogger(c.Logger.With().Str("domain-service", "geocoding").Logger()),
			),
			MasterProfile: masterprofile.NewService(
				c.Repos.MasterProfile,
				masterprofile.WithLogger(c.Logger.With().Str("domain-service", "master-profile").Logger()),
			),
			MasterService: masterservice.NewService(
				c.Repos.MasterProfile,
				c.Repos.MasterService,
				masterservice.WithLogger(c.Logger.With().Str("domain-service", "master-service").Logger()),
			),
			Resource: resource.NewService(
				c.Infrastructure.S3Manager,
				resource.WithLogger(c.Logger.With().Str("domain-service", "resource").Logger()),
			),
			Websocket: ws.NewService(
				c.Infrastructure.WSPool,
				ws.WithLogger(c.Logger.With().Str("domain-service", "websocket").Logger()),
			),
		}

		return nil
	}
}
