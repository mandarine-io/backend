package local

import (
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/internal/infrastructure/locale/local"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"os"
)

var (
	bundle    locale.Bundle
	localeDir string
)

type LocalizeSuite struct {
	suite.Suite
}

func (s *LocalizeSuite) BeforeEach(t provider.T) {
	t.Title("Render - before each")
	t.Feature("Local localizer")

	var err error
	localeDir, err = os.MkdirTemp("", "local_locale_*")
	t.Require().NoError(err)

	file, err := os.Create(localeDir + "/en.json")
	t.Require().NoError(err)

	defer func(file *os.File) {
		err := file.Close()
		t.Require().NoError(err)
	}(file)

	_, err = file.WriteString("{\"test.message\":\"test\"}")
	t.Require().NoError(err)

	bundle, err = local.NewBundle(localeDir)
	t.Require().NoError(err)
}

func (s *LocalizeSuite) AfterEach(t provider.T) {
	t.Title("Render - after each")
	t.Feature("Local localizer")

	err := os.RemoveAll(localeDir)
	t.Require().NoError(err)
}

func (s *LocalizeSuite) Test_Localize_Success(t provider.T) {
	t.Title("Localize - success")
	t.Severity(allure.NORMAL)
	t.Feature("Local localizer")
	t.Tags("Positive")
	t.Parallel()

	localizer := bundle.NewLocalizer("en")

	res := localizer.Localize("test.message", nil, 0)
	t.Require().Equal("test", res)
}

func (s *LocalizeSuite) Test_RenderHTML_NotFound(t provider.T) {
	t.Title("Localize - not found")
	t.Severity(allure.CRITICAL)
	t.Feature("Local localizer")
	t.Tags("Negative")
	t.Parallel()

	localizer := bundle.NewLocalizer("en")

	res := localizer.Localize("test.not_found", nil, 0)
	t.Require().Equal("test.not_found", res)
}
