package registry

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/minio/minio-go/v7"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/pkg/locale"
	"mandarine/pkg/logging"
	"mandarine/pkg/oauth"
	"mandarine/pkg/oauth/google"
	"mandarine/pkg/oauth/mailru"
	"mandarine/pkg/oauth/yandex"
	"mandarine/pkg/smtp"
	"mandarine/pkg/storage/cache/db_cacher"
	"mandarine/pkg/storage/cache/manager"
	cacheResource "mandarine/pkg/storage/cache/resource"
	"mandarine/pkg/storage/database"
	dbResource "mandarine/pkg/storage/database/resource"
	"mandarine/pkg/storage/s3"
	s3Resource "mandarine/pkg/storage/s3/resource"
	"mandarine/pkg/template"
	"time"
)

type Container struct {
	Config *config.Config

	Logger         *slog.Logger
	Bundle         *i18n.Bundle
	RedisClient    *redis.Client
	DB             *gorm.DB
	S3             *minio.Client
	HttpClient     *resty.Client
	OauthProviders map[string]oauth.Provider

	TemplateEngine template.Engine
	CacheManager   manager.CacheManager
	S3Client       s3.Client
	SmtpSender     smtp.Sender

	Repositories *Repositories
	Services     *Services
	Handlers     Handlers
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) MustInitialize(cfg *config.Config) {
	// Setup config
	c.Config = cfg

	// Setup logger
	c.Logger = slog.Default()

	// Setup locale
	localeConfig := mapAppLocaleConfigToLocaleConfig(&cfg.Locale)
	c.Bundle = locale.MustLoadLocales(localeConfig)

	// Setup template engine
	templateConfig := mapAppTemplateConfigToTemplateConfig(&cfg.Template)
	c.TemplateEngine = template.MustLoadTemplates(templateConfig)

	// Setup cache manager
	redisConfig := mapAppRedisConfigToRedisConfig(&cfg.Redis)
	c.RedisClient = cacheResource.MustConnectRedis(redisConfig)
	c.CacheManager = manager.NewRedisCacheManager(c.RedisClient, time.Duration(cfg.Cache.TTL)*time.Second)

	// Setup database
	postgresConfig := mapAppPostgresConfigToPostgresConfig(&cfg.Postgres)
	c.DB = dbResource.MustConnectPostgres(postgresConfig)
	err := database.UseCachePlugin(c.DB, db_cacher.NewDbCacher(c.CacheManager))
	if err != nil {
		slog.Warn("Database cache setup error", logging.ErrorAttr(err))
	}

	// Migrate database
	err = database.Migrate(dbResource.GetDSN(postgresConfig), cfg.Migrations.Path)
	if err != nil {
		slog.Warn("Database migration error", logging.ErrorAttr(err))
	}

	// Setup S3
	minioConfig := mapAppMinioConfigToMinioConfig(&cfg.Minio)
	c.S3 = s3Resource.MustConnectMinio(minioConfig)
	c.S3Client = s3.NewClient(c.S3, cfg.Minio.BucketName)

	// Setup SMTP sender
	smtpConfig := mapAppSmtpConfigToSmtpConfig(&cfg.SMTP)
	c.SmtpSender = smtp.MustNewSender(smtpConfig)

	// Setup HTTP client
	c.HttpClient = resty.New()

	// Setup OAuth clients
	c.OauthProviders = map[string]oauth.Provider{
		google.ProviderKey: google.NewOAuthGoogleProvider(cfg.OAuthClient.Google.ClientID, cfg.OAuthClient.Google.ClientSecret),
		yandex.ProviderKey: yandex.NewOAuthYandexProvider(cfg.OAuthClient.Yandex.ClientID, cfg.OAuthClient.Yandex.ClientSecret),
		mailru.ProviderKey: mailru.NewOAuthMailRuProvider(cfg.OAuthClient.MailRu.ClientID, cfg.OAuthClient.MailRu.ClientSecret),
	}

	// Setup CSR components
	c.Repositories = newGormRepositories(c.DB)
	c.Services = newServices(c)
	c.Handlers = newHandlers(c)
}

func (c *Container) Close() error {
	var errs []error

	if err := dbResource.Close(c.DB); err != nil {
		slog.Error("Postgres connection closing error", logging.ErrorAttr(err))
		errs = append(errs, err)
	} else {
		slog.Info("Postgres connection is closed")
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
