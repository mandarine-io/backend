package middleware

import (
	"github.com/gin-gonic/gin"
	"mandarine/pkg/rest/dto"
	"net/http"
)

var (
	ErrAccessDenied = dto.NewI18nError("access denied", "errors.access_denied")
)

type RequireRoleFactory func(...string) gin.HandlerFunc

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authUser, err := GetAuthUser(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		for _, role := range roles {
			if role == authUser.Role {
				return
			}
		}

		_ = c.AbortWithError(http.StatusForbidden, ErrAccessDenied)
	}
}
