package mock

import (
	"context"
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type OAuthProviderMock struct {
	mock.Mock
}

func (m *OAuthProviderMock) GetConsentPageUrl(redirectUrl string) (string, string) {
	args := m.Called(redirectUrl)
	return args.String(0), args.String(1)
}

func (m *OAuthProviderMock) ExchangeCodeToToken(ctx context.Context, code string, redirectUrl string) (*oauth2.Token, error) {
	args := m.Called(ctx, code, redirectUrl)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *OAuthProviderMock) GetUserInfo(ctx context.Context, token *oauth2.Token) (oauth.UserInfo, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(oauth.UserInfo), args.Error(1)
}
