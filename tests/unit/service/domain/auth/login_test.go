package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/util/security"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type LoginSuite struct {
	suite.Suite
}

func (s *LoginSuite) Test_ErrUserNotFound(t provider.T) {
	t.Title("Login returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Login")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.LoginInput{Login: "test@example.com", Password: "password123"}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByUsernameOrEmail", ctx, req.Login, mock.Anything).
		Once().Return(nil, nil)

	resp, err := svc.Login(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotFound, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *LoginSuite) Test_ErrorFindingUser(t provider.T) {
	t.Title("Login returns DB FindingUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Login")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.LoginInput{Login: "test@example.com", Password: "password123"}
	expectedErr := errors.New("database error")

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByUsernameOrEmail", ctx, req.Login, mock.Anything).
		Once().Return(nil, expectedErr)

	resp, err := svc.Login(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *LoginSuite) Test_ErrBadCredentials(t provider.T) {
	t.Title("Login returns BadCredentials error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Login")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.LoginInput{Login: "test@example.com", Password: "password123"}
	userEntity := &entity.User{
		Email:    req.Login,
		Password: "hashedpassword",
	}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByUsernameOrEmail", ctx, req.Login, mock.Anything).
		Once().Return(userEntity, nil)

	resp, err := svc.Login(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrBadCredentials, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *LoginSuite) Test_ErrUserIsBlocked(t provider.T) {
	t.Title("Login returns UserIsBlocked error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("Login")
	t.Tags("Negative")

	ctx := context.Background()
	req := v0.LoginInput{Login: "test@example.com", Password: "password123"}
	hashPassword, _ := security.HashPassword("password123")
	userEntity := &entity.User{
		Email:     req.Login,
		Password:  hashPassword,
		IsEnabled: false,
	}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByUsernameOrEmail", ctx, req.Login, mock.Anything).
		Once().Return(userEntity, nil)

	resp, err := svc.Login(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserIsBlocked, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *LoginSuite) Test_Success(t provider.T) {
	t.Title("Login returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("Login")
	t.Tags("Positive")

	ctx := context.Background()
	req := v0.LoginInput{Login: "test@example.com", Password: "password123"}
	hashPassword, _ := security.HashPassword("password123")
	accessToken := "access_token"
	refreshToken := "refresh_token"
	userEntity := &entity.User{
		Email:     req.Login,
		Password:  hashPassword,
		IsEnabled: true,
	}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByUsernameOrEmail", ctx, req.Login, mock.Anything).
		Once().Return(userEntity, nil)
	jwtServiceMock.On("GenerateTokens", ctx, userEntity).Once().Return(accessToken, refreshToken, nil)

	resp, err := svc.Login(ctx, req)

	t.Require().NoError(err)
	t.Require().Equal(accessToken, resp.AccessToken)
	t.Require().Equal(refreshToken, resp.RefreshToken)
}
