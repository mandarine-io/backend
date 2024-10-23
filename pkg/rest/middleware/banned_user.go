package middleware

import (
	"github.com/gin-gonic/gin"
	"mandarine/pkg/rest/dto"
	"net/http"
)

var (
	ErrBannedUser = dto.NewI18nError("banned user", "errors.banned_user")
)

// BannedUserMiddleware
//
// Use strictly after adding JWT middleware
func BannedUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser, err := GetAuthUser(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if !authUser.IsEnabled {
			_ = c.AbortWithError(http.StatusForbidden, ErrBannedUser)
			return
		}
	}
}
