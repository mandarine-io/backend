package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowMethods: []string{
				http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
				http.MethodDelete, http.MethodOptions, http.MethodHead,
			},
			AllowHeaders: []string{
				"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Frame-Options",
				"X-XSS-Protection", "Strict-Transport-Security", "Referrer-Policy", "X-Content-Type-Options",
				"Content-Security-Policy", "Permissions-Policy",
			},
			AllowCredentials: true,
			MaxAge:           24 * time.Hour,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		},
	)
}
