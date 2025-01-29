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
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type SetPasswordSuite struct {
	suite.Suite
}

func (s *SetPasswordSuite) Test_Success(t provider.T) {
	t.Title("SetPassword returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Account service")
	t.Feature("SetPassword")
	t.Tags("Positive")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: true,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

	// Act
	req := v0.SetPasswordInput{
		Password: "newpassword",
	}
	err := svc.SetPassword(ctx, userID, req)

	// Assert
	t.Require().NoError(err)

}

func (s *SetPasswordSuite) Test_UserNotFound(t provider.T) {
	t.Title("SetPassword returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("SetPassword")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	// Act
	req := v0.SetPasswordInput{
		Password: "newpassword",
	}
	err := svc.SetPassword(ctx, userID, req)

	// Assert
	t.Require().Error(err)
	t.Require().ErrorIs(err, domain.ErrUserNotFound)

}

func (s *SetPasswordSuite) Test_DbErrorDuringFindUserById(t provider.T) {
	t.Title("SetPassword returns DB error during FindUserById")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("SetPassword")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, expectedErr)

	// Act
	req := v0.SetPasswordInput{
		Password: "newpassword",
	}
	err := svc.SetPassword(ctx, userID, req)

	// Arrange
	t.Require().Error(err)
	t.Require().ErrorIs(err, expectedErr)

}

func (s *SetPasswordSuite) Test_PasswordAlreadySet(t provider.T) {
	t.Title("SetPassword returns PasswordAlreadySet error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("SetPassword")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: false,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	// Act
	req := v0.SetPasswordInput{
		Password: "newpassword",
	}

	err := svc.SetPassword(ctx, userID, req)

	// Assert
	t.Require().Error(err)
	t.Require().ErrorIs(err, domain.ErrPasswordIsSet)

}

func (s *SetPasswordSuite) Test_ErrorHashPassword(t provider.T) {
	t.Title("SetPassword returns HashPassword error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("SetPassword")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: true,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	// Act
	req := v0.SetPasswordInput{
		Password: strings.Repeat("1", 1000),
	}
	err := svc.SetPassword(ctx, userID, req)

	// Assert
	t.Require().Error(err)
	t.Require().ErrorIs(err, bcrypt.ErrPasswordTooLong)

}

func (s *SetPasswordSuite) Test_DbErrorDuringUpdateUser(t provider.T) {
	t.Title("SetPassword returns DB error during UpdateUser")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("SetPassword")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: true,
	}
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

	// Act
	req := v0.SetPasswordInput{
		Password: "newpassword",
	}
	err := svc.SetPassword(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

}
