package factory

import (
	"errors"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/mandarine-io/backend/third_party/geocoding/chained"
	graphhopper "github.com/mandarine-io/backend/third_party/geocoding/graphhoper"
	"github.com/mandarine-io/backend/third_party/geocoding/here"
	"github.com/mandarine-io/backend/third_party/geocoding/locationiq"
	"github.com/mandarine-io/backend/third_party/geocoding/osm/nominatim"
	"github.com/mandarine-io/backend/third_party/geocoding/yandex"
)

var (
	ErrUnsupportedGeocodingProvider = errors.New("unsupported geocoding provider")
)

func NewProviderByKey(key string, apiKey string, opts ...geocoding.Option) (geocoding.Provider, error) {
	switch key {
	case here.ProviderKey:
		return here.NewProvider(apiKey, opts...), nil
	case graphhopper.ProviderKey:
		return graphhopper.NewProvider(apiKey, opts...), nil
	case nominatim.ProviderKey:
		return nominatim.NewProvider(opts...), nil
	case locationiq.ProviderKey:
		return locationiq.NewProvider(apiKey, opts...), nil
	case yandex.ProviderKey:
		return yandex.NewProvider(apiKey, opts...), nil
	default:
		return nil, ErrUnsupportedGeocodingProvider
	}
}

func NewProviderChained(providers ...geocoding.Provider) geocoding.Provider {
	return chained.NewProvider(providers...)
}
