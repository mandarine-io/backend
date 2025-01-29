package handler

import (
	"github.com/gin-gonic/gin"
)

type APIHandler interface {
	RegisterRoutes(router *gin.Engine)
}
