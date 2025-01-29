package chained

import (
	"atomicgo.dev/robin"
	"context"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/rs/zerolog/log"
)

type provider struct {
	providers *robin.Loadbalancer[geocoding.Provider]
	len       int
}

func NewProvider(providers ...geocoding.Provider) geocoding.Provider {
	return &provider{
		providers: robin.NewLoadbalancer(providers),
		len:       len(providers),
	}
}

func (p *provider) Geocode(ctx context.Context, address string, config geocoding.GeocodeConfig) (
	[]*geocoding.Location,
	error,
) {
	for i := 0; i < p.len; i++ {
		nextProvider := p.providers.Next()
		loc, err := nextProvider.Geocode(ctx, address, config)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to geocode address")
			continue
		}
		return loc, nil
	}
	return nil, geocoding.ErrGeocodeProvidersUnavailable
}

func (p *provider) ReverseGeocode(
	ctx context.Context,
	loc geocoding.Location,
	config geocoding.ReverseGeocodeConfig,
) ([]*geocoding.Address, error) {
	for i := 0; i < p.len; i++ {
		nextProvider := p.providers.Next()
		addr, err := nextProvider.ReverseGeocode(ctx, loc, config)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to reverse geocode address")
			continue
		}
		return addr, nil
	}
	return nil, geocoding.ErrGeocodeProvidersUnavailable
}
