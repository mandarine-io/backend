package yandex

import (
	"encoding/json"
	"github.com/mandarine-io/Backend/pkg/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

const (
	ProviderKey = "yandex"
)

func NewProvider(clientID string, clientSecret string) oauth.Provider {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{},
		Endpoint:     yandex.Endpoint,
	}
	return oauth.NewProvider(oauthConfig, "https://login.yandex.ru/info", UnmarshalJSON)
}

//////////////////// Marshall User Info ////////////////////

type UserInfo struct {
	DefaultEmail string `json:"default_email"`
	DisplayName  string `json:"display_name"`
}

func UnmarshalJSON(data []byte) (oauth.UserInfo, error) {
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return oauth.UserInfo{}, err
	}

	return oauth.UserInfo{
		Username:        userInfo.DisplayName,
		Email:           userInfo.DefaultEmail,
		IsEmailVerified: true,
	}, nil
}
