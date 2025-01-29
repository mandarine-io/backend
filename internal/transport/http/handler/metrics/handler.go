package metrics

import (
	"github.com/gin-gonic/gin"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type handler struct {
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(opts ...Option) apihandler.APIHandler {
	h := &handler{
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	h.logger.Debug().Msg("register metrics routes")
	router.GET("/metrics/prometheus", h.getMetrics)
}

// health godoc
//
//	@Id				metrics
//	@Summary		Metrics in Prometheus format
//	@Description	Request for getting Prometheus metrics
//	@Tags			Metrics API
//	@Produce		text/plain; charset=utf-8
//	@Success		200	{object}	string
//	@Router			/metrics/prometheus [get]
func (h *handler) getMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
