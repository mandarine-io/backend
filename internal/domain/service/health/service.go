package health

import (
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	check2 "github.com/mandarine-io/Backend/internal/domain/service/health/check"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
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
	resp := make([]dto.HealthOutput, 0)
	if s.db != nil {
		check := check2.NewGormCheck(s.db)
		resp = append(resp, dto.HealthOutput{
			Name: check.Name(),
			Pass: check.Pass(),
		})
	}
	if s.minio != nil {
		check := check2.NewMinioCheck(s.minio)
		resp = append(resp, dto.HealthOutput{
			Name: check.Name(),
			Pass: check.Pass(),
		})
	}
	if s.sender != nil {
		check := check2.NewSmtpCheck(s.sender)
		resp = append(resp, dto.HealthOutput{
			Name: check.Name(),
			Pass: check.Pass(),
		})
	}
	if s.cacheRdb != nil {
		check := check2.NewRedisCheck(s.cacheRdb)
		resp = append(resp, dto.HealthOutput{
			Name: "redis - cache",
			Pass: check.Pass(),
		})
	}
	if s.pubSubRdb != nil {
		check := check2.NewRedisCheck(s.pubSubRdb)
		resp = append(resp, dto.HealthOutput{
			Name: "redis - pub/sub",
			Pass: check.Pass(),
		})
	}

	return resp
}
