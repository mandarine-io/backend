package middleware

import (
	"github.com/gin-gonic/gin"
	"mandarine/pkg/rest/dto"
	"net/http"
)

var (
	ErrDeletedUser = dto.NewI18nError("deleted user", "errors.deleted_user")
)

// DeletedUserMiddleware
//
// Use strictly after adding JWT middleware
func DeletedUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser, err := GetAuthUser(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if authUser.IsDeleted {
			_ = c.AbortWithError(http.StatusForbidden, ErrDeletedUser)
			return
		}
	}
}
