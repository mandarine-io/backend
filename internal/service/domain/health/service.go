package health

import (
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/domain/health/check"
	"github.com/mandarine-io/backend/pkg/model/health"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type svc struct {
	db          *gorm.DB
	cacheRDB    redis.UniversalClient
	pubsubRDB   redis.UniversalClient
	minioClient *minio.Client
	dialer      *gomail.Dialer
	logger      zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(
	db *gorm.DB,
	cacheRDB redis.UniversalClient,
	pubsubRDB redis.UniversalClient,
	minioClient *minio.Client,
	dialer *gomail.Dialer,
	opts ...Option,
) domain.HealthService {
	s := &svc{
		db:          db,
		cacheRDB:    cacheRDB,
		pubsubRDB:   pubsubRDB,
		minioClient: minioClient,
		dialer:      dialer,
		logger:      zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *svc) Health() []health.HealthOutput {
	s.logger.Info().Msg("health")

	resp := make([]health.HealthOutput, 0)

	if s.db != nil {
		c, err := check.NewGormCheck(s.db, check.WithGormLogger(s.logger))
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to GORM DB health check")
		}

		resp = append(
			resp, health.HealthOutput{
				Name: "database",
				Pass: err == nil && c.Pass(),
			},
		)
	}

	if s.dialer != nil {
		c, err := check.NewSMTPCheck(s.dialer, check.WithSMTPLogger(s.logger))
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to SMTP health check")
		}

		resp = append(
			resp, health.HealthOutput{
				Name: "smtp",
				Pass: err == nil && c.Pass(),
			},
		)
	}

	if s.minioClient != nil {
		c, err := check.NewMinioCheck(s.minioClient, check.WithMinioLogger(s.logger))
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to MinIO S3 health check")
		}

		resp = append(
			resp, health.HealthOutput{
				Name: "s3",
				Pass: err == nil && c.Pass(),
			},
		)
	}

	if s.cacheRDB != nil {
		c, err := check.NewRedisCheck(s.cacheRDB, check.WithRedisLogger(s.logger))
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to Redis cache health check")
		}

		resp = append(
			resp, health.HealthOutput{
				Name: "cache",
				Pass: err == nil && c.Pass(),
			},
		)
	}

	if s.pubsubRDB != nil {
		c, err := check.NewRedisCheck(s.pubsubRDB, check.WithRedisLogger(s.logger))
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to Redis pub/sub health check")
		}

		resp = append(
			resp, health.HealthOutput{
				Name: "pub/sub",
				Pass: err == nil && c.Pass(),
			},
		)
	}

	return resp
}
