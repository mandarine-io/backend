package util

import (
	"github.com/gin-gonic/gin"
)

func ErrorWithStatus(ctx *gin.Context, status int, err error) error {
	ctx.Status(status)
	return ctx.Error(err)
}
