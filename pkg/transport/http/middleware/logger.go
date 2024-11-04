package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	requestIDCtx       = "request-id"
	requestIDHeaderKey = "X-Request-Id"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup logger middleware")
	return func(c *gin.Context) {
		// Request
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		host := c.Request.Host
		userAgent := c.Request.UserAgent()
		ip := c.ClientIP()

		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		requestID := c.GetHeader(requestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header(requestIDHeaderKey, requestID)
		}
		c.Set(requestIDCtx, requestID)

		reqEvent := log.Info().
			Str("id", requestID).
			Str("method", method).
			Str("host", host).
			Str("path", path).
			Str("query", query).
			Interface("params", params).
			Str("ip", ip).
			Str("user-agent", userAgent)
		if log.Logger.GetLevel() <= zerolog.DebugLevel {
			reqEvent.Interface("headers", c.Request.Header)
			reqEvent.Interface("body", c.Request.Body)

			c.Writer = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		}

		reqEvent.Msg("Incoming request")

		// Process
		c.Next()

		// Response
		latency := time.Since(start)
		status := c.Writer.Status()

		respEvent := log.Info().
			Str("request-id", requestID).
			Str("method", method).
			Str("host", host).
			Str("path", path).
			Str("query", query).
			Interface("params", params).
			Str("ip", ip).
			Str("user-agent", userAgent).
			Dur("latency", latency).
			Int("status", status)
		if log.Logger.GetLevel() <= zerolog.DebugLevel {
			respEvent.Interface("body", c.Writer.(*bodyLogWriter).body.String())
		}

		respEvent.Msg("Outcoming response")
	}
}
