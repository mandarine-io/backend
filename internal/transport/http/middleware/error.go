package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

func ErrorMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup error middleware")
	logger := log.With().Str("middleware", "error").Logger()

	return func(c *gin.Context) {
		c.Next()

		// get the last error
		logger.Debug().Msg("get the last error")
		lastErr := c.Errors.Last()
		if lastErr == nil {
			logger.Debug().Msg("no error found")
			return
		}

		logger.Debug().Msg("found the last error")
		err := lastErr.Err
		logger.Error().Stack().Err(err).Msg("failed to handle request")

		// get localizer
		logger.Debug().Msg("get localizer")
		localizerAny, _ := c.Get(LocalizerKey)
		localizer := localizerAny.(locale.Localizer)

		// build error response
		logger.Debug().Msg("build error response")
		status := c.Writer.Status()
		var errorOutput v0.ErrorOutput

		var (
			validErrs validator.ValidationErrors
			i18nErr   v0.I18nError
			syntaxErr *json.SyntaxError
		)
		switch {
		case errors.As(err, &validErrs):
			errorOutput = v0.NewErrorOutput(
				convertValidationErrorsToString(validErrs, localizer), status, c.Request.URL.Path,
			)
		case errors.As(err, &i18nErr):
			errorOutput = v0.NewErrorOutput(
				localizer.Localize(i18nErr.Tag(), i18nErr.Args(), -1), status, c.Request.URL.Path,
			)
		case errors.Is(err, io.EOF):
			errorOutput = v0.NewErrorOutput(
				localizer.Localize("errors.empty_body", nil, -1), status, c.Request.URL.Path,
			)
		case errors.As(err, &syntaxErr):
			errorOutput = v0.NewErrorOutput(
				localizer.Localize("errors.syntax_error", nil, -1), status, c.Request.URL.Path,
			)
		default:
			errorOutput = v0.NewErrorOutput(
				localizer.Localize("errors.internal_error", nil, -1), status, c.Request.URL.Path,
			)
		}

		c.JSON(status, errorOutput)
	}
}

func convertValidationErrorsToString(validErrs validator.ValidationErrors, localizer locale.Localizer) string {
	errStrs := make([]string, len(validErrs))
	for i, validErr := range validErrs {
		i18nTag := "errors.validation." + validErr.Tag()
		message := localizer.Localize(i18nTag, map[string]string{"param": validErr.Param()}, -1)
		errStrs[i] = fmt.Sprintf("%s: %s", validErr.StructField(), message)
	}

	return strings.Join(errStrs, "; ")
}
