package locationiq

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/http"
)

const (
	ProviderKey              = "locationiq"
	DefaultGeocodeURL        = "https://locationiq.org/v1/search.php"
	DefaultReverseGeocodeURL = "https://locationiq.org/v1/reverse.php"
)

func NewProvider(apiKey string) geocoding.Provider {
	return &http.Provider{
		EndpointBuilder:       &endpointBuilder{DefaultGeocodeURL, DefaultReverseGeocodeURL, apiKey},
		ResponseParserFactory: func() http.ResponseParser { return &responseParser{} },
	}
}
