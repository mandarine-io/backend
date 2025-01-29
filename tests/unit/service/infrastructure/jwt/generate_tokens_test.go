package jwt

import (
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"strings"
)

type GenerateTokensSuite struct {
	suite.Suite
}

func (s *GenerateTokensSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("JWT service")
	t.Feature("GenerateTokens")
	t.Tags("Positive")

	accessToken, refreshToken, err := svc.GenerateTokens(ctx, &entity.User{})

	t.Require().NoError(err)
	t.Require().NotEmpty(accessToken)
	t.Require().Len(strings.Split(accessToken, "."), 3)
	t.Require().NotEmpty(refreshToken)
	t.Require().Len(strings.Split(refreshToken, "."), 3)
}
