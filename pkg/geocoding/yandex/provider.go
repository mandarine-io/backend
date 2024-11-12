package yandex

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/http"
)

const (
	ProviderKey              = "yandex"
	DefaultGeocodeURL        = "https://geocode-maps.yandex.ru/1.x"
	DefaultReverseGeocodeURL = "https://geocode-maps.yandex.ru/1.x"
)

func NewProvider(apiKey string) geocoding.Provider {
	return &http.Provider{
		EndpointBuilder:       &endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		ResponseParserFactory: func() http.ResponseParser { return &responseParser{} },
	}
}
