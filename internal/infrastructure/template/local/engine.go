package local

import (
	"bytes"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/template"
	"github.com/mandarine-io/backend/internal/util/file"
	"github.com/rs/zerolog"
	htmltemplate "html/template"
	syspath "path"
	texttemplate "text/template"
)

type Option func(*engine) error

func WithLogger(logger zerolog.Logger) Option {
	return func(e *engine) error {
		e.logger = logger
		return nil
	}
}

type engine struct {
	templates map[string]string
	logger    zerolog.Logger
}

func NewEngine(path string, opts ...Option) (template.Engine, error) {
	e := &engine{
		templates: make(map[string]string),
		logger:    zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	e.logger.Debug().Msgf("load templates from %s", path)

	files, err := file.GetFilesFromDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get template files: %w", err)
	}

	for _, f := range files {
		filePath := syspath.Join(path, f)

		e.logger.Debug().Msgf("load template file: %s", filePath)
		tmplName := file.GetFileNameWithoutExt(f)
		e.templates[tmplName] = filePath
	}

	return e, nil
}

func (e *engine) RenderHTML(name string, args any) (string, error) {
	e.logger.Debug().Msgf("search template: %s", name)
	tmplPath, ok := e.templates[name]
	if !ok {
		return "", template.ErrTemplateNotFound
	}

	e.logger.Debug().Msgf("render template: %s", name)
	tmpl, err := htmltemplate.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var tmplBuf bytes.Buffer
	if err := tmpl.Execute(&tmplBuf, args); err != nil {
		return "", err
	}

	return tmplBuf.String(), nil
}

func (e *engine) RenderText(name string, args any) (string, error) {
	e.logger.Debug().Msgf("search template: %s", name)
	tmplPath, ok := e.templates[name]
	if !ok {
		return "", template.ErrTemplateNotFound
	}

	e.logger.Debug().Msgf("render template: %s", name)
	tmpl, err := texttemplate.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var tmplBuf bytes.Buffer
	if err := tmpl.Execute(&tmplBuf, args); err != nil {
		return "", err
	}

	return tmplBuf.String(), nil
}
