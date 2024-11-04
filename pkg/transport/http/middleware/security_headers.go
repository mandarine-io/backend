package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SecurityHeadersMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup security headers middleware")
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header(
			"Content-Security-Policy",
			"default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';",
		)
		c.Header(
			"Permissions-Policy",
			"geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()",
		)
	}
}
