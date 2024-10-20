package locale

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
	"mandarine/internal/api/helper/file"
	"mandarine/pkg/logging"
	"os"
	"path"
	"path/filepath"
)

type Config struct {
	Path     string
	Language string
}

func MustLoadLocales(cfg *Config) *i18n.Bundle {
	// Parse default language tag
	tag, err := language.Parse(cfg.Language)
	if err != nil {
		slog.Warn("Default language parsing error", logging.ErrorAttr(err))
		tag = language.English
	}
	slog.Info("Default language: " + tag.String())

	// Read locale files
	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := file.GetFilesFromDir(cfg.Path)
	if err != nil {
		slog.Error("Locales set up error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	for _, f := range files {
		filePath := path.Join(cfg.Path, f)
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			slog.Error("Locales set up error", logging.ErrorAttr(err))
			os.Exit(1)
		}

		slog.Info("Read translation file: " + absFilePath)
		_, err = bundle.LoadMessageFile(absFilePath)
		if err != nil {
			slog.Error("Locales set up error", logging.ErrorAttr(err))
			os.Exit(1)
		}
	}

	return bundle
}

func Localize(localizer *i18n.Localizer, tag string) string {
	message, err := localizer.Localize(
		&i18n.LocalizeConfig{
			MessageID: tag,
		},
	)
	if err != nil {
		slog.Warn("Localize error", logging.ErrorAttr(err))
		message = tag
	}
	return message
}

func LocalizeWithArgs(localizer *i18n.Localizer, tag string, args interface{}) string {
	message, err := localizer.Localize(
		&i18n.LocalizeConfig{
			MessageID:    tag,
			TemplateData: args,
		},
	)
	if err != nil {
		slog.Warn("Localize error", logging.ErrorAttr(err))
		message = tag
	}
	return message
}
