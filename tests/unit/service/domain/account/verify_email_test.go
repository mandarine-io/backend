package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type VerifyEmailSuite struct {
	suite.Suite
}

func (s *VerifyEmailSuite) Test_VerifyEmail_Success(t provider.T) {
	t.Title("Successfully verifies email")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Positive")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)
	otpServiceMock.On("DeleteDataByCode", ctx, "email_verify", "123456").Once().Return(nil)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().NoError(err)

}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrCacheEntryNotFound(t provider.T) {
	t.Title("Returns error when cache entry is not found")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Once().Return(nil)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrGetCache(t provider.T) {
	t.Title("Returns error when GetCache fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "654321",
	}
	cacheError := errors.New("cache error")
	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "654321", mock.Anything).Once().Return(cacheError)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().Equal(cacheError, err)
}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrUserNotFound(t provider.T) {
	t.Title("Returns error when user is not found")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().Equal(domain.ErrUserNotFound, err)

}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrFindUserById(t provider.T) {
	t.Title("Returns error when FindUserById fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}
	findUserError := errors.New("database error")

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, findUserError)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().Equal(findUserError, err)

}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrCheckEmail(t provider.T) {
	t.Title("Returns error when email in cache and user do not match")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:    userID,
		Email: "another@example.com",
	}
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)

}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrUpdateUser(t provider.T) {
	t.Title("Returns error when updating user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}
	updateUserError := errors.New("database error")

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(nil, updateUserError)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().Equal(updateUserError, err)

}

func (s *VerifyEmailSuite) Test_VerifyEmail_ErrDeleteInCache(t provider.T) {
	t.Title("Logs error when cache deletion fails but still succeeds")
	t.Severity(allure.NORMAL)
	t.Epic("Account service")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Email:           "test@example.com",
		IsEmailVerified: true,
	}
	req := v0.VerifyEmailInput{
		Email: "test@example.com",
		OTP:   "123456",
	}
	cacheError := errors.New("cache error")

	otpServiceMock.On("GetDataByCode", ctx, "email_verify", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			email := args.Get(3).(*string)
			*email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)
	otpServiceMock.On("DeleteDataByCode", ctx, "email_verify", "123456").Once().Return(cacheError)

	err := svc.VerifyEmail(ctx, userID, req)

	t.Require().NoError(err)
}
