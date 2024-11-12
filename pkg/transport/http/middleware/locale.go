package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
)

var (
	LangKey      = "lang"
	LocalizerKey = "localizer"
)

func LocaleMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	log.Debug().Msg("setup locale middleware")
	return func(c *gin.Context) {
		lang := language.English.String()
		localizer := i18n.NewLocalizer(bundle, lang)

		// Header
		log.Debug().Msg("get locale request")
		if headerLang := c.GetHeader("Accept-Language"); headerLang != "" {
			log.Debug().Msg("found locale in header")

			localizer = i18n.NewLocalizer(bundle, headerLang)
			tags, _, err := language.ParseAcceptLanguage(headerLang)
			if err != nil {
				log.Warn().Err(err).Msg("failed to parse locale")
			} else if len(tags) == 0 {
				log.Warn().Err(errors.New("empty locale")).Msg("failed to parse locale")
			} else {
				lang = tags[0].String()
			}
		}

		// Query Params
		if queryLang, ok := c.GetQuery("lang"); ok {
			log.Debug().Msg("found locale in query params")

			localizer = i18n.NewLocalizer(bundle, queryLang)
			lang = queryLang
		}

		c.Set(LangKey, lang)
		c.Set(LocalizerKey, localizer)
	}
}
