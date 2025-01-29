package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"time"
)

type DeleteAccountSuite struct {
	suite.Suite
}

func (s *DeleteAccountSuite) Test_Success(t provider.T) {
	t.Title("Successfully deletes a user account")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("DeleteAccount")
	t.Tags("Positive")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID: userID,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

	err := svc.DeleteAccount(ctx, userID)

	t.Require().NoError(err)

}

func (s *DeleteAccountSuite) Test_UserNotFound(t provider.T) {
	t.Title("Returns error when user is not found")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("DeleteAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	err := svc.DeleteAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotFound, err)

}

func (s *DeleteAccountSuite) Test_ErrorFindingUser(t provider.T) {
	t.Title("Returns error when finding user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("DeleteAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, expectedErr)

	err := svc.DeleteAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)

}

func (s *DeleteAccountSuite) Test_UserAlreadyDeleted(t provider.T) {
	t.Title("Returns error when user account is already deleted")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("DeleteAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	deletedAt := time.Now()
	userEntity := &entity.User{
		ID:        userID,
		DeletedAt: &deletedAt,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	err := svc.DeleteAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserAlreadyDeleted, err)

}

func (s *DeleteAccountSuite) Test_ErrorUpdateUser(t provider.T) {
	t.Title("Returns error when updating user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("DeleteAccount")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:        userID,
		DeletedAt: nil,
	}
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

	err := svc.DeleteAccount(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)

}
