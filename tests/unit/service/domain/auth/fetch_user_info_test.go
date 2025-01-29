package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/third_party/oauth"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type FetchUserInfoSuite struct {
	suite.Suite
}

func (s *FetchUserInfoSuite) Test_NotSupportedProvider(t provider.T) {
	t.Title("FetchUserInfo not support provider")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("FetchUserInfo")
	t.Tags("Negative")

	input := v0.FetchUserInfoInput{Code: "someCode"}
	_, err := svc.FetchUserInfo(context.Background(), "unsupported", input)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrInvalidProvider, err)
}

func (s *FetchUserInfoSuite) Test_Success(t provider.T) {
	t.Title("FetchUserInfo return success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("FetchUserInfo")
	t.Tags("Positive")

	input := v0.FetchUserInfoInput{Code: "someCode"}
	expectedUserInfo := oauth.UserInfo{Email: "test@example.com"}

	oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(
		&oauth2.Token{},
		nil,
	).Once()
	oauthProviderMock.On("GetUserInfo", mock.Anything, mock.Anything).Return(expectedUserInfo, nil).Once()

	userInfo, err := svc.FetchUserInfo(context.Background(), "mock", input)

	t.Require().NoError(err)
	t.Require().Equal(expectedUserInfo, userInfo)
}

func (s *FetchUserInfoSuite) Test_ErrorExchangingCodeToToken(t provider.T) {
	t.Title("FetchUserInfo return ExchangingCodeToToken error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("FetchUserInfo")
	t.Tags("Negative")

	input := v0.FetchUserInfoInput{Code: "someCode"}
	expectedError := errors.New("exchange error")

	oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(
		nil,
		expectedError,
	).Once()

	_, err := svc.FetchUserInfo(context.Background(), "mock", input)

	t.Require().Error(err)
	t.Require().Equal(expectedError, err)
}

func (s *FetchUserInfoSuite) Test_ErrorGettingUserInfo(t provider.T) {
	t.Title("FetchUserInfo return GettingUserInfo error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("FetchUserInfo")
	t.Tags("Negative")

	input := v0.FetchUserInfoInput{Code: "someCode"}
	token := &oauth2.Token{}
	expectedError := errors.New("user info error")

	oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(token, nil).Once()
	oauthProviderMock.On("GetUserInfo", mock.Anything, token).Return(oauth.UserInfo{}, expectedError).Once()

	_, err := svc.FetchUserInfo(context.Background(), "mock", input)

	t.Require().Error(err)
	t.Require().Equal(expectedError, err)
}
