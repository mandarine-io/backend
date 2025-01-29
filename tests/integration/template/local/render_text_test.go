package local

import (
	"github.com/mandarine-io/backend/internal/infrastructure/template"
	"github.com/mandarine-io/backend/internal/infrastructure/template/local"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"os"
)

var (
	renderTextEngine      template.Engine
	renderTextTemplateDir string
)

type RenderTextSuite struct {
	suite.Suite
}

func (s *RenderTextSuite) BeforeEach(t provider.T) {
	t.Title("Render - before each")
	t.Feature("Local template")

	var err error
	renderTextTemplateDir, err = os.MkdirTemp("", "local_template_engine_RenderText_*")
	t.Require().NoError(err)

	file, err := os.Create(renderTextTemplateDir + "/template.gotmpl")
	t.Require().NoError(err)

	defer func(file *os.File) {
		err := file.Close()
		t.Require().NoError(err)
	}(file)

	_, err = file.WriteString("Timestamp: {{.Timestamp}}")
	t.Require().NoError(err)

	renderTextEngine, err = local.NewEngine(renderTextTemplateDir)
	t.Require().NoError(err)
}

func (s *RenderTextSuite) AfterEach(t provider.T) {
	t.Title("Render - after each")
	t.Feature("Local template")

	err := os.RemoveAll(renderTextTemplateDir)
	t.Require().NoError(err)
}

func (s *RenderTextSuite) Test_RenderText_Success(t provider.T) {
	t.Title("Render - success")
	t.Severity(allure.NORMAL)
	t.Feature("Local template")
	t.Tags("Positive")
	t.Parallel()

	res, err := renderTextEngine.RenderText("template", map[string]any{"Timestamp": "2022-01-01 00:00:00"})
	t.Require().NoError(err)
	t.Require().Equal("Timestamp: 2022-01-01 00:00:00", res)
}

func (s *RenderTextSuite) Test_RenderText_NotFound(t provider.T) {
	t.Title("Render - not found")
	t.Severity(allure.CRITICAL)
	t.Feature("Local template")
	t.Tags("Negative")
	t.Parallel()

	_, err := renderTextEngine.RenderText("template_non_exist", map[string]any{"Timestamp": "2022-01-01 00:00:00"})
	t.Require().Error(err)
	t.Require().ErrorIs(err, template.ErrTemplateNotFound)
}
