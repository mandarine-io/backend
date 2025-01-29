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
	renderHTMLEngine      template.Engine
	renderHTMLTemplateDir string
)

type RenderHTMLSuite struct {
	suite.Suite
}

func (s *RenderHTMLSuite) BeforeAll(t provider.T) {
	t.Title("Render - before each")
	t.Feature("Local template")

	var err error
	renderHTMLTemplateDir, err = os.MkdirTemp("", "local_template_engine_RenderHTML_*")
	t.Require().NoError(err)

	file, err := os.Create(renderHTMLTemplateDir + "/template.html")
	t.Require().NoError(err)

	defer func(file *os.File) {
		err := file.Close()
		t.Require().NoError(err)
	}(file)

	_, err = file.WriteString("<h1>Timestamp: {{.Timestamp}}</h1>")
	t.Require().NoError(err)

	renderHTMLEngine, err = local.NewEngine(renderHTMLTemplateDir)
	t.Require().NoError(err)
}

func (s *RenderHTMLSuite) AfterAll(t provider.T) {
	t.Title("Render - after each")
	t.Feature("Local template")

	err := os.RemoveAll(renderHTMLTemplateDir)
	t.Require().NoError(err)
}

func (s *RenderHTMLSuite) Test_RenderHTML_Success(t provider.T) {
	t.Title("Render - success")
	t.Severity(allure.NORMAL)
	t.Feature("Local template")
	t.Tags("Positive")
	t.Parallel()

	res, err := renderHTMLEngine.RenderHTML("template", map[string]any{"Timestamp": "2022-01-01 00:00:00"})
	t.Require().NoError(err)
	t.Require().Equal("<h1>Timestamp: 2022-01-01 00:00:00</h1>", res)
}

func (s *RenderHTMLSuite) Test_RenderHTML_NotFound(t provider.T) {
	t.Title("Render - not found")
	t.Severity(allure.CRITICAL)
	t.Feature("Local template")
	t.Tags("Negative")
	t.Parallel()

	_, err := renderHTMLEngine.RenderHTML("template_non_exist", map[string]any{"Timestamp": "2022-01-01 00:00:00"})
	t.Require().Error(err)
	t.Require().ErrorIs(err, template.ErrTemplateNotFound)
}
