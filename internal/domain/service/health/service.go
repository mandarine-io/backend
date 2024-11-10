package health

import (
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	check2 "github.com/mandarine-io/Backend/internal/domain/service/health/check"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	checks2 "github.com/tavsec/gin-healthcheck/checks"
	"gorm.io/gorm"
)

type svc struct {
	db        *gorm.DB
	cacheRdb  redis.UniversalClient
	pubSubRdb redis.UniversalClient
	minio     *minio.Client
	sender    smtp.Sender
}

func NewService(db *gorm.DB, cacheRdb redis.UniversalClient, pubSubRdb redis.UniversalClient, s3 *minio.Client, sender smtp.Sender) service.HealthService {
	return &svc{
		db:        db,
		cacheRdb:  cacheRdb,
		pubSubRdb: pubSubRdb,
		minio:     s3,
		sender:    sender,
	}
}

func (s *svc) Health() []dto.HealthOutput {
	checks := make([]checks2.Check, 0)
	if s.db != nil {
		checks = append(checks, check2.NewGormCheck(s.db))
	}
	if s.minio != nil {
		checks = append(checks, check2.NewMinioCheck(s.minio))
	}
	if s.sender != nil {
		checks = append(checks, check2.NewSmtpCheck(s.sender))
	}
	if s.cacheRdb != nil {
		checks = append(checks, check2.NewRedisCheck(s.cacheRdb))
	}
	if s.pubSubRdb != nil {
		checks = append(checks, check2.NewRedisCheck(s.pubSubRdb))
	}

	resp := make([]dto.HealthOutput, 0)
	for _, check := range checks {
		resp = append(resp, dto.HealthOutput{
			Name: check.Name(),
			Pass: check.Pass(),
		})
	}

	return resp
}
