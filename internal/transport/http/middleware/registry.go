package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
)

var (
	Registry RouteMiddlewareRegistry
)

type RouteMiddlewareRegistry struct {
	Auth        gin.HandlerFunc
	UserRole    gin.HandlerFunc
	AdminRole   gin.HandlerFunc
	BannedUser  gin.HandlerFunc
	DeletedUser gin.HandlerFunc
}

func InitRegistry(jwtClient infrastructure.JWTService) {
	Registry = RouteMiddlewareRegistry{
		Auth:        JWTAuthMiddleware(jwtClient),
		UserRole:    UserRoleMiddleware(),
		AdminRole:   AdminRoleMiddleware(),
		BannedUser:  BannedUserMiddleware(),
		DeletedUser: DeletedUserMiddleware(),
	}
}
