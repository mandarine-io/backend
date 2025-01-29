package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
)

const (
	subsystem   = "backend"
	resultKey   = "requests_total"
	durationKey = "requests_duration_seconds"
)

var (
	requestResult = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      resultKey,
			Help:      "Number of HTTP requests, partitioned by status code, method, and path.",
		},
		[]string{"path", "method"},
	)

	requestLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: subsystem,
			Name:      durationKey,
			Help:      "Request latency in seconds. Broken down by verb, and path.",
		},
		[]string{"path", "method"},
	)

	once sync.Once
)

type MetricsAdapter interface {
	IncrementRequestTotal(r *http.Request)
	UpdateRequestLatency(r *http.Request, duration float64)
}

type defaultMetricAdapter struct {
	logger zerolog.Logger
}

type Option func(*defaultMetricAdapter)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *defaultMetricAdapter) {
		h.logger = logger
	}
}

func NewMetricAdapter(opts ...Option) MetricsAdapter {
	d := &defaultMetricAdapter{
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(d)
	}

	once.Do(
		func() {
			d.logger.Info().Msg("register custom metrics")
			prometheus.MustRegister(requestResult)
			prometheus.MustRegister(requestLatency)
		},
	)

	return d
}

func (d *defaultMetricAdapter) IncrementRequestTotal(r *http.Request) {
	path, method := d.extractPathAndMethod(r)

	d.logger.Debug().Msgf("increment request total, path: %s, method: %s", path, method)
	requestResult.WithLabelValues(path, method).Add(1)
}

func (d *defaultMetricAdapter) UpdateRequestLatency(r *http.Request, duration float64) {
	path, method := d.extractPathAndMethod(r)

	d.logger.Debug().Msgf("update request latency, path: %s, method: %s, duration: %f", path, method, duration)
	requestLatency.WithLabelValues(path, method).Set(duration)
}

func (d *defaultMetricAdapter) extractPathAndMethod(req *http.Request) (string, string) {
	path := "none"
	method := "none"

	if req.URL.Path != "" {
		path = req.URL.Path
	}

	if req.Method != "" {
		method = req.Method
	}

	return path, method
}
