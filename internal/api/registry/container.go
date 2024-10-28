package registry

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/pkg/locale"
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/mandarine-io/Backend/pkg/oauth/google"
	"github.com/mandarine-io/Backend/pkg/oauth/mailru"
	"github.com/mandarine-io/Backend/pkg/oauth/yandex"
	redis3 "github.com/mandarine-io/Backend/pkg/pubsub/redis"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/mandarine-io/Backend/pkg/storage/cache/db_cacher"
	"github.com/mandarine-io/Backend/pkg/storage/cache/manager"
	redis2 "github.com/mandarine-io/Backend/pkg/storage/cache/manager/redis"
	cacheResource "github.com/mandarine-io/Backend/pkg/storage/cache/resource"
	"github.com/mandarine-io/Backend/pkg/storage/database"
	dbResource "github.com/mandarine-io/Backend/pkg/storage/database/resource"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	s3Resource "github.com/mandarine-io/Backend/pkg/storage/s3/resource"
	"github.com/mandarine-io/Backend/pkg/template"
	"github.com/minio/minio-go/v7"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type Container struct {
	Config *config.Config

	Logger        *zerolog.Logger
	Bundle        *i18n.Bundle
	RedisClient   *redis.Client
	DB            *gorm.DB
	S3            *minio.Client
	WebsocketPool *websocket.Pool
	HttpClient    *resty.Client

	OauthProviders map[string]oauth.Provider
	TemplateEngine template.Engine
	CacheManager   manager.CacheManager
	S3Client       s3.Client
	SmtpSender     smtp.Sender
	PubSub         pubsub.Agent

	Repositories *Repositories
	Services     *Services
	Handlers     Handlers
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

	// Setup cache manager
	log.Debug().Msg("setup cache manager")
	redisConfig := mapAppRedisConfigToRedisConfig(&cfg.Redis)
	c.RedisClient = cacheResource.MustConnectRedis(redisConfig)
	c.CacheManager = redis2.NewCacheManager(c.RedisClient, time.Duration(cfg.Cache.TTL)*time.Second)

	// Setup database
	log.Debug().Msg("setup database")
	postgresConfig := mapAppPostgresConfigToPostgresConfig(&cfg.Postgres)
	c.DB = dbResource.MustConnectPostgres(postgresConfig)
	err := database.UseCachePlugin(c.DB, db_cacher.NewDbCacher(c.CacheManager))
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to use cache plugin")
	}

	// Migrate database
	log.Debug().Msg("migrate database")
	err = database.Migrate(dbResource.GetDSN(postgresConfig), cfg.Migrations.Path)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to migrate database")
	}

	// Setup S3
	log.Debug().Msg("setup s3")
	minioConfig := mapAppMinioConfigToMinioConfig(&cfg.Minio)
	c.S3 = s3Resource.MustConnectMinio(minioConfig)
	c.S3Client = s3.NewClient(c.S3, cfg.Minio.BucketName)

	// Setup SMTP sender
	log.Debug().Msg("setup smtp sender")
	smtpConfig := mapAppSmtpConfigToSmtpConfig(&cfg.SMTP)
	c.SmtpSender = smtp.MustNewSender(smtpConfig)

	// Setup HTTP client
	log.Debug().Msg("setup http client")
	c.HttpClient = resty.New()

	// Setup OAuth providers
	log.Debug().Msg("setup oauth providers")
	c.OauthProviders = map[string]oauth.Provider{
		google.ProviderKey: google.NewOAuthGoogleProvider(cfg.OAuthClient.Google.ClientID, cfg.OAuthClient.Google.ClientSecret),
		yandex.ProviderKey: yandex.NewOAuthYandexProvider(cfg.OAuthClient.Yandex.ClientID, cfg.OAuthClient.Yandex.ClientSecret),
		mailru.ProviderKey: mailru.NewOAuthMailRuProvider(cfg.OAuthClient.MailRu.ClientID, cfg.OAuthClient.MailRu.ClientSecret),
	}

	// Setup CSR components
	log.Debug().Msg("setup csr components")
	c.Repositories = newGormRepositories(c.DB)
	c.Services = newServices(c)
	c.Handlers = newHandlers(c)
}

func (c *Container) Close() error {
	var errs []error

	if err := dbResource.Close(c.DB); err != nil {
		log.Error().Stack().Err(err).Msg("failed to close postgres connection")
		errs = append(errs, err)
	} else {
		log.Info().Msg("postgres connection is closed")
	}

	if err := c.RedisClient.Close(); err != nil {
		slog.Error("Redis connection closing error", logging.ErrorAttr(err))
		errs = append(errs, err)
	} else {
		slog.Info("Redis connection is closed")
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

func mapAppRedisConfigToRedisConfig(cfg *config.RedisConfig) *cacheResource.RedisConfig {
	return &cacheResource.RedisConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
	}
}

func mapAppPostgresConfigToPostgresConfig(cfg *config.PostgresConfig) *dbResource.PostgresConfig {
	return &dbResource.PostgresConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		DBName:   cfg.DBName,
	}
}

func mapAppMinioConfigToMinioConfig(cfg *config.MinioConfig) *s3Resource.MinioConfig {
	return &s3Resource.MinioConfig{
		Host:       cfg.Host,
		Port:       cfg.Port,
		AccessKey:  cfg.AccessKey,
		SecretKey:  cfg.SecretKey,
		BucketName: cfg.BucketName,
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
