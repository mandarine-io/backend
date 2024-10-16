package template

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"mandarine/internal/api/helper/file"
	"mandarine/pkg/logging"
	"os"
	"path"
	"path/filepath"
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
		slog.Error("Templates set up error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	tmplEngine := &engine{
		templates: make(map[string]string),
	}
	for _, f := range files {
		filePath := path.Join(cfg.Path, f)
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			slog.Error("Templates set up error", logging.ErrorAttr(err))
			os.Exit(1)
		}

		tmplName := file.GetFileNameWithoutExt(f)
		slog.Info("Read template file: " + absFilePath)
		tmplEngine.templates[tmplName] = absFilePath
	}

	return tmplEngine
}

func (t *engine) Render(tmplName string, args any) (string, error) {
	slog.Debug("Search template: " + tmplName)
	tmplPath, ok := t.templates[tmplName]
	if !ok {
		return "", fmt.Errorf("template \"%s\"not found", tmplName)
	}

	slog.Debug("Execute template: " + tmplName)
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
