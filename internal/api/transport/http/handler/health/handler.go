package health

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/api/service/health"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	svc *health.Service
}

func NewHandler(svc *health.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes godoc
//
//	@Id				Health
//	@Summary		Health
//	@Description	Request for getting health status. In response will be status of all check (database, minio, smtp, redis).
//	@Tags			Metrics API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]dto.HealthOutput
//	@Router			/health [get]
func (h *Handler) RegisterRoutes(router *gin.Engine, _ handler.RouteMiddlewares) {
	log.Debug().Msg("register healthcheck routes")
	router.GET("/health", h.Health)
}

func (h *Handler) Health(c *gin.Context) {
	log.Debug().Msg("handle health")

	resp := h.svc.Health()

	healthy := true
	for _, v := range resp {
		if !v.Pass {
			healthy = false
			break
		}
	}

	if healthy {
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusServiceUnavailable, resp)
	}
}
