package common

import (
	"github.com/gin-gonic/gin"
	"mandarine/pkg/rest/dto"
	"net/http"
)

func NoMethod(ctx *gin.Context) {
	_ = ctx.AbortWithError(http.StatusMethodNotAllowed, dto.NewI18nError("method not allowed", "errors.method_not_allowed"))
}

func NoRoute(ctx *gin.Context) {
	_ = ctx.AbortWithError(http.StatusNotFound, dto.NewI18nError("route not found", "errors.route_not_found"))
}
