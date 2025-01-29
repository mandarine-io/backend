package otp

import (
	"errors"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"time"
)

type SaveWithCodeSuite struct {
	suite.Suite
}

func (s *SaveWithCodeSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("OTP service")
	t.Feature("SaveWithCode")
	t.Tags("Positive")

	prefix := "prefix"
	code := "123456"
	data := "data"

	managerMock.On("SetWithExpiration", ctx, mock.Anything, data, time.Duration(cfg.TTL)*time.Second).Once().Return(nil)

	err := svc.SaveWithCode(ctx, prefix, code, data)

	t.Require().NoError(err)
}

func (s *SaveWithCodeSuite) Test_ErrSettingCache(t provider.T) {
	t.Title("Returns setting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("SaveWithCode")
	t.Tags("Negative")

	prefix := "prefix"
	code := "123456"
	data := "data"
	cacheErr := errors.New("cache error")

	managerMock.On(
		"SetWithExpiration",
		ctx,
		mock.Anything,
		data,
		time.Duration(cfg.TTL)*time.Second,
	).Once().Return(cacheErr)

	err := svc.SaveWithCode(ctx, prefix, code, data)

	t.Require().Error(err)
	t.Require().ErrorIs(err, cacheErr)
}
