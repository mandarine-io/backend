package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrBannedUser = v0.NewI18nError("banned user", "errors.banned_user")
)

// BannedUserMiddleware
//
// Use strictly after adding JWT middleware
func BannedUserMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup banned user middleware")
	logger := log.With().Str("middleware", "banned-user").Logger()

	return func(c *gin.Context) {
		logger.Debug().Msg("check banned user")

		authUser, err := GetAuthUser(c)
		if err != nil {
			logger.Error().Stack().Err(err).Msg("failed to get auth user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !authUser.IsEnabled {
			logger.Error().Stack().Err(ErrBannedUser).Msg("user is banned")
			_ = c.AbortWithError(http.StatusForbidden, ErrBannedUser)
			return
		}
	}
}
