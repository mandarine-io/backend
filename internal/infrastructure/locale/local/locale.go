package local

import (
	"encoding/json"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/internal/util/file"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	syspath "path"
)

type bundle struct {
	bundle      *i18n.Bundle
	logger      zerolog.Logger
	defaultLang language.Tag
}

type Option func(bundle *bundle) error

func WithLogger(logger zerolog.Logger) Option {
	return func(b *bundle) error {
		b.logger = logger
		return nil
	}
}

func WithDefaultLang(lang string) Option {
	return func(b *bundle) error {
		tag, err := language.Parse(lang)
		if err != nil {
			return fmt.Errorf("failed to parse default language: %w", err)
		}

		b.defaultLang = tag
		return nil
	}
}

func NewBundle(path string, opts ...Option) (locale.Bundle, error) {
	b := &bundle{
		logger:      zerolog.Nop(),
		defaultLang: language.English,
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	b.bundle = i18n.NewBundle(b.defaultLang)
	b.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	b.bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	b.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	b.logger.Debug().Msg("read locale files")

	files, err := file.GetFilesFromDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get locale files: %w", err)
	}

	for _, f := range files {
		filePath := syspath.Join(path, f)

		b.logger.Debug().Msgf("load translation file: %s", filePath)
		_, err = b.bundle.LoadMessageFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load locale file: %w", err)
		}
	}

	return b, nil
}

type localLocalizer struct {
	localizer *i18n.Localizer
	logger    zerolog.Logger
}

func (b *bundle) NewLocalizer(lang string) locale.Localizer {
	return &localLocalizer{
		localizer: i18n.NewLocalizer(b.bundle, lang),
		logger:    b.logger,
	}
}

func (l *localLocalizer) Localize(tag string, args any, pluralCount int) string {
	l.logger.Debug().Msgf("localize message by tag: %s", tag)

	cfg := &i18n.LocalizeConfig{
		MessageID:    tag,
		TemplateData: args,
	}

	if pluralCount >= 0 {
		cfg.PluralCount = pluralCount
	}

	message, err := l.localizer.Localize(cfg)
	if err != nil {
		l.logger.Warn().Err(err).Msg("failed to localize message by tag, use tag as fallback")
		message = tag
	}

	return message
}
