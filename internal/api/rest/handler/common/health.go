package common

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	healthcheck "github.com/tavsec/gin-healthcheck"
	checks2 "github.com/tavsec/gin-healthcheck/checks"
	healthcheckconfig "github.com/tavsec/gin-healthcheck/config"
	"gorm.io/gorm"
	"mandarine/pkg/smtp"
	"net/http"
)

// SetupHealthcheck godoc
//
//	@Id				Health
//	@Summary		Health
//	@Description	Request for getting health status. In response will be status of all checks (database, minio, smtp, redis).
//	@Tags			Metrics API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]dto.HealthResponse
//	@Router			/health [get]
func SetupHealthcheck(router *gin.Engine, db *gorm.DB, rdb *redis.Client, minio *minio.Client, smtp smtp.Sender) error {
	cfg := healthcheckconfig.Config{
		HealthPath:  "/health",
		Method:      "GET",
		StatusOK:    http.StatusOK,
		StatusNotOK: http.StatusServiceUnavailable,
		FailureNotification: struct {
			Threshold uint32
			Chan      chan error
		}{
			Threshold: 1,
		},
	}

	checks := []checks2.Check{
		NewGormCheck(db),
		NewMinioCheck(minio),
		NewSmtpCheck(smtp),
		checks2.NewRedisCheck(rdb),
	}
	return healthcheck.New(router, cfg, checks)
}

///////////////// MinIO /////////////////

type MinioCheck struct {
	client *minio.Client
}

func NewMinioCheck(client *minio.Client) *MinioCheck {
	return &MinioCheck{client: client}
}

func (r *MinioCheck) Pass() bool {
	return r.client.IsOnline()
}

func (r *MinioCheck) Name() string {
	return "minio"
}

///////////////// GORM /////////////////

type GormCheck struct {
	db *gorm.DB
}

func NewGormCheck(db *gorm.DB) *GormCheck {
	return &GormCheck{db: db}
}

func (r *GormCheck) Pass() bool {
	sqlDB, err := r.db.DB()
	if err != nil {
		return false
	}

	err = sqlDB.Ping()
	return err == nil
}

func (r *GormCheck) Name() string {
	return "database"
}

///////////////// SMTP /////////////////

type SmtpCheck struct {
	smtp smtp.Sender
}

func NewSmtpCheck(smtp smtp.Sender) *SmtpCheck {
	return &SmtpCheck{smtp: smtp}
}

func (r *SmtpCheck) Pass() bool {
	return r.smtp.HealthCheck()
}

func (r *SmtpCheck) Name() string {
	return "smtp"
}
