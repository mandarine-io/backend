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
)

type GetAccountSuite struct {
	suite.Suite
}

func (s *GetAccountSuite) Test_Success(t provider.T) {
	t.Title("GetAccount returns account successfully")
	t.Severity(allure.NORMAL)
	t.Epic("Account service")
	t.Feature("GetAccount")
	t.Tags("Positive")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:       userID,
		Username: "test_user",
		Email:    "test@example.com",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Return(userEntity, nil).Once()

	// Act
	result, err := svc.GetAccount(ctx, userID)

	// Assert
	t.Require().NoError(err)
	t.Require().Equal(userEntity.Username, result.Username)
	t.Require().Equal(userEntity.Email, result.Email)

}

func (s *GetAccountSuite) Test_UserNotFound(t provider.T) {
	t.Title("GetAccount returns error when user not found")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("GetAccount")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	userRepoMock.On("FindUserByID", ctx, userID).Return(nil, nil).Once()

	// Act
	_, err := svc.GetAccount(ctx, userID)

	// Assert
	t.Require().Error(err)
	t.Require().ErrorIs(err, domain.ErrUserNotFound)

}

func (s *GetAccountSuite) Test_DbErrorDuringFindUserById(t provider.T) {
	t.Title("GetAccount returns db error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("GetAccount")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("db error")

	userRepoMock.On("FindUserByID", ctx, userID).Return(nil, expectedErr).Once()

	// Act
	_, err := svc.GetAccount(ctx, userID)

	// Assert
	t.Require().Error(err)
	t.Require().ErrorIs(err, expectedErr)

}
