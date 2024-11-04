package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrDeletedUser = dto.NewI18nError("deleted user", "errors.deleted_user")
)

// DeletedUserMiddleware
//
// Use strictly after adding JWT middleware
func DeletedUserMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup deleted user middleware")
	return func(c *gin.Context) {
		log.Debug().Msg("check deleted user")

		authUser, err := GetAuthUser(c)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to get auth user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if authUser.IsDeleted {
			log.Error().Stack().Err(ErrDeletedUser).Msg("user is deleted")
			_ = c.AbortWithError(http.StatusForbidden, ErrDeletedUser)
			return
		}
	}
}
