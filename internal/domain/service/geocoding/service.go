package geocoding

import (
	"context"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/geocoding/mapper"
	cache2 "github.com/mandarine-io/Backend/internal/helper/cache"
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/mandarine-io/Backend/pkg/geocoding/chained"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"strconv"
)

const (
	geocodeCachePrefix  = "geocode"
	reverseGeocodeCache = "reverse_geocode"
)

type svc struct {
	provider geocoding.Provider
	manager  cache.Manager
}

func NewService(provider geocoding.Provider, manager cache.Manager) service.GeocodingService {
	return &svc{provider: provider, manager: manager}
}

func (s *svc) Geocode(ctx context.Context, input dto.GeocodingInput, lang language.Tag) (dto.GeocodingOutput, error) {
	log.Info().Msgf("geocode address: %s", input.Address)

	var (
		locs []*geocoding.Location
		err  error
	)

	// Get from cache
	key := cache2.CreateCacheKey(geocodeCachePrefix, input.Address)
	err = s.manager.Get(ctx, key, &locs)
	if err == nil {
		return mapper.MapLocationsToGeocodingOutput(locs), nil
	}
	if !errors.Is(err, cache.ErrCacheEntryNotFound) {
		log.Warn().Stack().Err(err).Msg("failed to get geocode from cache")
	}

	// Get from provider
	locs, err = s.provider.GeocodeWithContext(ctx, input.Address, geocoding.GeocodeConfig{Lang: lang, Limit: input.Limit})
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to geocode address")
		if errors.Is(err, chained.ErrGeocodeProvidersUnavailable) {
			return dto.GeocodingOutput{}, service.ErrGeocodeProvidersUnavailable
		}
		return dto.GeocodingOutput{}, err
	}

	// Save	to cache
	err = s.manager.Set(ctx, key, locs)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to set geocode to cache")
	}

	return mapper.MapLocationsToGeocodingOutput(locs), nil
}

func (s *svc) ReverseGeocode(
	ctx context.Context, input dto.ReverseGeocodingInput, lang language.Tag,
) (dto.ReverseGeocodingOutput, error) {
	log.Info().Msgf("reverse geocode latitude: %f, longitude: %f", input.Latitude, input.Longitude)

	var (
		addrs []*geocoding.Address
		err   error
	)

	// Get from cache
	key := cache2.CreateCacheKey(
		reverseGeocodeCache,
		strconv.FormatFloat(input.Latitude, 'g', -1, 64),
		strconv.FormatFloat(input.Longitude, 'g', -1, 64),
	)
	err = s.manager.Get(ctx, key, &addrs)
	if err == nil {
		return mapper.MapAddressesToGeocodingOutput(addrs), nil
	}
	if !errors.Is(err, cache.ErrCacheEntryNotFound) {
		log.Warn().Stack().Err(err).Msg("failed to get reverse geocode from cache")
	}

	// Get from provider
	addrs, err = s.provider.ReverseGeocodeWithContext(
		ctx,
		geocoding.Location{Lat: input.Latitude, Lng: input.Longitude},
		geocoding.ReverseGeocodeConfig{Lang: lang, Limit: input.Limit, Zoom: 18},
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to reverse geocode address")
		if errors.Is(err, chained.ErrGeocodeProvidersUnavailable) {
			return dto.ReverseGeocodingOutput{}, service.ErrGeocodeProvidersUnavailable
		}
		return dto.ReverseGeocodingOutput{}, err
	}

	// Save	to cache
	err = s.manager.Set(ctx, key, addrs)
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to set reverse geocode to cache")
	}

	return mapper.MapAddressesToGeocodingOutput(addrs), nil
}
