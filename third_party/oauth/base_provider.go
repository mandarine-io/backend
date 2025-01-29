package oauth

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

type Provider interface {
	GetConsentPageURL(redirectURL string) (string, string)
	ExchangeCodeToToken(ctx context.Context, code string, redirectURL string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}

type UnmarshallUserInfo = func(bytes []byte) (UserInfo, error)

type baseProvider struct {
	oauthConfig        *oauth2.Config
	userInfoURL        string
	unmarshallUserInfo UnmarshallUserInfo
	logger             zerolog.Logger
}

type Option func(*baseProvider)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *baseProvider) {
		p.logger = logger
	}
}

func NewBaseProvider(
	oauthConfig *oauth2.Config,
	userInfoURL string,
	unmarshallUserInfo UnmarshallUserInfo,
	opts ...Option,
) Provider {
	p := &baseProvider{
		oauthConfig:        oauthConfig,
		userInfoURL:        userInfoURL,
		unmarshallUserInfo: unmarshallUserInfo,
		logger:             zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (c *baseProvider) GetConsentPageURL(redirectURL string) (string, string) {
	c.logger.Debug().Msg("get consent page url")

	oauthState := uuid.New().String()
	redirectURISetter := oauth2.SetAuthURLParam("redirect_uri", redirectURL)

	return c.oauthConfig.AuthCodeURL(oauthState, redirectURISetter), oauthState
}

func (c *baseProvider) ExchangeCodeToToken(ctx context.Context, code string, redirectURL string) (
	*oauth2.Token,
	error,
) {
	c.logger.Debug().Msg("exchange code to token")

	redirectURISetter := oauth2.SetAuthURLParam("redirect_uri", redirectURL)

	return c.oauthConfig.Exchange(ctx, code, redirectURISetter)
}

func (c *baseProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error) {
	// Create request
	c.logger.Debug().Msg("send request to get user info")
	req, _ := http.NewRequestWithContext(ctx, "GET", c.userInfoURL, nil)

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
		err = res.Body.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("failed to close response body")
		}
	}()

	// Check status
	if res.StatusCode >= 400 {
		return UserInfo{}, ErrUserInfoNotReceived
	}

	// Convert response body
	c.logger.Debug().Msg("read response body")
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return UserInfo{}, err
	}

	c.logger.Debug().Msg("unmarshal user info to unnecessary struct")
	return c.unmarshallUserInfo(bytes)
}
