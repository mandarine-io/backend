package nominatim

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/http"
)

const (
	ProviderKey              = "osm_nominatim"
	DefaultGeocodeURL        = "https://nominatim.openstreetmap.org/search"
	DefaultReverseGeocodeURL = "https://nominatim.openstreetmap.org/reverse"
)

func NewProvider() geocoding.Provider {
	return &http.Provider{
		EndpointBuilder:       &endpointBuilder{DefaultGeocodeURL, DefaultReverseGeocodeURL},
		ResponseParserFactory: func() http.ResponseParser { return &responseParser{} },
	}
}
