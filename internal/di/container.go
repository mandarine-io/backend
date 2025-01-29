package di

import (
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/mandarine-io/backend/internal/infrastructure/smtp"
	"github.com/mandarine-io/backend/internal/infrastructure/template"
	"github.com/mandarine-io/backend/internal/infrastructure/websocket"
	"github.com/mandarine-io/backend/internal/observability"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/scheduler"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/mandarine-io/backend/third_party/oauth"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type Initializer func() error

type Finalizer func() error

type Infrastructure struct {
	DB          *gorm.DB
	CacheRDB    redis.UniversalClient
	PubSubRDB   redis.UniversalClient
	MinioClient *minio.Client
	SMTPDialer  *gomail.Dialer

	LocaleBundle   locale.Bundle
	TemplateEngine template.Engine
	CacheManager   cache.Manager
	S3Manager      s3.Manager
	SMTPSender     smtp.Sender
	PubSubAgent    pubsub.Agent
	Scheduler      *scheduler.Scheduler
	WSPool         *websocket.Pool
}

type Repositories struct {
	MasterProfile repo.MasterProfileRepository
	MasterService repo.MasterServiceRepository
	User          repo.UserRepository
}

type InfrastructureServices struct {
	JWT infrastructure.JWTService
	OTP infrastructure.OTPService
}

type DomainServices struct {
	Account       domain.AccountService
	Auth          domain.AuthService
	Health        domain.HealthService
	Geocoding     domain.GeocodingService
	MasterProfile domain.MasterProfileService
	MasterService domain.MasterServiceService
	Resource      domain.ResourceService
	Websocket     domain.WebsocketService
}

type ThirdParties struct {
	Geocoding map[string]geocoding.Provider
	OAuth     map[string]oauth.Provider
}

type Handlers []handler.APIHandler

type Container struct {
	Config             config.Config
	Logger             zerolog.Logger
	Infrastructure     Infrastructure
	Metrics            observability.MetricsAdapter
	Repos              Repositories
	InfrastructureSVCs InfrastructureServices
	DomainSVCs         DomainServices
	Handlers           Handlers
	ThirdParties       ThirdParties

	initializers []Initializer
	finalizers   []Finalizer
}

func NewContainer(cfg config.Config) *Container {
	return &Container{
		Logger: log.Logger,
		Config: cfg,
	}
}

func (c *Container) RegisterInitializers(initializers ...Initializer) {
	c.initializers = append(c.initializers, initializers...)
}

func (c *Container) RegisterFinalizers(finalizers ...Finalizer) {
	c.finalizers = append(c.finalizers, finalizers...)
}

func (c *Container) Initialize() error {
	for _, initialize := range c.initializers {
		err := initialize()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) Finalize() error {
	for _, finalize := range c.finalizers {
		err := finalize()
		if err != nil {
			return err
		}
	}

	return nil
}
