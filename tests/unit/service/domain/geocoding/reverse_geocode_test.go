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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
)

type ReverseGeocodeSuite struct {
	suite.Suite
}

func (s *ReverseGeocodeSuite) Test_Success_CacheHit(t provider.T) {
	t.Title("Returns success - cache hit")
	t.Severity(allure.NORMAL)
	t.Epic("Geocoding service")
	t.Feature("ReverseGeocode")
	t.Tags("Positive")

	input := v0.ReverseGeocodingInput{
		Lat: decimal.NewFromFloat(1),
		Lng: decimal.NewFromFloat(1),
	}
	addrs := []*geocoding2.Address{
		{FormattedAddress: "address"},
	}
	lang := language.English

	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			*args.Get(2).(*[]*geocoding2.Address) = addrs
		},
	).Return(nil).Once()

	output, err := svc.ReverseGeocode(ctx, input, lang)

	t.Require().NoError(err)
	t.Require().Equal(len(addrs), output.Count)
	t.Require().Equal(addrs[0].FormattedAddress, output.Data[0].FormattedAddress)
}

func (s *ReverseGeocodeSuite) Test_Success_CacheMiss(t provider.T) {
	t.Title("Returns success - cache miss")
	t.Severity(allure.NORMAL)
	t.Epic("Geocoding service")
	t.Feature("ReverseGeocode")
	t.Tags("Positive")

	input := v0.ReverseGeocodingInput{
		Lat: decimal.NewFromFloat(1),
		Lng: decimal.NewFromFloat(1),
	}
	addrs := []*geocoding2.Address{
		{FormattedAddress: "address"},
	}
	lang := language.English

	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
	cacheManagerMock.On("Set", ctx, mock.Anything, mock.Anything).Return(nil).Once()
	geocodingProviderMock.On("ReverseGeocode", ctx, mock.Anything, mock.Anything).Return(addrs, nil).Once()

	output, err := svc.ReverseGeocode(ctx, input, lang)

	t.Require().NoError(err)
	t.Require().Equal(len(addrs), output.Count)
	t.Require().Equal(addrs[0].FormattedAddress, output.Data[0].FormattedAddress)
}

func (s *ReverseGeocodeSuite) Test_ErrServicesUnavailable(t provider.T) {
	t.Title("Returns ErrServicesUnavailable")
	t.Severity(allure.CRITICAL)
	t.Epic("Geocoding service")
	t.Feature("ReverseGeocode")
	t.Tags("Negative")

	input := v0.ReverseGeocodingInput{
		Lat: decimal.NewFromFloat(1),
		Lng: decimal.NewFromFloat(1),
	}
	lang := language.English

	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Return(cache.ErrCacheEntryNotFound).Once()
	geocodingProviderMock.On("ReverseGeocode", ctx, mock.Anything, mock.Anything).
		Return(nil, geocoding2.ErrGeocodeProvidersUnavailable).Once()

	_, err := svc.ReverseGeocode(ctx, input, lang)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrGeocodeProvidersUnavailable, err)
}

func (s *ReverseGeocodeSuite) Test_UnexpectedErr(t provider.T) {
	t.Title("Returns UnexpectedErr")
	t.Severity(allure.CRITICAL)
	t.Epic("Geocoding service")
	t.Feature("ReverseGeocode")
	t.Tags("Negative")

	input := v0.ReverseGeocodingInput{
		Lat: decimal.NewFromFloat(1),
		Lng: decimal.NewFromFloat(1),
	}
	lang := language.English

	expectedErr := errors.New("unexpected error")
	cacheManagerMock.On("Get", ctx, mock.Anything, mock.Anything).Return(expectedErr).Once()
	geocodingProviderMock.On("ReverseGeocode", ctx, mock.Anything, mock.Anything).
		Return(nil, expectedErr).Once()

	_, err := svc.ReverseGeocode(ctx, input, lang)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
