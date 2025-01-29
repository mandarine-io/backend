package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"time"
)

type RestoreAccountSuite struct {
	suite.Suite
}

func (s *RestoreAccountSuite) Test_Success(t provider.T) {
	t.Title("Successfully restores a deleted account")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("RestoreAccount")
	t.Tags("Positive")

	ctx := context.Background()
	userID := uuid.New()
	deletedAt := time.Now()
	userEntity := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		IsEnabled: true,
		DeletedAt: &deletedAt,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

	resp, err := svc.RestoreAccount(ctx, userID)

	t.Require().NoError(err)
	t.Require().Equal(userEntity.Email, resp.Email)
	t.Require().Equal(userEntity.IsEnabled, resp.IsEnabled)

}

func (s *RestoreAccountSuite) Test_UserNotFound(t provider.T) {
	t.Title("Returns error when user is not found")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("RestoreAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	resp, err := svc.RestoreAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotFound, err)
	t.Require().Equal(v0.AccountOutput{}, resp)

}

func (s *RestoreAccountSuite) Test_ErrorFindingUser(t provider.T) {
	t.Title("Returns error when finding user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("RestoreAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, expectedErr)

	resp, err := svc.RestoreAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
	t.Require().Equal(v0.AccountOutput{}, resp)

}

func (s *RestoreAccountSuite) Test_UserNotDeleted(t provider.T) {
	t.Title("Returns error when user account is not marked as deleted")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("RestoreAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:        userID,
		DeletedAt: nil,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	resp, err := svc.RestoreAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotDeleted, err)
	t.Require().Equal(v0.AccountOutput{}, resp)

}

func (s *RestoreAccountSuite) Test_ErrorUpdateUser(t provider.T) {
	t.Title("Returns error when updating user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("RestoreAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	deletedAt := time.Now()
	userEntity := &entity.User{
		ID:        userID,
		DeletedAt: &deletedAt,
	}
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

	resp, err := svc.RestoreAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
	t.Require().Equal(v0.AccountOutput{}, resp)

}
