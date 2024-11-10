package handler

import (
	"github.com/gin-gonic/gin"
)

type RouteMiddlewares struct {
	Auth        gin.HandlerFunc
	RoleFactory func(...string) gin.HandlerFunc
	BannedUser  gin.HandlerFunc
	DeletedUser gin.HandlerFunc
}

type ApiHandler interface {
	RegisterRoutes(router *gin.Engine, middlewares RouteMiddlewares)
}
