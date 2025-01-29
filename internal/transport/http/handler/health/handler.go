package health

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/service/domain"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	_ "github.com/mandarine-io/backend/pkg/model/health"
	_ "github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"net/http"
)

type handler struct {
	logger zerolog.Logger
	svc    domain.HealthService
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.HealthService, opts ...Option) apihandler.APIHandler {
	h := &handler{
		svc:    svc,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	h.logger.Debug().Msg("register health routes")

	healthRouter := router.Group("/health")
	{
		healthRouter.GET("", h.health)
		healthRouter.GET("/readiness", h.healthReadiness)
		healthRouter.GET("/liveness", h.healthLiveness)
	}
}

// health godoc
//
//	@Id				health
//	@Summary		Health
//	@Description	Request for getting health. Alias healthReadiness
//	@Tags			Metrics API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200	{object}	[]health.HealthOutput
//	@Router			/health [get]
func (h *handler) health(c *gin.Context) {
	h.healthReadiness(c)
}

// healthReadiness godoc
//
//	@Id				healthReadiness
//	@Summary		Health readiness
//	@Description	Request for getting health readiness. In response will be status of all check (database, s3, smtp, cache, pubsub).
//	@Tags			Metrics API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200	{object}	[]health.HealthOutput
//	@Failure		503	{object}	[]v0.ErrorOutput
//	@Router			/health/readiness [get]
func (h *handler) healthReadiness(c *gin.Context) {
	h.logger.Debug().Msg("handle health readiness")

	resp := h.svc.Health()

	for _, v := range resp {
		if !v.Pass {
			c.JSON(http.StatusServiceUnavailable, resp)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

// healthLiveness godoc
//
//	@Id				healthLiveness
//	@Summary		Health liveness
//	@Description	Request for getting health liveness.
//	@Tags			Metrics API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200
//	@Router			/health/liveness [get]
func (h *handler) healthLiveness(c *gin.Context) {
	h.logger.Debug().Msg("handle health liveness")
	c.Status(http.StatusOK)
}
