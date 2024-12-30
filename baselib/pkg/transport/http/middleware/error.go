package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mandarine-io/baselib/pkg/locale"
	"github.com/mandarine-io/baselib/pkg/transport/http/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

func ErrorMiddleware() gin.HandlerFunc {
	log.Debug().Msg("setup register error middleware")
	return func(c *gin.Context) {
		c.Next()

		// get the last error
		log.Debug().Msg("get the last error")
		lastErr := c.Errors.Last()
		if lastErr == nil {
			log.Debug().Msg("no error found")
			return
		}

		log.Debug().Msg("found the last error")
		err := lastErr.Err

		// get localizer
		log.Debug().Msg("get localizer")
		localizerAny, _ := c.Get("localizer")
		localizer := localizerAny.(*i18n.Localizer)

		// build error response
		log.Debug().Msg("build error response")
		status := c.Writer.Status()
		var errorResponse model.ErrorResponse

		var (
			validErrs validator.ValidationErrors
			i18nErr   model.I18nError
		)
		switch {
		case errors.As(err, &validErrs):
			errorResponse = model.NewErrorResponse(
				convertValidationErrorsToString(validErrs, localizer), status, c.Request.URL.Path,
			)
		case errors.As(err, &i18nErr):
			errorResponse = model.NewErrorResponse(
				locale.LocalizeWithArgs(localizer, i18nErr.Tag(), i18nErr.Args()), status, c.Request.URL.Path,
			)
		case errors.Is(err, io.EOF):
			errorResponse = model.NewErrorResponse(
				locale.Localize(localizer, "errors.empty_body"), status, c.Request.URL.Path,
			)
		default:
			errorResponse = model.NewErrorResponse(
				locale.Localize(localizer, "errors.internal_error"), status, c.Request.URL.Path,
			)
		}

		c.JSON(status, errorResponse)
	}
}

func convertValidationErrorsToString(validErrs validator.ValidationErrors, localizer *i18n.Localizer) string {
	errStrs := make([]string, len(validErrs))
	for i, validErr := range validErrs {
		i18nTag := "errors.validation." + validErr.Tag()
		message := locale.LocalizeWithArgs(localizer, i18nTag, map[string]string{"param": validErr.Param()})
		errStrs[i] = fmt.Sprintf("%s: %s", validErr.StructField(), message)
	}

	return strings.Join(errStrs, "; ")
}
