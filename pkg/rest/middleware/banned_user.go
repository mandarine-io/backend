package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrBannedUser = dto.NewI18nError("banned user", "errors.banned_user")
)

// BannedUserMiddleware
//
// Use strictly after adding JWT middleware
func BannedUserMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup banned user middleware")
	return func(c *gin.Context) {
		log.Debug().Msg("check banned user")

		authUser, err := GetAuthUser(c)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to get auth user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !authUser.IsEnabled {
			log.Error().Stack().Err(ErrBannedUser).Msg("user is banned")
			_ = c.AbortWithError(http.StatusForbidden, ErrBannedUser)
			return
		}
	}
}
