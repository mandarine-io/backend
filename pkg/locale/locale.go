package locale

import (
	"encoding/json"
	"github.com/mandarine-io/Backend/pkg/helper/file"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"path"
	"path/filepath"
)

type Config struct {
	Path     string
	Language string
}

func MustLoadLocales(cfg *Config) *i18n.Bundle {
	// Parse default language tag
	log.Debug().Msg("parse default language")
	tag, err := language.Parse(cfg.Language)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse default language")
		tag = language.English
	}
	log.Info().Msgf("set default language: %s", tag)

	// Read locale files
	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	log.Debug().Msg("read locale files")
	files, err := file.GetFilesFromDir(cfg.Path)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get locale files")
	}

	for _, f := range files {
		filePath := path.Join(cfg.Path, f)
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			log.Fatal().Stack().Err(err).Msgf("failed to get absolute path of locale file: %s", filePath)
		}

		log.Info().Msgf("load translation file: %s", absFilePath)
		_, err = bundle.LoadMessageFile(absFilePath)
		if err != nil {
			log.Fatal().Stack().Err(err).Msg("failed to load locale file")
		}
	}

	return bundle
}

func Localize(localizer *i18n.Localizer, tag string) string {
	log.Debug().Msgf("localize message by tag: %s", tag)
	message, err := localizer.Localize(
		&i18n.LocalizeConfig{
			MessageID: tag,
		},
	)
	if err != nil {
		log.Warn().Err(err).Msg("failed to localize message by tag, use tag as fallback")
		message = tag
	}
	return message
}

func LocalizeWithArgs(localizer *i18n.Localizer, tag string, args interface{}) string {
	log.Debug().Msgf("localize message by tag: %s", tag)
	message, err := localizer.Localize(
		&i18n.LocalizeConfig{
			MessageID:    tag,
			TemplateData: args,
		},
	)
	if err != nil {
		log.Warn().Err(err).Msg("failed to localize message by tag, use tag as fallback")
		message = tag
	}
	return message
}
