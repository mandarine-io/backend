package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	LocalizerKey = "localizer"
)

func LocaleMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		localizer := i18n.NewLocalizer(bundle)

		// Header
		if lang := c.GetHeader("Accept-Language"); lang != "" {
			localizer = i18n.NewLocalizer(bundle, lang)
		}

		// Query Params
		if lang, ok := c.GetQuery("lang"); ok {
			localizer = i18n.NewLocalizer(bundle, lang)
		}

		c.Set(LocalizerKey, localizer)
	}
}
