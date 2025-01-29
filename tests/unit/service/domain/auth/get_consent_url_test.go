package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

const (
	redirectURL = "https://example.com/callback"
)

type GetConsentPageURLSuite struct {
	suite.Suite
}

func (s *GetConsentPageURLSuite) Test_NotSupportedProvider(t provider.T) {
	t.Title("GetConsentPageURL not support provider")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("GetConsentPageURL")
	t.Tags("Negative")

	_, err := svc.GetConsentPageURL(context.Background(), "unsupported", redirectURL)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrInvalidProvider, err)
}

func (s *GetConsentPageURLSuite) Test_Success(t provider.T) {
	t.Title("GetConsentPageURL return success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("GetConsentPageURL")
	t.Tags("Positive")

	oauthProviderMock.On("GetConsentPageURL", redirectURL).
		Return("consentURL", "oauthState").Once()

	result, err := svc.GetConsentPageURL(context.Background(), "mock", redirectURL)

	t.Require().NoError(err)
	t.Require().Equal("consentURL", result.ConsentPageURL)
	t.Require().Equal("oauthState", result.OauthState)
}
