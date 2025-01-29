package local

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

type LocalLocalizerSuite struct {
	suite.Suite
}

func TestLocalLocalizerSuite(t *testing.T) {
	suite.RunSuite(t, new(LocalLocalizerSuite))
}

func (s *LocalLocalizerSuite) Test(t provider.T) {
	s.RunSuite(t, new(LocalizeSuite))
}
