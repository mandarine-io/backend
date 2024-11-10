package template

import (
	"bytes"
	"fmt"
	"github.com/mandarine-io/Backend/pkg/helper/file"
	"github.com/rs/zerolog/log"
	"html/template"
	"path"
	"path/filepath"
)

var (
	ErrTemplateNotFound = fmt.Errorf("template not found")
)

type Engine interface {
	Render(tmplName string, args any) (string, error)
}

type engine struct {
	templates map[string]string
}

type Config struct {
	Path string
}

func MustLoadTemplates(cfg *Config) Engine {
	files, err := file.GetFilesFromDir(cfg.Path)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get templates from directory")
	}

	tmplEngine := &engine{
		templates: make(map[string]string),
	}
	for _, f := range files {
		filePath := path.Join(cfg.Path, f)
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			log.Fatal().Stack().Err(err).Msgf("failed to get absolute path for file %s", filePath)
		}

		tmplName := file.GetFileNameWithoutExt(f)
		log.Info().Msgf("read template file: %s", absFilePath)
		tmplEngine.templates[tmplName] = absFilePath
	}

	return tmplEngine
}

func (t *engine) Render(tmplName string, args any) (string, error) {
	log.Debug().Msgf("search template: %s", tmplName)
	tmplPath, ok := t.templates[tmplName]
	if !ok {
		return "", ErrTemplateNotFound
	}

	log.Debug().Msgf("render template: %s", tmplName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	var tmplBuf bytes.Buffer
	if err := tmpl.Execute(&tmplBuf, args); err != nil {
		return "", err
	}

	return tmplBuf.String(), nil
}
