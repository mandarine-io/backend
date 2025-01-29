package graphhopper

import (
	"github.com/mandarine-io/backend/third_party/geocoding"
)

const (
	ProviderKey              = "graphhopper"
	DefaultGeocodeURL        = "https://graphhopper.com/api/1/geocode"
	DefaultReverseGeocodeURL = "https://graphhopper.com/api/1/geocode"
)

func NewProvider(apiKey string, opts ...geocoding.Option) geocoding.Provider {
	return geocoding.NewBaseProvider(
		&endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		func() geocoding.ResponseParser { return &responseParser{} },
		opts...,
	)
}
