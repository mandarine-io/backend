package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

var (
	LangKey      = "lang"
	LocalizerKey = "localizer"
)

func LocaleMiddleware(bundle locale.Bundle) gin.HandlerFunc {
	log.Debug().Msg("setup locale middleware")
	logger := log.With().Str("middleware", "locale").Logger()

	return func(c *gin.Context) {
		lang := language.English.String()

		// Header
		logger.Debug().Msg("get locale request")
		if headerLang := c.GetHeader("Accept-Language"); headerLang != "" {
			tags, _, err := language.ParseAcceptLanguage(headerLang)
			if err == nil && len(tags) > 0 {
				lang = tags[0].String()
			}
		}

		// Query Params
		if queryLang, ok := c.GetQuery("lang"); ok {
			lang = queryLang
		}

		localizer := bundle.NewLocalizer(lang)

		c.Set(LangKey, lang)
		c.Set(LocalizerKey, localizer)
	}
}
