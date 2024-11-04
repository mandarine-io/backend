package factory

import (
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/mandarine-io/Backend/pkg/oauth/google"
	"github.com/mandarine-io/Backend/pkg/oauth/mailru"
	"github.com/mandarine-io/Backend/pkg/oauth/yandex"
)

func NewProviderByKey(key string, clientID string, clientSecret string) oauth.Provider {
	switch key {
	case google.ProviderKey:
		return google.NewProvider(clientID, clientSecret)
	case yandex.ProviderKey:
		return yandex.NewProvider(clientID, clientSecret)
	case mailru.ProviderKey:
		return mailru.NewProvider(clientID, clientSecret)
	default:
		return nil
	}
}
