package factory

import (
	"errors"
	"github.com/mandarine-io/backend/third_party/oauth"
	"github.com/mandarine-io/backend/third_party/oauth/google"
	"github.com/mandarine-io/backend/third_party/oauth/mailru"
	"github.com/mandarine-io/backend/third_party/oauth/yandex"
)

var (
	ErrUnsupportedOAuthProvider = errors.New("unsupported OAuth provider")
)

func NewProviderByKey(key string, clientID string, clientSecret string, opts ...oauth.Option) (oauth.Provider, error) {
	switch key {
	case google.ProviderKey:
		return google.NewProvider(clientID, clientSecret, opts...), nil
	case yandex.ProviderKey:
		return yandex.NewProvider(clientID, clientSecret, opts...), nil
	case mailru.ProviderKey:
		return mailru.NewProvider(clientID, clientSecret, opts...), nil
	default:
		return nil, ErrUnsupportedOAuthProvider
	}
}
