package factory

import (
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/chained"
	graphhopper "github.com/mandarine-io/Backend/pkg/geocoding/graphhoper"
	"github.com/mandarine-io/Backend/pkg/geocoding/here"
	"github.com/mandarine-io/Backend/pkg/geocoding/locationiq"
	"github.com/mandarine-io/Backend/pkg/geocoding/osm/nominatim"
	"github.com/mandarine-io/Backend/pkg/geocoding/yandex"
)

func NewProviderByKey(key string, apiKey string) geocoding.Provider {
	switch key {
	case here.ProviderKey:
		return here.NewProvider(apiKey)
	case graphhopper.ProviderKey:
		return graphhopper.NewProvider(apiKey)
	case nominatim.ProviderKey:
		return nominatim.NewProvider()
	case locationiq.ProviderKey:
		return locationiq.NewProvider(apiKey)
	case yandex.ProviderKey:
		return yandex.NewProvider(apiKey)
	default:
		return nil
	}
}

func NewProviderChained(providers ...geocoding.Provider) geocoding.Provider {
	return chained.NewProvider(providers...)
}
