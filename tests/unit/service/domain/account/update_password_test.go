package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/util/security"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UpdatePasswordSuite struct {
	suite.Suite
}

func (s *UpdatePasswordSuite) Test_Success(t provider.T) {
	t.Title("Successfully updates the password for a valid user and input")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdatePassword")
	t.Tags("Positive")

	ctx := context.Background()
	userID := uuid.New()
	hashPassword, _ := security.HashPassword("oldpassword")
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: false,
		Password:       hashPassword,
	}
	req := v0.UpdatePasswordInput{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

	err := svc.UpdatePassword(ctx, userID, req)

	t.Require().NoError(err)
}

func (s *UpdatePasswordSuite) Test_UserNotFound(t provider.T) {
	t.Title("Returns error when user is not found")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdatePassword")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	req := v0.UpdatePasswordInput{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	err := svc.UpdatePassword(ctx, userID, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotFound, err)
}

func (s *UpdatePasswordSuite) Test_ErrorFindingUser(t provider.T) {
	t.Title("Returns error when finding user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdatePassword")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("database error")
	req := v0.UpdatePasswordInput{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, expectedErr)

	err := svc.UpdatePassword(ctx, userID, req)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *UpdatePasswordSuite) Test_IncorrectOldPassword(t provider.T) {
	t.Title("Returns error when old password does not match")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdatePassword")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	hashPassword, _ := security.HashPassword("oldpassword")
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: false,
		Password:       hashPassword,
	}
	req := v0.UpdatePasswordInput{
		OldPassword: "wrongoldpassword",
		NewPassword: "newpassword",
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	err := svc.UpdatePassword(ctx, userID, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrIncorrectOldPassword, err)
}

func (s *UpdatePasswordSuite) Test_ErrorHashPassword(t provider.T) {
	t.Title("Returns error when hashing the new password fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdatePassword")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	hashPassword, _ := security.HashPassword("oldpassword")
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: false,
		Password:       hashPassword,
	}
	req := v0.UpdatePasswordInput{
		OldPassword: "oldpassword",
		NewPassword: strings.Repeat("1", 1000),
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	err := svc.UpdatePassword(ctx, userID, req)

	t.Require().Error(err)
	t.Require().Equal(bcrypt.ErrPasswordTooLong, err)
}

func (s *UpdatePasswordSuite) Test_ErrorUpdateUser(t provider.T) {
	t.Title("Returns error when updating user fails")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdatePassword")
	t.Tags("Negative")

	ctx := context.Background()
	userID := uuid.New()
	hashPassword, _ := security.HashPassword("oldpassword")
	userEntity := &entity.User{
		ID:             userID,
		IsPasswordTemp: false,
		Password:       hashPassword,
	}
	req := v0.UpdatePasswordInput{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}
	expectedErr := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

	err := svc.UpdatePassword(ctx, userID, req)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
