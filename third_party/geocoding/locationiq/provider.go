package locationiq

import (
	"github.com/mandarine-io/backend/third_party/geocoding"
)

const (
	ProviderKey              = "locationiq"
	DefaultGeocodeURL        = "https://locationiq.org/v1/search.php"
	DefaultReverseGeocodeURL = "https://locationiq.org/v1/reverse.php"
)

func NewProvider(apiKey string, opts ...geocoding.Option) geocoding.Provider {
	return geocoding.NewBaseProvider(
		&endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		func() geocoding.ResponseParser { return &responseParser{} },
		opts...,
	)
}
