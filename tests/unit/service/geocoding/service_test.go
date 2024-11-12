package auth_test

import (
	"context"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/geocoding"
	geocoding2 "github.com/mandarine-io/Backend/pkg/geocoding"
	mock2 "github.com/mandarine-io/Backend/pkg/geocoding/mock"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	mock1 "github.com/mandarine-io/Backend/pkg/storage/cache/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"testing"
)

var (
	ctx = context.Background()

	cacheManager = new(mock1.ManagerMock)
	provider     = new(mock2.ProviderMock)
	svc          = geocoding.NewService([]geocoding2.Provider{provider}, cacheManager)
)

func Test_GeocodingService_Geocode(t *testing.T) {
	input := dto.GeocodingInput{Address: "address"}
	locs := []*geocoding2.Location{
		{Lat: 1.0, Lng: 1.0},
	}
	lang := language.English

	t.Run("Geocoding services unavailable", func(t *testing.T) {
		cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
		provider.On("GeocodeWithContext", ctx, input.Address, mock.Anything).
			Return(nil, geocoding2.ErrGeocodeProvidersUnavailable).Once()

		_, err := svc.Geocode(ctx, input, lang)

		assert.Error(t, err)
		assert.Equal(t, service.ErrGeocodeProvidersUnavailable, err)
	})

	t.Run("Unexpected error", func(t *testing.T) {
		expectedErr := errors.New("unexpected error")
		cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Return(expectedErr).Once()
		provider.On("GeocodeWithContext", ctx, input.Address, mock.Anything).
			Return(nil, expectedErr).Once()

		_, err := svc.Geocode(ctx, input, lang)

		assert.Error(t, err)
		assert.Equal(t, service.ErrGeocodeProvidersUnavailable, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Run("Cache hit", func(t *testing.T) {
			cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				*args.Get(2).(*[]*geocoding2.Location) = locs
			}).Return(nil).Once()

			output, err := svc.Geocode(ctx, input, lang)

			assert.NoError(t, err)
			assert.Equal(t, len(locs), output.Count)
			assert.Equal(t, locs[0].Lat, output.Data[0].Latitude)
			assert.Equal(t, locs[0].Lng, output.Data[0].Longitude)
		})

		t.Run("Cache miss", func(t *testing.T) {
			cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
			cacheManager.On("Set", ctx, mock.Anything, mock.Anything).Return(nil).Once()
			provider.On("GeocodeWithContext", ctx, input.Address, mock.Anything).Return(locs, nil).Once()

			output, err := svc.Geocode(ctx, input, lang)

			assert.NoError(t, err)
			assert.Equal(t, len(locs), output.Count)
			assert.Equal(t, locs[0].Lat, output.Data[0].Latitude)
			assert.Equal(t, locs[0].Lng, output.Data[0].Longitude)
		})
	})
}

func Test_ReverseGeocodingService_Geocode(t *testing.T) {
	input := dto.ReverseGeocodingInput{Point: "1.0,1.0"}
	addrs := []*geocoding2.Address{
		{FormattedAddress: "address"},
	}
	lang := language.English

	t.Run("Geocoding services unavailable", func(t *testing.T) {
		cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
		provider.On("ReverseGeocodeWithContext", ctx, mock.Anything, mock.Anything).
			Return(nil, geocoding2.ErrGeocodeProvidersUnavailable).Once()

		_, err := svc.ReverseGeocode(ctx, input, lang)

		assert.Error(t, err)
		assert.Equal(t, service.ErrGeocodeProvidersUnavailable, err)
	})

	t.Run("Unexpected error", func(t *testing.T) {
		expectedErr := errors.New("unexpected error")
		cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Return(expectedErr).Once()
		provider.On("ReverseGeocodeWithContext", ctx, mock.Anything, mock.Anything).
			Return(nil, expectedErr).Once()

		_, err := svc.ReverseGeocode(ctx, input, lang)

		assert.Error(t, err)
		assert.Equal(t, service.ErrGeocodeProvidersUnavailable, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Run("Cache hit", func(t *testing.T) {
			cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				*args.Get(2).(*[]*geocoding2.Address) = addrs
			}).Return(nil).Once()

			output, err := svc.ReverseGeocode(ctx, input, lang)

			assert.NoError(t, err)
			assert.Equal(t, len(addrs), output.Count)
			assert.Equal(t, addrs[0].FormattedAddress, output.Data[0].FormattedAddress)
		})

		t.Run("Cache miss", func(t *testing.T) {
			cacheManager.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
			cacheManager.On("Set", ctx, mock.Anything, mock.Anything).Return(nil).Once()
			provider.On("ReverseGeocodeWithContext", ctx, mock.Anything, mock.Anything).Return(addrs, nil).Once()

			output, err := svc.ReverseGeocode(ctx, input, lang)

			assert.NoError(t, err)
			assert.Equal(t, len(addrs), output.Count)
			assert.Equal(t, addrs[0].FormattedAddress, output.Data[0].FormattedAddress)
		})
	})
}
