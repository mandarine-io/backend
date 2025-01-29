package template

import (
	"fmt"
)

var (
	ErrTemplateNotFound = fmt.Errorf("template not found")
)

type Engine interface {
	RenderText(name string, args any) (string, error)
	RenderHTML(name string, args any) (string, error)
}
