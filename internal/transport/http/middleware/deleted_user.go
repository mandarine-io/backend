package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrDeletedUser = v0.NewI18nError("deleted user", "errors.deleted_user")
)

// DeletedUserMiddleware
//
// Use strictly after adding JWT middleware
func DeletedUserMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup deleted user middleware")
	logger := log.With().Str("middleware", "deleted-user").Logger()

	return func(c *gin.Context) {
		logger.Debug().Msg("check deleted user")

		authUser, err := GetAuthUser(c)
		if err != nil {
			logger.Error().Stack().Err(err).Msg("failed to get auth user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if authUser.IsDeleted {
			logger.Error().Stack().Err(ErrDeletedUser).Msg("user is deleted")
			_ = c.AbortWithError(http.StatusForbidden, ErrDeletedUser)
			return
		}
	}
}
