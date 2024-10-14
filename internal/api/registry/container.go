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
	"mandarine/internal/api/persistence/repo"
	gorm2 "mandarine/internal/api/persistence/repo/gorm"
	"mandarine/internal/api/rest/handler"
	"mandarine/internal/api/service"
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

	Repositories *repo.Repositories
	Services     *service.Services
	Handlers     *handler.Handlers
}

func MustNewContainer(cfg *config.Config) *Container {
	// Setup locale
	localeConfig := mapAppLocaleConfigToLocaleConfig(&cfg.Locale)
	bundle := locale.MustLoadLocales(localeConfig)

	// Setup template engine
	templateConfig := mapAppTemplateConfigToTemplateConfig(&cfg.Template)
	templateEngine := template.MustLoadTemplates(templateConfig)

	// Setup cache manager
	redisConfig := mapAppRedisConfigToRedisConfig(&cfg.Redis)
	redisClient := cacheResource.MustConnectRedis(redisConfig)
	cacheManager := manager.NewRedisCacheManager(redisClient, time.Duration(cfg.Cache.TTL)*time.Second)

	// Setup database
	postgresConfig := mapAppPostgresConfigToPostgresConfig(&cfg.Postgres)
	postgresDB := dbResource.MustConnectPostgres(postgresConfig)
	err := database.UseCachePlugin(postgresDB, db_cacher.NewDbCacher(cacheManager))
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
	minio := s3Resource.MustConnectMinio(minioConfig)
	s3Client := s3.NewClient(minio, cfg.Minio.BucketName)

	// Setup SMTP sender
	smtpConfig := mapAppSmtpConfigToSmtpConfig(&cfg.SMTP)
	smtpSender := smtp.MustNewSender(smtpConfig)

	// Setup HTTP client
	httpClient := resty.New()

	// Setup OAuth clients
	oauthProviders := map[string]oauth.Provider{
		google.ProviderKey: google.NewOAuthGoogleProvider(cfg.OAuthClient.Google.ClientID, cfg.OAuthClient.Google.ClientSecret),
		yandex.ProviderKey: yandex.NewOAuthYandexProvider(cfg.OAuthClient.Yandex.ClientID, cfg.OAuthClient.Yandex.ClientSecret),
		mailru.ProviderKey: mailru.NewOAuthMailRuProvider(cfg.OAuthClient.MailRu.ClientID, cfg.OAuthClient.MailRu.ClientSecret),
	}

	// Setup CSR components
	repos := gorm2.NewRepositories(postgresDB)
	services := service.NewServices(
		repos,
		cacheManager,
		smtpSender,
		templateEngine,
		oauthProviders,
		cfg,
	)
	handlers := handler.NewHandlers(services, cfg)

	return &Container{
		Config: cfg,

		Logger:         slog.Default(),
		Bundle:         bundle,
		RedisClient:    redisClient,
		DB:             postgresDB,
		S3:             minio,
		HttpClient:     httpClient,
		OauthProviders: oauthProviders,

		TemplateEngine: templateEngine,
		CacheManager:   cacheManager,
		S3Client:       s3Client,
		SmtpSender:     smtpSender,

		Repositories: repos,
		Services:     services,
		Handlers:     handlers,
	}
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
