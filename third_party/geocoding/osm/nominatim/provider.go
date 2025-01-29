package nominatim

import (
	"github.com/mandarine-io/backend/third_party/geocoding"
)

const (
	ProviderKey              = "osm_nominatim"
	DefaultGeocodeURL        = "https://nominatim.openstreetmap.org/search"
	DefaultReverseGeocodeURL = "https://nominatim.openstreetmap.org/reverse"
)

func NewProvider(opts ...geocoding.Option) geocoding.Provider {
	return geocoding.NewBaseProvider(
		&endpointBuilder{DefaultGeocodeURL, DefaultReverseGeocodeURL},
		func() geocoding.ResponseParser { return &responseParser{} },
		opts...,
	)
}
