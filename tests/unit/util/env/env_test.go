package env

import (
	"github.com/mandarine-io/backend/internal/util/env"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"os"
	"testing"
)

type EnvUtilSuite struct {
	suite.Suite
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(EnvUtilSuite))
}

func (s *EnvUtilSuite) Test_GetEnvWithDefault_EnvExists(t provider.T) {
	t.Title("GetEnvWithDefault")
	t.Severity(allure.NORMAL)
	t.Tag("positive")
	t.Parallel()

	expected := "value"
	err := os.Setenv("TEST_ENV_1", expected)
	t.Require().NoError(err)

	result := env.GetEnvWithDefault("TEST_ENV_1", "default")
	t.Require().Equal(expected, result)
}

func (s *EnvUtilSuite) Test_GetEnvWithDefault_EnvNotExists(t provider.T) {
	t.Title("GetEnvWithDefault - env not exists")
	t.Severity(allure.NORMAL)
	t.Tag("positive")
	t.Parallel()

	expected := "default"
	result := env.GetEnvWithDefault("TEST_ENV_2", expected)
	t.Require().Equal(expected, result)
}
