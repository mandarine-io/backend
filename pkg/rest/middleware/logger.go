package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"mandarine/pkg/logging"
	"net/http"
	"time"
)

const (
	requestIDCtx       = "request-id"
	requestIDHeaderKey = "X-Request-Id"
)

func LoggerMiddleware() gin.HandlerFunc {
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

		requestAttributes := []slog.Attr{
			slog.String("id", requestID),
			slog.String("method", method),
			slog.String("host", host),
			slog.String("path", path),
			slog.String("query", query),
			slog.Any("params", params),
			slog.String("ip", ip),
			slog.String("user-agent", userAgent),
		}
		attributes := []slog.Attr{
			{
				Key:   "request",
				Value: slog.GroupValue(requestAttributes...),
			},
		}
		slog.LogAttrs(c.Request.Context(), slog.LevelInfo, "Incoming request", attributes...)

		// Process
		c.Next()

		// Response
		latency := time.Since(start)
		status := c.Writer.Status()

		responseAttributes := []slog.Attr{
			slog.String("request-id", requestID),
			slog.Duration("latency", latency),
			slog.Int("status", status),
		}

		attributes = []slog.Attr{
			{
				Key:   "response",
				Value: slog.GroupValue(responseAttributes...),
			},
		}

		level := slog.LevelInfo
		msg := "Outcoming response"
		if status >= http.StatusBadRequest {
			level = slog.LevelError
			attributes = append(attributes, logging.ErrorStringAttr(c.Errors.String()))
		}

		slog.LogAttrs(c.Request.Context(), level, msg, attributes...)
	}
}
