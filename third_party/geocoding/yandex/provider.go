package yandex

import (
	"github.com/mandarine-io/backend/third_party/geocoding"
)

const (
	ProviderKey              = "yandex"
	DefaultGeocodeURL        = "https://geocode-maps.yandex.ru/1.x"
	DefaultReverseGeocodeURL = "https://geocode-maps.yandex.ru/1.x"
)

func NewProvider(apiKey string, opts ...geocoding.Option) geocoding.Provider {
	return geocoding.NewBaseProvider(
		&endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		func() geocoding.ResponseParser { return &responseParser{} },
		opts...,
	)
}
