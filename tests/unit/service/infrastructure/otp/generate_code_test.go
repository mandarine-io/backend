package otp

import (
	"github.com/mandarine-io/backend/config"
	otp2 "github.com/mandarine-io/backend/internal/service/infrastructure/otp"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type GenerateCodeSuite struct {
	suite.Suite
}

func (s *GenerateCodeSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("OTP service")
	t.Feature("GenerateCode")
	t.Tags("Positive")

	otp, err := svc.GenerateCode(ctx)

	t.Require().NoError(err)
	t.Require().Len(otp, cfg.Length)
	t.Require().Regexp("^\\d*$", otp)
}

func (s *GenerateCodeSuite) Test_ErrLength(t provider.T) {
	t.Title("Returns OTP length error")
	t.Severity(allure.CRITICAL)
	t.Epic("OTP service")
	t.Feature("GenerateCode")
	t.Tags("Negative")

	newCfg := config.OTPConfig{
		TTL:    120,
		Length: -1,
	}
	newSvc := otp2.NewService(managerMock, newCfg)

	_, err := newSvc.GenerateCode(ctx)

	t.Require().Error(err)
	t.Require().ErrorIs(err, otp2.ErrNegativeOTPLength)
}
