package google

import (
	"encoding/json"
	"github.com/mandarine-io/backend/third_party/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ProviderKey = "google"
)

func NewProvider(clientID string, clientSecret string, opts ...oauth.Option) oauth.Provider {
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return oauth.NewBaseProvider(oauthConfig, "https://www.googleapis.com/oauth2/v2/userinfo", UnmarshalJSON, opts...)
}

//////////////////// Marshall User Info ////////////////////

type UserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

func UnmarshalJSON(data []byte) (oauth.UserInfo, error) {
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return oauth.UserInfo{}, err
	}

	return oauth.UserInfo{
		Username:        userInfo.Name,
		Email:           userInfo.Email,
		IsEmailVerified: userInfo.VerifiedEmail,
	}, nil
}
