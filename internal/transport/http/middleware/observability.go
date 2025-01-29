package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/observability"
	"github.com/rs/zerolog/log"
	"regexp"
	"time"
)

func ObservabilityMiddleware(adapter observability.MetricsAdapter, ignorePathRegexpStrs ...string) gin.HandlerFunc {
	log.Debug().Msg("setup observability middleware")
	logger := log.With().Str("middleware", "observability").Logger()

	ignorePathRegexps := make([]*regexp.Regexp, 0)
	for _, regexpStr := range ignorePathRegexpStrs {
		logger.Debug().Msgf("compile ingnore regexp: %s", regexpStr)
		compiledRegexp, err := regexp.Compile(regexpStr)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to compiled regexp")
			continue
		}

		ignorePathRegexps = append(ignorePathRegexps, compiledRegexp)
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		for _, pathRegexp := range ignorePathRegexps {
			match := pathRegexp.MatchString(path)
			if match {
				return
			}
		}

		c.Next()
		latency := time.Since(start)

		adapter.IncrementRequestTotal(c.Request)
		adapter.UpdateRequestLatency(c.Request, latency.Seconds())
	}
}
