package health

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/domain/service"
	apihandler "github.com/mandarine-io/Backend/internal/transport/http/handler"
	"github.com/rs/zerolog/log"
	"net/http"
)

type handler struct {
	svc service.HealthService
}

func NewHandler(svc service.HealthService) apihandler.ApiHandler {
	return &handler{svc: svc}
}

func (h *handler) RegisterRoutes(router *gin.Engine, _ apihandler.RouteMiddlewares) {
	log.Debug().Msg("register healthcheck routes")
	router.GET("/health", h.health)
	router.GET("/health/readiness", h.healthReadiness)
	router.GET("/health/liveness", h.healthLiveness)
}

// @Id				health
// @Summary		Health
// @Description	Request for getting health. Alias healthReadiness
// @Tags			Metrics API
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]dto.HealthOutput
// @Router			/health [get]
func (h *handler) health(c *gin.Context) {
	h.healthReadiness(c)
}

// healthReadiness godoc
//
//	@Id				healthReadiness
//	@Summary		Health readiness
//	@Description	Request for getting health readiness. In response will be status of all check (database, s3, smtp, cache, pubsub).
//	@Tags			Metrics API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]dto.HealthOutput
//	@Router			/health/readiness [get]
func (h *handler) healthReadiness(c *gin.Context) {
	log.Debug().Msg("handle health readiness")

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
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/health/liveness [get]
func (h *handler) healthLiveness(c *gin.Context) {
	log.Debug().Msg("handle health liveness")
	c.Status(http.StatusOK)
}
