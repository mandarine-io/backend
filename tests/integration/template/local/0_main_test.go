package local

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

type LocalTemplateEngineSuite struct {
	suite.Suite
}

func TestLocalTemplateEngineSuite(t *testing.T) {
	suite.RunSuite(t, new(LocalTemplateEngineSuite))
}

func (s *LocalTemplateEngineSuite) Test(t provider.T) {
	s.RunSuite(t, new(RenderTextSuite))
	s.RunSuite(t, new(RenderHTMLSuite))
}
