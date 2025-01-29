package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware() gin.HandlerFunc {
	log.Debug().Str("middleware", "recovery").Msg("setup recovery middleware")
	logger := log.With().Str("middleware", "recovery").Logger()

	return func(ctx *gin.Context) {
		defer func() {
			errRaw := recover()

			// get error
			err, ok := errRaw.(error)
			if !ok || err == nil {
				return
			}

			logger.Error().Msgf("catch panic: %s\n%s", err, string(debug.Stack()))

			// get localizer
			logger.Debug().Msg("get localizer")
			localizerAny, _ := ctx.Get(LocalizerKey)
			localizer := localizerAny.(locale.Localizer)

			// build error output
			errorOutput := v0.NewErrorOutput(
				localizer.Localize("errors.internal_error", nil, -1),
				http.StatusInternalServerError,
				ctx.Request.URL.Path,
			)

			ctx.JSON(http.StatusInternalServerError, errorOutput)
		}()
		ctx.Next()
	}
}
