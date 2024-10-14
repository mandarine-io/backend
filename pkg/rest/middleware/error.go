package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"mandarine/pkg/locale"
	"mandarine/pkg/rest/dto"
	"strings"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()
		if lastErr == nil {
			return
		}
		err := lastErr.Err

		localizerAny, _ := c.Get("localizer")
		localizer := localizerAny.(*i18n.Localizer)

		status := c.Writer.Status()
		var errorResponse dto.ErrorResponse

		var validationErrors validator.ValidationErrors
		var i18nError dto.I18nError
		switch {
		case errors.As(err, &validationErrors):
			var validErr validator.ValidationErrors
			errors.As(err, &validErr)
			errorResponse = dto.NewErrorResponse(
				convertValidationErrorsToString(validErr, localizer), status, c.Request.URL.Path,
			)
		case errors.As(err, &i18nError):
			var i18nErr dto.I18nError
			errors.As(err, &i18nErr)
			errorResponse = dto.NewErrorResponse(
				locale.LocalizeWithArgs(localizer, i18nErr.Tag(), i18nErr.Args()), status, c.Request.URL.Path,
			)
		default:
			errorResponse = dto.NewErrorResponse(
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
