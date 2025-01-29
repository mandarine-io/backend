package geocoding

import (
	"context"
	mock1 "github.com/mandarine-io/backend/internal/infrastructure/cache/mock"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/domain/geocoding"
	mock2 "github.com/mandarine-io/backend/third_party/geocoding/mock"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	ctx = context.Background()

	cacheManagerMock      *mock1.ManagerMock
	geocodingProviderMock *mock2.ProviderMock
	svc                   domain.GeocodingService
)

func init() {
	cacheManagerMock = new(mock1.ManagerMock)
	geocodingProviderMock = new(mock2.ProviderMock)
	svc = geocoding.NewService(cacheManagerMock, geocodingProviderMock)
}

type GeocodingServiceSuite struct {
	suite.Suite
}

func TestGeocodingServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(GeocodingServiceSuite))
}

func (s *GeocodingServiceSuite) Test(t provider.T) {
	s.RunSuite(t, new(GeocodeSuite))
	s.RunSuite(t, new(ReverseGeocodeSuite))
}
