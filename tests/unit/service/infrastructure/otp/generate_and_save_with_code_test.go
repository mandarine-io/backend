package otp

import (
	"errors"
	"github.com/mandarine-io/backend/config"
	otp2 "github.com/mandarine-io/backend/internal/service/infrastructure/otp"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"time"
)

type GenerateAndSaveWithCodeSuite struct {
	suite.Suite
}

func (s *GenerateAndSaveWithCodeSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("OTP service")
	t.Feature("GenerateAndSaveWithCode")
	t.Tags("Positive")

	prefix := "prefix"
	data := "data"
	managerMock.On("SetWithExpiration", ctx, mock.Anything, data, time.Duration(cfg.TTL)*time.Second).Once().Return(nil)

	otp, err := svc.GenerateAndSaveWithCode(ctx, prefix, data)

	t.Require().NoError(err)
	t.Require().Len(otp, cfg.Length)
	t.Require().Regexp("^\\d*$", otp)
}

func (s *GenerateAndSaveWithCodeSuite) Test_ErrLength(t provider.T) {
	t.Title("Returns OTP length error")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("GenerateAndSaveWithCode")
	t.Tags("Negative")

	newCfg := config.OTPConfig{
		TTL:    120,
		Length: -1,
	}
	newSvc := otp2.NewService(managerMock, newCfg)

	_, err := newSvc.GenerateAndSaveWithCode(ctx, "", "")

	t.Require().Error(err)
	t.Require().ErrorIs(err, otp2.ErrNegativeOTPLength)
}

func (s *GenerateAndSaveWithCodeSuite) Test_ErrSettingCache(t provider.T) {
	t.Title("Returns setting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("GenerateAndSaveWithCode")
	t.Tags("Negative")

	prefix := "prefix"
	data := "data"
	cacheErr := errors.New("cache error")

	managerMock.On(
		"SetWithExpiration",
		ctx,
		mock.Anything,
		data,
		time.Duration(cfg.TTL)*time.Second,
	).Once().Return(cacheErr)

	_, err := svc.GenerateAndSaveWithCode(ctx, prefix, data)

	t.Require().Error(err)
	t.Require().ErrorIs(err, cacheErr)
}
