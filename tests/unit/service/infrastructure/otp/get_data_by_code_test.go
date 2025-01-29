package otp

import (
	"errors"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	infra "github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type GetDataByCodeSuite struct {
	suite.Suite
}

func (s *GetDataByCodeSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("OTP service")
	t.Feature("GetDataByCode")
	t.Tags("Positive")

	prefix := "prefix"
	code := "123456"
	data := ""
	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			dataPtr := args.Get(2).(*string)
			*dataPtr = "data"
		},
	).Once().Return(nil)

	err := svc.GetDataByCode(ctx, prefix, code, &data)

	t.Require().NoError(err)
	t.Require().Equal("data", data)
}

func (s *GetDataByCodeSuite) Test_CacheEntryNotFound(t provider.T) {
	t.Title("Returns cache entry not found")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("GetDataByCode")
	t.Tags("Negative")

	prefix := "prefix"
	code := "123456"
	data := ""

	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Return(cache.ErrCacheEntryNotFound)

	err := svc.GetDataByCode(ctx, prefix, code, &data)

	t.Require().Error(err)
	t.Require().ErrorIs(err, infra.ErrInvalidOrExpiredOTP)
}

func (s *GetDataByCodeSuite) Test_ErrGettingCache(t provider.T) {
	t.Title("Returns getting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("GetDataByCode")
	t.Tags("Negative")

	prefix := "prefix"
	code := "123456"
	data := ""
	cacheErr := errors.New("cache error")

	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Return(cacheErr)

	err := svc.GetDataByCode(ctx, prefix, code, &data)

	t.Require().Error(err)
	t.Require().ErrorIs(err, cacheErr)
}
