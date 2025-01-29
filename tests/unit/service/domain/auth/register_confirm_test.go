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
)

type RegisterConfirmSuite struct {
	suite.Suite
}

func (s *RegisterConfirmSuite) Test_Success(t provider.T) {
	t.Title("RegisterConfirm returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Positive")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}

	userEntity := &entity.User{Email: "test@example.com", Username: "testuser"}
	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("ExistsUserByUsernameOrEmail", ctx, "", "test@example.com").Once().Return(false, nil)
	userRepoMock.On("CreateUser", mock.Anything, mock.Anything).Once().Return(userEntity, nil)
	otpServiceMock.On("DeleteDataByCode", ctx, "register", "123456").Once().Return(nil)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().NoError(err)
}

func (s *RegisterConfirmSuite) Test_InvalidOTP(t provider.T) {
	t.Title("RegisterConfirm returns InvalidOTP error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	otpServiceMock.On(
		"GetDataByCode",
		ctx,
		"register",
		"123456",
		mock.Anything,
	).Once().Return(infrastructure.ErrInvalidOrExpiredOTP)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}

func (s *RegisterConfirmSuite) Test_InvalidEmail(t provider.T) {
	t.Title("RegisterConfirm returns InvalidEmail error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "invalid@example.com"
		},
	).Once().Return(nil)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(infrastructure.ErrInvalidOrExpiredOTP, err)
}

func (s *RegisterConfirmSuite) Test_ExistsUser(t provider.T) {
	t.Title("RegisterConfirm returns ExistsUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("ExistsUserByUsernameOrEmail", ctx, "", "test@example.com").Once().Return(true, nil)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrDuplicateUser, err)
}

func (s *RegisterConfirmSuite) Test_ErrorExistsUser(t provider.T) {
	t.Title("RegisterConfirm returns DB ExistsUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	dbError := errors.New("cache error")
	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("ExistsUserByUsernameOrEmail", ctx, "", "test@example.com").Once().Return(false, dbError)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(dbError, err)
}

func (s *RegisterConfirmSuite) Test_DuplicateUser(t provider.T) {
	t.Title("RegisterConfirm returns DuplicateUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("ExistsUserByUsernameOrEmail", ctx, "", "test@example.com").Once().Return(false, nil)
	userRepoMock.On("CreateUser", ctx, mock.Anything).Once().Return(nil, repo.ErrDuplicateUser)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrDuplicateUser, err)
}

func (s *RegisterConfirmSuite) Test_ErrorSavingUser(t provider.T) {
	t.Title("RegisterConfirm returns DB SavingUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("ExistsUserByUsernameOrEmail", ctx, "", "test@example.com").Once().Return(false, nil)
	userRepoMock.On("CreateUser", mock.Anything, mock.Anything).Once().Return(nil, errors.New("db error"))

	err := svc.RegisterConfirm(ctx, req)

	t.Require().Error(err)
	t.Require().Equal("db error", err.Error())
}

func (s *RegisterConfirmSuite) Test_WarnDeleteCache(t provider.T) {
	t.Title("RegisterConfirm returns ErrDeleteCache error")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RegisterConfirm")
	t.Tags("Positive")

	ctx := context.Background()
	req := v0.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}
	userEntity := &entity.User{Email: "test@example.com", Username: "testuser"}
	cacheErr := errors.New("cache error")

	otpServiceMock.On("GetDataByCode", ctx, "register", "123456", mock.Anything).Run(
		func(args mock.Arguments) {
			args.Get(3).(*v0.RegisterInput).Email = "test@example.com"
		},
	).Once().Return(nil)
	userRepoMock.On("ExistsUserByUsernameOrEmail", ctx, "", "test@example.com").Once().Return(false, nil)
	userRepoMock.On("CreateUser", mock.Anything, mock.Anything).Once().Return(userEntity, nil)
	otpServiceMock.On("DeleteDataByCode", ctx, "register", "123456").Once().Return(cacheErr)

	err := svc.RegisterConfirm(ctx, req)

	t.Require().NoError(err)
}
