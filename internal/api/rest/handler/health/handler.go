package health

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/api/rest/handler"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	healthcheck "github.com/tavsec/gin-healthcheck"
	checks2 "github.com/tavsec/gin-healthcheck/checks"
	healthcheckconfig "github.com/tavsec/gin-healthcheck/config"
	"gorm.io/gorm"
	"net/http"
)

type Response struct {
	Name string `json:"name"`
	Pass bool   `json:"pass"`
}

type Handler struct {
	db    *gorm.DB
	rdb   *redis.Client
	minio *minio.Client
	smtp  smtp.Sender
}

func NewHandler(db *gorm.DB, rdb *redis.Client, minio *minio.Client, smtp smtp.Sender) *Handler {
	return &Handler{db: db, rdb: rdb, minio: minio, smtp: smtp}
}

// RegisterRoutes godoc
//
//	@Id				Health
//	@Summary		Health
//	@Description	Request for getting health status. In response will be status of all checks (database, minio, smtp, redis).
//	@Tags			Metrics API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]Response
//	@Router			/health [get]
func (h *Handler) RegisterRoutes(router *gin.Engine, _ handler.RouteMiddlewares) {
	log.Debug().Msg("register healthcheck routes")

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
		NewGormCheck(h.db),
		NewMinioCheck(h.minio),
		NewSmtpCheck(h.smtp),
		checks2.NewRedisCheck(h.rdb),
	}

	err := healthcheck.New(router, cfg, checks)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to setup healthcheck")
	}
}

///////////////// MinIO /////////////////

type MinioCheck struct {
	client *minio.Client
}

func NewMinioCheck(client *minio.Client) *MinioCheck {
	return &MinioCheck{client: client}
}

func (r *MinioCheck) Pass() bool {
	log.Debug().Msg("check minio connection")
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
	log.Debug().Msg("check gorm connection")

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
	log.Debug().Msg("check smtp connection")
	return r.smtp.HealthCheck()
}

func (r *SmtpCheck) Name() string {
	return "smtp"
}
