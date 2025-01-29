package auth

import (
	"context"
	"github.com/google/uuid"
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
	"gorm.io/gorm"
	"time"
)

type RefreshTokensSuite struct {
	suite.Suite
}

func (s *RefreshTokensSuite) Test_InvalidJwtToken(t provider.T) {
	t.Title("RefreshTokens returns InvalidJwtToken error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RefreshTokens")
	t.Tags("Negative")

	ctx := context.Background()
	jwtServiceMock.On("GetRefreshTokenClaims", ctx, "invalid_refresh_token").
		Once().Return(infrastructure.RefreshTokenClaims{}, infrastructure.ErrInvalidJWTToken)

	req := v0.RefreshTokensInput{
		RefreshToken: "invalid_refresh_token",
	}
	resp, err := svc.RefreshTokens(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(infrastructure.ErrInvalidJWTToken, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *RefreshTokensSuite) Test_UserNotFound(t provider.T) {
	t.Title("RefreshTokens returns UserNotFound error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RefreshTokens")
	t.Tags("Negative")

	ctx := context.Background()
	claims := infrastructure.RefreshTokenClaims{
		UserID: uuid.New(),
		JTI:    uuid.New().String(),
		Exp:    time.Now().Unix(),
	}
	jwtServiceMock.On("GetRefreshTokenClaims", ctx, "refreshToken").Once().Return(claims, nil)

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByID", ctx, mock.Anything, mock.Anything).Once().Return(nil, nil)

	req := v0.RefreshTokensInput{
		RefreshToken: "refreshToken",
	}
	resp, err := svc.RefreshTokens(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotFound, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *RefreshTokensSuite) Test_ErrFindingUser(t provider.T) {
	t.Title("RefreshTokens returns DB FindingUser error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RefreshTokens")
	t.Tags("Negative")

	ctx := context.Background()
	expectedErr := errors.New("database error")
	claims := infrastructure.RefreshTokenClaims{
		UserID: uuid.New(),
		JTI:    uuid.New().String(),
		Exp:    time.Now().Unix(),
	}
	jwtServiceMock.On("GetRefreshTokenClaims", ctx, "refreshToken").Once().Return(claims, nil)

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByID", ctx, mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

	req := v0.RefreshTokensInput{
		RefreshToken: "refreshToken",
	}
	resp, err := svc.RefreshTokens(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *RefreshTokensSuite) Test_UserIsBanned(t provider.T) {
	t.Title("RefreshTokens returns UserIsBanned error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RefreshTokens")
	t.Tags("Negative")

	ctx := context.Background()
	userId := uuid.New()
	claims := infrastructure.RefreshTokenClaims{
		UserID: userId,
		JTI:    uuid.New().String(),
		Exp:    time.Now().Unix(),
	}
	userEntity := &entity.User{
		IsEnabled: false,
	}
	jwtServiceMock.On("GetRefreshTokenClaims", ctx, "refreshToken").Once().Return(claims, nil)

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByID", ctx, mock.Anything, mock.Anything).Once().Return(userEntity, nil)

	req := v0.RefreshTokensInput{
		RefreshToken: "refreshToken",
	}
	resp, err := svc.RefreshTokens(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserIsBlocked, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *RefreshTokensSuite) Test_ErrGenerateTokens(t provider.T) {
	t.Title("RefreshTokens returns GenerateTokens error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth service")
	t.Feature("RefreshTokens")
	t.Tags("Negative")

	ctx := context.Background()
	userId := uuid.New()
	expectedErr := errors.New("jwt error")
	claims := infrastructure.RefreshTokenClaims{
		UserID: userId,
		JTI:    uuid.New().String(),
		Exp:    time.Now().Unix(),
	}
	userEntity := &entity.User{
		IsEnabled: true,
	}
	jwtServiceMock.On("GetRefreshTokenClaims", ctx, "refreshToken").Once().Return(claims, nil)

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByID", ctx, mock.Anything, mock.Anything).Once().Return(userEntity, nil)
	jwtServiceMock.On("GenerateTokens", ctx, userEntity).Once().Return("", "", expectedErr)

	req := v0.RefreshTokensInput{
		RefreshToken: "refreshToken",
	}
	resp, err := svc.RefreshTokens(ctx, req)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
	t.Require().Equal(v0.JwtTokensOutput{}, resp)
}

func (s *RefreshTokensSuite) Test_Success(t provider.T) {
	t.Title("RefreshTokens returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Auth service")
	t.Feature("RefreshTokens")
	t.Tags("Positive")

	ctx := context.Background()
	userId := uuid.New()
	accessToken := "access_token"
	refreshToken := "refresh_token"
	claims := infrastructure.RefreshTokenClaims{
		UserID: userId,
		JTI:    uuid.New().String(),
		Exp:    time.Now().Unix(),
	}
	userEntity := &entity.User{
		IsEnabled: true,
	}
	jwtServiceMock.On("GetRefreshTokenClaims", ctx, "refreshToken").Once().Return(claims, nil)

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	userRepoMock.On("WithRolePreload").Once().Return(scope)
	userRepoMock.On("FindUserByID", ctx, mock.Anything, mock.Anything).Once().Return(userEntity, nil)
	jwtServiceMock.On("GenerateTokens", ctx, userEntity).Once().Return(accessToken, refreshToken, nil)

	req := v0.RefreshTokensInput{
		RefreshToken: "refreshToken",
	}
	resp, err := svc.RefreshTokens(ctx, req)

	t.Require().NoError(err)
	t.Require().Equal(accessToken, resp.AccessToken)
	t.Require().Equal(refreshToken, resp.RefreshToken)
}
