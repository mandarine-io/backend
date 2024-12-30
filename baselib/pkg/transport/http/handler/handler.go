package handler

import (
	"github.com/gin-gonic/gin"
)

type ApiHandler interface {
	RegisterRoutes(router *gin.Engine)
}
