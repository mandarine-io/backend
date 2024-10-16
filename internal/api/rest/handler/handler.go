package handler

import (
	"github.com/gin-gonic/gin"
	"mandarine/pkg/rest/middleware"
)

type ApiHandler interface {
	RegisterRoutes(router *gin.Engine, requireAuth middleware.RequireAuth, requireRole middleware.RequireRoleFactory)
}
