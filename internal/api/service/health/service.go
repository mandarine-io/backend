package health

import (
	check2 "github.com/mandarine-io/Backend/internal/api/service/health/check"
	"github.com/mandarine-io/Backend/internal/api/service/health/dto"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	checks2 "github.com/tavsec/gin-healthcheck/checks"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	cacheRdb  *redis.Client
	pubSubRdb *redis.Client
	s3        *minio.Client
	sender    smtp.Sender
}

func NewService(db *gorm.DB, cacheRdb *redis.Client, pubSubRdb *redis.Client, s3 *minio.Client, sender smtp.Sender) *Service {
	return &Service{
		db:        db,
		cacheRdb:  cacheRdb,
		pubSubRdb: pubSubRdb,
		s3:        s3,
		sender:    sender,
	}
}

func (s *Service) Health() []dto.HealthOutput {
	checks := make([]checks2.Check, 0)
	if s.db != nil {
		checks = append(checks, check2.NewGormCheck(s.db))
	}
	if s.s3 != nil {
		checks = append(checks, check2.NewMinioCheck(s.s3))
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
