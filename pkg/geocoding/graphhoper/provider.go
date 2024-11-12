package graphhopper

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/http"
)

const (
	ProviderKey              = "graphhopper"
	DefaultGeocodeURL        = "https://graphhopper.com/api/1/geocode"
	DefaultReverseGeocodeURL = "https://graphhopper.com/api/1/geocode"
)

func NewProvider(apiKey string) geocoding.Provider {
	return &http.Provider{
		EndpointBuilder:       &endpointBuilder{apiKey, DefaultGeocodeURL, DefaultReverseGeocodeURL},
		ResponseParserFactory: func() http.ResponseParser { return &responseParser{} },
	}
}
