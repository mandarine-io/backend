package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrAccessDenied = dto.NewI18nError("access denied", "errors.access_denied")
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	log.Debug().Msg("setup role middleware")
	return func(c *gin.Context) {
		log.Debug().Msg("check role")

		authUser, err := GetAuthUser(c)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to get auth user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		for _, role := range roles {
			if role == authUser.Role {
				return
			}
		}

		log.Error().Stack().Err(ErrAccessDenied).Msg("access denied")
		_ = c.AbortWithError(http.StatusForbidden, ErrAccessDenied)
	}
}
