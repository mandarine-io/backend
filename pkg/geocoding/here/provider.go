package here

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/http"
)

const (
	ProviderKey              = "here"
	DefaultGeocodeURL        = "https://geocode.search.hereapi.com/v1/geocode"
	DefaultReverseGeocodeURL = "https://revgeocode.search.hereapi.com/v1/revgeocode"
)

func NewProvider(apiKey string) geocoding.Provider {
	return &http.Provider{
		EndpointBuilder:       &endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		ResponseParserFactory: func() http.ResponseParser { return &responseParser{} },
	}
}
