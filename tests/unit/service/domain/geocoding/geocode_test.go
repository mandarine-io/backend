package geocoding

import (
	"errors"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	geocoding2 "github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
)

type GeocodeSuite struct {
	suite.Suite
}

func (s *GeocodeSuite) Test_Success_CacheHit(t provider.T) {
	t.Title("Returns success - cache hit")
	t.Severity(allure.NORMAL)
	t.Epic("Geocoding service")
	t.Feature("Geocode")
	t.Tags("Positive")

	input := v0.GeocodingInput{Address: "address"}
	locs := []*geocoding2.Location{
		{Lat: 1.0, Lng: 1.0},
	}
	lang := language.English

	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			*args.Get(2).(*[]*geocoding2.Location) = locs
		},
	).Return(nil).Once()

	output, err := svc.Geocode(ctx, input, lang)

	t.Require().NoError(err)
	t.Require().Equal(len(locs), output.Count)
	t.Require().Equal(locs[0].Lat, output.Data[0].Latitude)
	t.Require().Equal(locs[0].Lng, output.Data[0].Longitude)
}

func (s *GeocodeSuite) Test_Success_CacheMiss(t provider.T) {
	t.Title("Returns success - cache miss")
	t.Severity(allure.NORMAL)
	t.Epic("Geocoding service")
	t.Feature("Geocode")
	t.Tags("Positive")

	input := v0.GeocodingInput{Address: "address"}
	locs := []*geocoding2.Location{
		{Lat: 1.0, Lng: 1.0},
	}
	lang := language.English

	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
	cacheManagerMock.On("Set", ctx, mock.Anything, mock.Anything).Return(nil).Once()
	geocodingProviderMock.On("Geocode", ctx, input.Address, mock.Anything).Return(locs, nil).Once()

	output, err := svc.Geocode(ctx, input, lang)

	t.Require().NoError(err)
	t.Require().Equal(len(locs), output.Count)
	t.Require().Equal(locs[0].Lat, output.Data[0].Latitude)
	t.Require().Equal(locs[0].Lng, output.Data[0].Longitude)
}

func (s *GeocodeSuite) Test_ErrServicesUnavailable(t provider.T) {
	t.Title("Returns ErrServicesUnavailable")
	t.Severity(allure.CRITICAL)
	t.Epic("Geocoding service")
	t.Feature("Geocode")
	t.Tags("Negative")

	input := v0.GeocodingInput{Address: "address"}
	lang := language.English

	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
	geocodingProviderMock.On("Geocode", ctx, input.Address, mock.Anything).
		Return(nil, geocoding2.ErrGeocodeProvidersUnavailable).Once()

	_, err := svc.Geocode(ctx, input, lang)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrGeocodeProvidersUnavailable, err)
}

func (s *GeocodeSuite) Test_UnexpectedErr(t provider.T) {
	t.Title("Returns UnexpectedErr")
	t.Severity(allure.CRITICAL)
	t.Epic("Geocoding service")
	t.Feature("Geocode")
	t.Tags("Negative")

	input := v0.GeocodingInput{Address: "address"}
	lang := language.English

	expectedErr := errors.New("unexpected error")
	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Return(expectedErr).Once()
	geocodingProviderMock.On("Geocode", ctx, input.Address, mock.Anything).
		Return(nil, expectedErr).Once()

	_, err := svc.Geocode(ctx, input, lang)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
