package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
)

var (
	LocalizerKey = "localizer"
)

func LocaleMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	log.Debug().Msg("setup locale middleware")
	return func(c *gin.Context) {
		localizer := i18n.NewLocalizer(bundle)

		// Header
		log.Debug().Msg("get locale request")
		if lang := c.GetHeader("Accept-Language"); lang != "" {
			log.Debug().Msg("found locale in header")
			localizer = i18n.NewLocalizer(bundle, lang)
		}

		// Query Params
		if lang, ok := c.GetQuery("lang"); ok {
			log.Debug().Msg("found locale in query params")
			localizer = i18n.NewLocalizer(bundle, lang)
		}

		c.Set(LocalizerKey, localizer)
	}
}
