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
)

type UpdateUsernameSuite struct {
	suite.Suite
}

func (s *UpdateUsernameSuite) Test_Success(t provider.T) {
	t.Title("UpdateUsername returns account successfully")
	t.Severity(allure.NORMAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Positive")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Username:        "old-username",
		IsEmailVerified: true,
	}

	userRepoMock.On("ExistsUserByUsername", ctx, "username").Once().Return(false, nil)
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "username",
	}
	resp, err := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().NoError(err)
	t.Require().Equal("username", resp.Username)
}

func (s *UpdateUsernameSuite) Test_UserNotFound(t provider.T) {
	t.Title("UpdateUsername returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, nil)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "username",
	}
	resp, err := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().Equal(domain.ErrUserNotFound, err)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateUsernameSuite) Test_DbErrorDuringFindUserById(t provider.T) {
	t.Title("UpdateUsername returns DB error during FindUserById")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	err := errors.New("database error")

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(nil, err)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "username",
	}
	resp, err1 := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().Equal(err, err1)
	t.Require().Equal(v0.AccountOutput{}, resp)
}

func (s *UpdateUsernameSuite) Test_UsernameNotChanged(t provider.T) {
	t.Title("UpdateUsername returns old username because username not changed")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Username:        "old-username",
		IsEmailVerified: true,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "old-username",
	}
	resp, err := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().NoError(err)
	t.Require().Equal("old-username", resp.Username)
}

func (s *UpdateUsernameSuite) Test_ErrDuplicateUsername(t provider.T) {
	t.Title("UpdateUsername returns ErrDuplicateUsername")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Username:        "old-username",
		IsEmailVerified: true,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("ExistsUserByUsername", ctx, "username").Once().Return(true, nil)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "username",
	}
	resp, err := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().Equal(v0.AccountOutput{}, resp)
	t.Require().Equal(domain.ErrDuplicateUsername, err)
}

func (s *UpdateUsernameSuite) Test_DbErrDuringExistsUserByUsername(t provider.T) {
	t.Title("UpdateUsername returns DB error during ExistsUserByUsername")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	err := errors.New("database error")

	userEntity := &entity.User{
		ID:              userID,
		Username:        "old-username",
		IsEmailVerified: true,
	}

	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("ExistsUserByUsername", ctx, "username").Once().Return(true, err)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "username",
	}
	resp, err1 := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().Equal(v0.AccountOutput{}, resp)
	t.Require().Equal(err, err1)
}

func (s *UpdateUsernameSuite) Test_DbErrDuringUpdateUser(t provider.T) {
	t.Title("UpdateUsername returns DB error during UpdateUser")
	t.Severity(allure.CRITICAL)
	t.Epic("Account service")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	userEntity := &entity.User{
		ID:              userID,
		Username:        "old-username",
		IsEmailVerified: true,
	}

	err := errors.New("database error")
	userRepoMock.On("FindUserByID", ctx, userID).Once().Return(userEntity, nil)
	userRepoMock.On("ExistsUserByUsername", ctx, "username").Once().Return(false, nil)
	userRepoMock.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, err)

	// Act
	req := v0.UpdateUsernameInput{
		Username: "username",
	}
	resp, err1 := svc.UpdateUsername(ctx, userID, req)

	// Assert
	t.Require().Equal(err, err1)
	t.Require().Equal(v0.AccountOutput{}, resp)
}
