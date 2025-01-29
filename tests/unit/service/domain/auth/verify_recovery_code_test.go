package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
)

type VerifyRecoveryCodeSuite struct {
	suite.Suite
}

func (s *VerifyRecoveryCodeSuite) Test_Success(t provider.T) {
	t.Title("VerifyRecoveryCode returns Success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("VerifyRecoveryCode")
	t.Tags("Positive")

	input := v0.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}
	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = input.Email
		},
	).Once().Return(nil)

	err := svc.VerifyRecoveryCode(context.Background(), input)

	t.Require().NoError(err)
}

func (s *VerifyRecoveryCodeSuite) Test_InvalidOrExpiredOtp(t provider.T) {
	t.Title("VerifyRecoveryCode returns InvalidOrExpiredOtp error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("VerifyRecoveryCode")
	t.Tags("Negative")

	input := v0.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "wrong"}
	otpServiceMock.On(
		"GetDataByCode",
		mock.Anything,
		"recovery_password",
		"wrong",
		mock.Anything,
	).Once().Return(infrastructure.ErrInvalidOrExpiredOTP)

	err := svc.VerifyRecoveryCode(context.Background(), input)

	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}

func (s *VerifyRecoveryCodeSuite) Test_CacheEntryNotFound(t provider.T) {
	t.Title("VerifyRecoveryCode returns CacheEntryNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("VerifyRecoveryCode")
	t.Tags("Negative")

	input := v0.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}
	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "123456", mock.Anything).Once().Return(nil)

	err := svc.VerifyRecoveryCode(context.Background(), input)

	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}
