package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ResetPasswordSuite struct {
	suite.Suite
}

func (s *ResetPasswordSuite) Test_Success(t provider.T) {
	t.Title("ResetPassword returns Success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("ResetPassword")
	t.Tags("Positive")

	input := v0.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
	userEntity := &entity.User{Email: "test@example.com"}

	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = input.Email
		},
	).Once().Return(nil)
	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, input.Email, mock.Anything).Return(userEntity, nil).Once()
	userRepoMock.On("UpdateUser", mock.Anything, userEntity).Return(userEntity, nil).Once()
	otpServiceMock.On("DeleteDataByCode", mock.Anything, "recovery_password", "123456").Once().Return(nil)

	err := svc.ResetPassword(context.Background(), input)

	t.Require().NoError(err)
}

func (s *ResetPasswordSuite) Test_InvalidOrExpiredOtp(t provider.T) {
	t.Title("ResetPassword returns InvalidOrExpiredOtp error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("ResetPassword")
	t.Tags("Negative")

	input := v0.ResetPasswordInput{Email: "test@example.com", OTP: "wrong", Password: "newpassword"}
	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "wrong", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = input.Email
		},
	).Once().Return(infrastructure.ErrInvalidOrExpiredOTP)

	err := svc.ResetPassword(context.Background(), input)

	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}

func (s *ResetPasswordSuite) Test_CacheEntryNotFound(t provider.T) {
	t.Title("ResetPassword returns CacheEntryNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("ResetPassword")
	t.Tags("Negative")

	input := v0.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "123456", mock.Anything).Once().Return(nil)

	err := svc.ResetPassword(context.Background(), input)

	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}

func (s *ResetPasswordSuite) Test_UserNotFound(t provider.T) {
	t.Title("ResetPassword returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("ResetPassword")
	t.Tags("Negative")

	input := v0.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}

	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = input.Email
		},
	).Once().Return(nil)
	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, input.Email, mock.Anything).Return(nil, nil).Once()

	err := svc.ResetPassword(context.Background(), input)

	t.Require().Equal(domain.ErrUserNotFound, err)
}

func (s *ResetPasswordSuite) Test_ErrorUpdatingUser(t provider.T) {
	t.Title("ResetPassword returns DB UpdatingUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("ResetPassword")
	t.Tags("Negative")

	input := v0.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
	userEntity := &entity.User{Email: "test@example.com"}

	otpServiceMock.On("GetDataByCode", mock.Anything, "recovery_password", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = input.Email
		},
	).Once().Return(nil)
	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, input.Email, mock.Anything).Return(userEntity, nil).Once()
	userRepoMock.On("UpdateUser", mock.Anything, userEntity).
		Return(userEntity, errors.New("update error")).Once()

	err := svc.ResetPassword(context.Background(), input)

	t.Require().Error(err)
}
