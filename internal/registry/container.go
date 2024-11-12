package registry

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/locale"
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/mandarine-io/Backend/pkg/pubsub"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/mandarine-io/Backend/pkg/storage/database/postgres"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	"github.com/mandarine-io/Backend/pkg/template"
	"github.com/mandarine-io/Backend/pkg/websocket"
	"github.com/minio/minio-go/v7"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Container struct {
	Config             *config.Config
	Logger             *zerolog.Logger
	Bundle             *i18n.Bundle
	DB                 *gorm.DB
	WebsocketPool      *websocket.Pool
	HttpClient         *resty.Client
	OauthProviders     map[string]oauth.Provider
	GeocodingProviders map[string]geocoding.Provider
	SmtpSender         smtp.Sender
	TemplateEngine     template.Engine

	Cache struct {
		RDB     redis.UniversalClient
		Manager cache.Manager
	}
	S3 struct {
		Minio  *minio.Client
		Client s3.Client
	}
	PubSub struct {
		RDB   redis.UniversalClient
		Agent pubsub.Agent
	}

	Repos    *Repositories
	SVCs     *Services
	Handlers Handlers
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) MustInitialize(cfg *config.Config) {
	// Setup config
	log.Debug().Msg("setup config")
	c.Config = cfg

	// Setup logger
	log.Debug().Msg("setup logger")
	c.Logger = &log.Logger

	// Setup locale
	log.Debug().Msg("setup locale")
	localeConfig := mapAppLocaleConfigToLocaleConfig(&cfg.Locale)
	c.Bundle = locale.MustLoadLocales(localeConfig)

	// Setup template engine
	log.Debug().Msg("setup template engine")
	templateConfig := mapAppTemplateConfigToTemplateConfig(&cfg.Template)
	c.TemplateEngine = template.MustLoadTemplates(templateConfig)

	// Setup websocket pool
	log.Debug().Msg("setup websocket pool")
	c.WebsocketPool = websocket.NewPool(cfg.Websocket.PoolSize)

	// Setup SMTP sender
	log.Debug().Msg("setup smtp sender")
	smtpConfig := mapAppSmtpConfigToSmtpConfig(&cfg.SMTP)
	c.SmtpSender = smtp.MustNewSender(smtpConfig)

	// Setup HTTP client
	log.Debug().Msg("setup http client")
	c.HttpClient = resty.New()

	setupCacheManager(c)
	setupDatabase(c)
	setupS3(c)
	setupPubSub(c)
	setupOAuthClients(c)
	setupGeocodingClients(c)
	setupGormRepositories(c)
	setupServices(c)
	setupHandlers(c)
}

func (c *Container) Close() error {
	var errs []error

	if err := postgres.CloseGormDb(c.DB); err != nil {
		log.Error().Stack().Err(err).Msg("failed to close postgres connection")
		errs = append(errs, err)
	} else {
		log.Info().Msg("postgres connection is closed")
	}

	if c.Cache.RDB != nil {
		if err := c.Cache.RDB.Close(); err != nil {
			log.Error().Stack().Err(err).Msg("failed to close redis cache connection")
			errs = append(errs, err)
		} else {
			log.Info().Msg("redis cache connection is closed")
		}
	}

	if c.PubSub.RDB != nil {
		if err := c.PubSub.RDB.Close(); err != nil {
			log.Error().Stack().Err(err).Msg("failed to close redis pub/sub connection")
			errs = append(errs, err)
		} else {
			log.Info().Msg("redis pub/sub connection is closed")
		}
	}

	if err := c.PubSub.Agent.Close(); err != nil {
		log.Error().Stack().Err(err).Msg("failed to close pub/sub")
		errs = append(errs, err)
	} else {
		log.Info().Msg("pub/sub is closed")
	}

	if err := c.WebsocketPool.Close(); err != nil {
		log.Error().Stack().Err(err).Msg("failed to close websocket pool")
		errs = append(errs, err)
	} else {
		log.Info().Msg("websocket pool is closed")
	}

	return errors.Join(errs...)
}

func mapAppLocaleConfigToLocaleConfig(cfg *config.LocaleConfig) *locale.Config {
	return &locale.Config{
		Path:     cfg.Path,
		Language: cfg.Language,
	}
}

func mapAppTemplateConfigToTemplateConfig(cfg *config.TemplateConfig) *template.Config {
	return &template.Config{
		Path: cfg.Path,
	}
}

func mapAppSmtpConfigToSmtpConfig(cfg *config.SmtpConfig) *smtp.Config {
	return &smtp.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
		SSL:      cfg.SSL,
	}
}
