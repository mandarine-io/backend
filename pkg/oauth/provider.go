package oauth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io"
	"log/slog"
	"mandarine/pkg/logging"
	"net/http"
)

var (
	ErrUserInfoNotReceived = fmt.Errorf("user info not received")
)

type UnmarshallUserInfo = func(bytes []byte) (UserInfo, error)

type Provider interface {
	GetConsentPageUrl(redirectUrl string) (string, string)
	ExchangeCodeToToken(ctx context.Context, code string, redirectUrl string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}

type provider struct {
	oauthConfig        *oauth2.Config
	userInfoUrl        string
	unmarshallUserInfo UnmarshallUserInfo
}

func NewProvider(oauthConfig *oauth2.Config, userInfoUrl string, unmarshallUserInfo UnmarshallUserInfo) Provider {
	return &provider{
		oauthConfig:        oauthConfig,
		userInfoUrl:        userInfoUrl,
		unmarshallUserInfo: unmarshallUserInfo,
	}
}

func (c *provider) GetConsentPageUrl(redirectUrl string) (string, string) {
	oauthState := uuid.New().String()
	redirectUriSetter := oauth2.SetAuthURLParam("redirect_uri", redirectUrl)
	return c.oauthConfig.AuthCodeURL(oauthState, redirectUriSetter), oauthState
}

func (c *provider) ExchangeCodeToToken(ctx context.Context, code string, redirectUrl string) (*oauth2.Token, error) {
	redirectUriSetter := oauth2.SetAuthURLParam("redirect_uri", redirectUrl)
	return c.oauthConfig.Exchange(ctx, code, redirectUriSetter)
}

type UserInfo struct {
	Username        string
	Email           string
	IsEmailVerified bool
}

func (c *provider) GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error) {
	// Create request
	req, _ := http.NewRequest("GET", c.userInfoUrl, nil)

	query := req.URL.Query()
	query.Add("access_token", token.AccessToken)
	req.URL.RawQuery = query.Encode()

	// Send request
	client2 := c.oauthConfig.Client(ctx, token)
	res, err := client2.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			slog.Warn("Get user info error", logging.ErrorAttr(err))
		}
	}()

	// Check status
	if res.StatusCode >= 400 {
		return UserInfo{}, ErrUserInfoNotReceived
	}

	// Convert response body
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return UserInfo{}, err
	}

	return c.unmarshallUserInfo(bytes)
}
