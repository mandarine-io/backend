package auth

import (
	"context"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/third_party/oauth"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type RegisterOrLoginSuite struct {
	suite.Suite
}

func (s *RegisterOrLoginSuite) Test_SuccessUniqueUsername(t provider.T) {
	t.Title("RegisterOrLogin returns success with unique username")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RegisterOrLogin")
	t.Tags("Positive")

	userInfo := oauth.UserInfo{Username: "test", Email: "test@example.com"}
	userEntity := &entity.User{Email: "test@example.com"}
	accessToken := "access_token"
	refreshToken := "refresh_token"

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, userInfo.Email, mock.Anything).Return(nil, nil).Once()
	userRepoMock.On("CreateUser", mock.Anything, mock.Anything).Return(userEntity, nil).Once()
	userRepoMock.On("ExistsUserByUsername", mock.Anything, userInfo.Username).Return(false, nil).Once()
	jwtServiceMock.On("GenerateTokens", context.Background(), userEntity).Once().Return(accessToken, refreshToken, nil)

	result, err := svc.RegisterOrLogin(context.Background(), userInfo)

	t.Require().NoError(err)
	t.Require().Equal(result.AccessToken, accessToken)
	t.Require().Equal(result.RefreshToken, refreshToken)
}

func (s *RegisterOrLoginSuite) Test_SuccessNotUniqueUsername(t provider.T) {
	t.Title("RegisterOrLogin returns success with not unique username")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RegisterOrLogin")
	t.Tags("Positive")

	userInfo := oauth.UserInfo{Username: "test", Email: "test@example.com"}
	userEntity := &entity.User{Email: "test@example.com"}
	accessToken := "access_token"
	refreshToken := "refresh_token"

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, userInfo.Email, mock.Anything).Return(nil, nil).Once()
	userRepoMock.On("CreateUser", mock.Anything, mock.Anything).Return(userEntity, nil).Once()
	userRepoMock.On("ExistsUserByUsername", mock.Anything, userInfo.Username).Return(true, nil).Once()
	userRepoMock.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, nil).Once()
	jwtServiceMock.On("GenerateTokens", context.Background(), userEntity).Once().Return(accessToken, refreshToken, nil)

	result, err := svc.RegisterOrLogin(context.Background(), userInfo)

	t.Require().NoError(err)
	t.Require().Equal(result.AccessToken, accessToken)
	t.Require().Equal(result.RefreshToken, refreshToken)
}

func (s *RegisterOrLoginSuite) Test_SuccessExistingUser(t provider.T) {
	t.Title("RegisterOrLogin returns success with existent user")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RegisterOrLogin")
	t.Tags("Positive")

	userInfo := oauth.UserInfo{Email: "test@example.com"}
	userEntity := &entity.User{Email: "test@example.com", IsEnabled: true}
	accessToken := "access_token"
	refreshToken := "refresh_token"

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, userInfo.Email, mock.Anything).Return(userEntity, nil).Once()
	jwtServiceMock.On("GenerateTokens", context.Background(), userEntity).Once().Return(accessToken, refreshToken, nil)

	result, err := svc.RegisterOrLogin(context.Background(), userInfo)

	t.Require().NoError(err)
	t.Require().Equal(result.AccessToken, accessToken)
	t.Require().Equal(result.RefreshToken, refreshToken)
}

func (s *RegisterOrLoginSuite) Test_ErrorFindingUser(t provider.T) {
	t.Title("RegisterOrLogin returns DB FindingUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterOrLogin")
	t.Tags("Negative")

	userInfo := oauth.UserInfo{Email: "test@example.com"}
	expectedError := errors.New("repo error")

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, userInfo.Email, mock.Anything).Return(nil, expectedError).Once()

	_, err := svc.RegisterOrLogin(context.Background(), userInfo)

	t.Require().Error(err)
	t.Require().Equal(expectedError, err)
}

func (s *RegisterOrLoginSuite) Test_ErrorExistingUser(t provider.T) {
	t.Title("RegisterOrLogin returns DB ExistingUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterOrLogin")
	t.Tags("Negative")

	userInfo := oauth.UserInfo{Email: "test@example.com"}
	expectedError := errors.New("repo error")

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, userInfo.Email, mock.Anything).Return(nil, nil).Once()
	userRepoMock.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, expectedError).Once()

	_, err := svc.RegisterOrLogin(context.Background(), userInfo)

	t.Require().Error(err)
	t.Require().Equal(expectedError, err)
}

func (s *RegisterOrLoginSuite) Test_ErrorCreatingUser(t provider.T) {
	t.Title("RegisterOrLogin returns DB CreateUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RegisterOrLogin")
	t.Tags("Negative")

	userInfo := oauth.UserInfo{Email: "test@example.com"}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByEmail", mock.Anything, userInfo.Email, mock.Anything).Return(nil, nil).Once()
	userRepoMock.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, nil).Once()
	userRepoMock.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("create error")).Once()

	_, err := svc.RegisterOrLogin(context.Background(), userInfo)

	t.Require().Error(err)
	t.Require().Equal("create error", err.Error())
}
