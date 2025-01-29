package mailru

import (
	"encoding/json"
	"github.com/mandarine-io/backend/third_party/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/mailru"
)

const (
	ProviderKey = "mailru"
)

func NewProvider(clientID string, clientSecret string, opts ...oauth.Option) oauth.Provider {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"userinfo"},
		Endpoint:     mailru.Endpoint,
	}
	return oauth.NewBaseProvider(oauthConfig, "https://oauth.mail.ru/userinfo", UnmarshalJSON, opts...)
}

//////////////////// Marshall User Info ////////////////////

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func UnmarshalJSON(data []byte) (oauth.UserInfo, error) {
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return oauth.UserInfo{}, err
	}

	return oauth.UserInfo{
		Username:        userInfo.Name,
		Email:           userInfo.Email,
		IsEmailVerified: true,
	}, nil
}
