package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

var (
	ErrAccessDenied = v0.NewI18nError("access denied", "errors.access_denied")
)

func UserRoleMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleUser)
}

func AdminRoleMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleAdmin)
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	log.Debug().Msg("setup role middleware")
	logger := log.With().Str("middleware", "role").Logger()

	return func(c *gin.Context) {
		logger.Debug().Msg("check role")

		authUser, err := GetAuthUser(c)
		if err != nil {
			logger.Error().Stack().Err(err).Msg("failed to get auth user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		for _, role := range roles {
			if role == authUser.Role {
				return
			}
		}

		logger.Error().Stack().Err(ErrAccessDenied).Msg("access denied")
		_ = c.AbortWithError(http.StatusForbidden, ErrAccessDenied)
	}
}
