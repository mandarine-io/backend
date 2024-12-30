package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func RecoveryMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup recovery middleware")
	return func(ctx *gin.Context) {
		defer func() {
			errRaw := recover()
			err, ok := errRaw.(error)
			if !ok {
				return
			}
			if err != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			}
		}()
	}
}
