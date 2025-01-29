package here

import (
	"github.com/mandarine-io/backend/third_party/geocoding"
)

const (
	ProviderKey              = "here"
	DefaultGeocodeURL        = "https://geocode.search.hereapi.com/v1/geocode"
	DefaultReverseGeocodeURL = "https://revgeocode.search.hereapi.com/v1/revgeocode"
)

func NewProvider(apiKey string, opts ...geocoding.Option) geocoding.Provider {
	return geocoding.NewBaseProvider(
		&endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		func() geocoding.ResponseParser { return &responseParser{} },
		opts...,
	)
}
