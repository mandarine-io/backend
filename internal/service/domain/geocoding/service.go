package geocoding

import (
	"context"
	"errors"
	"github.com/mandarine-io/backend/internal/converter"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/service/domain"
	cachehelper "github.com/mandarine-io/backend/internal/util/cache"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
)

const (
	geocodeCachePrefix  = "geocode"
	reverseGeocodeCache = "reverse_geocode"
)

type svc struct {
	provider geocoding.Provider
	manager  cache.Manager
	logger   zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(manager cache.Manager, provider geocoding.Provider, opts ...Option) domain.GeocodingService {
	s := &svc{
		provider: provider,
		manager:  manager,
		logger:   zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *svc) Geocode(ctx context.Context, input v0.GeocodingInput, lang language.Tag) (
	v0.GeocodingOutput,
	error,
) {
	s.logger.Info().Msgf("geocode address: %s", input.Address)

	var (
		locs []*geocoding.Location
		err  error
	)

	// Get from cache
	key := cachehelper.CreateCacheKey(geocodeCachePrefix, input.Address)
	err = s.manager.Get(ctx, key, &locs)
	if err == nil {
		return converter.MapLocationsToGeocodingOutput(locs), nil
	}
	if !errors.Is(err, cache.ErrCacheEntryNotFound) {
		s.logger.Warn().Stack().Err(err).Msg("failed to get geocode from cache")
	}

	// Get from provider
	locs, err = s.provider.Geocode(ctx, input.Address, geocoding.GeocodeConfig{Lang: lang, Limit: input.Limit})
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to geocode address")
		if errors.Is(err, geocoding.ErrGeocodeProvidersUnavailable) {
			return v0.GeocodingOutput{}, domain.ErrGeocodeProvidersUnavailable
		}
		return v0.GeocodingOutput{}, err
	}

	// Save	to cache
	err = s.manager.Set(ctx, key, locs)
	if err != nil {
		s.logger.Warn().Stack().Err(err).Msg("failed to set geocode to cache")
	}

	return converter.MapLocationsToGeocodingOutput(locs), nil
}

func (s *svc) ReverseGeocode(
	ctx context.Context, input v0.ReverseGeocodingInput, lang language.Tag,
) (v0.ReverseGeocodingOutput, error) {
	s.logger.Info().Msgf("reverse geocode point: [%s, %s]", input.Lng, input.Lat)

	var (
		addrs []*geocoding.Address
		err   error
	)

	// Get from cache
	key := cachehelper.CreateCacheKey(
		reverseGeocodeCache,
		input.Lng.String(),
		input.Lat.String(),
	)
	err = s.manager.Get(ctx, key, &addrs)
	if err == nil {
		return converter.MapAddressesToGeocodingOutput(addrs), nil
	}
	if !errors.Is(err, cache.ErrCacheEntryNotFound) {
		s.logger.Warn().Stack().Err(err).Msg("failed to get reverse geocode from cache")
	}

	// Get from provider
	addrs, err = s.provider.ReverseGeocode(
		ctx,
		converter.MapLngLatToLocation(input.Lng, input.Lat),
		geocoding.ReverseGeocodeConfig{Lang: lang, Limit: input.Limit, Zoom: 18},
	)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to reverse geocode address")
		if errors.Is(err, geocoding.ErrGeocodeProvidersUnavailable) {
			return v0.ReverseGeocodingOutput{}, domain.ErrGeocodeProvidersUnavailable
		}
		return v0.ReverseGeocodingOutput{}, err
	}

	// Save	to cache
	err = s.manager.Set(ctx, key, addrs)
	if err != nil {
		s.logger.Warn().Stack().Err(err).Msg("failed to set reverse geocode to cache")
	}

	return converter.MapAddressesToGeocodingOutput(addrs), nil
}
