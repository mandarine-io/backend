package auth_test

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/model"
	mock2 "mandarine/internal/api/persistence/repo/mock"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.Level(10000),
			},
		),
	)
	slog.SetDefault(logger)
}

func Test_LoginService_Login(t *testing.T) {
	userRepo := new(mock2.UserRepositoryMock)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			JWT: config.JWTConfig{
				Secret:          "LMRdskYdRNdXA0m1YK3stPFWAciSiwkvQVOZNebYvFI=",
				AccessTokenTTL:  3600,
				RefreshTokenTTL: 86400,
			},
		},
	}

	loginService := auth.NewLoginService(userRepo, cfg)

	ctx := context.Background()
	req := dto.LoginInput{Login: "test@example.com", Password: "password123"}

	t.Run(
		"ErrUserNotFound", func(t *testing.T) {
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(nil, nil)

			resp, err := loginService.Login(ctx, req)

			assert.Equal(t, err, auth.ErrUserNotFound)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(nil, expectedErr)

			resp, err := loginService.Login(ctx, req)

			assert.Equal(t, err, expectedErr)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrBadCredentials", func(t *testing.T) {
			userEntity := &model.UserEntity{
				Email:    req.Login,
				Password: "hashedpassword",
			}
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(userEntity, nil)

			resp, err := loginService.Login(ctx, req)

			assert.Equal(t, err, auth.ErrBadCredentials)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrUserIsBlocked", func(t *testing.T) {
			hashPassword, _ := security.HashPassword("password123")
			userEntity := &model.UserEntity{
				Email:     req.Login,
				Password:  hashPassword,
				IsEnabled: false,
			}
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(userEntity, nil)

			resp, err := loginService.Login(ctx, req)

			assert.Equal(t, err, auth.ErrUserIsBlocked)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"Successful login", func(t *testing.T) {
			hashPassword, _ := security.HashPassword("password123")
			userEntity := &model.UserEntity{
				Email:     req.Login,
				Password:  hashPassword,
				IsEnabled: true,
			}
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(userEntity, nil)

			resp, err := loginService.Login(ctx, req)

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.AccessToken)
			assert.NotEmpty(t, resp.RefreshToken)
		},
	)
}

func Test_LoginService_RefreshTokens(t *testing.T) {
	userRepo := new(mock2.UserRepositoryMock)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			JWT: config.JWTConfig{
				Secret:          "LMRdskYdRNdXA0m1YK3stPFWAciSiwkvQVOZNebYvFI=",
				AccessTokenTTL:  3600,
				RefreshTokenTTL: 86400,
			},
		},
	}

	loginService := auth.NewLoginService(userRepo, cfg)

	ctx := context.Background()
	userEntity := &model.UserEntity{
		ID:    uuid.New(),
		Email: "example@mail.ru",
		Role: model.RoleEntity{
			Name: model.RoleAdmin,
		},
	}
	_, refreshToken, _ := security.GenerateTokens(cfg.Security.JWT, userEntity)

	t.Run(
		"InvalidJwtToken", func(t *testing.T) {
			resp, err := loginService.RefreshTokens(ctx, "invalid_refresh_token")

			assert.Equal(t, auth.ErrInvalidJwtToken, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrSubClaimsNotFound", func(t *testing.T) {
			refreshToken := jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				jwt.MapClaims{},
			)
			refreshTokenSigned, _ := refreshToken.SignedString([]byte(cfg.Security.JWT.Secret))
			resp, err := loginService.RefreshTokens(ctx, refreshTokenSigned)

			assert.Equal(t, auth.ErrInvalidJwtToken, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrInvalidUserUUID", func(t *testing.T) {
			refreshToken := jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				jwt.MapClaims{
					"sub": "wrong_uuid",
				},
			)
			refreshTokenSigned, _ := refreshToken.SignedString([]byte(cfg.Security.JWT.Secret))
			resp, err := loginService.RefreshTokens(ctx, refreshTokenSigned)

			assert.Equal(t, auth.ErrInvalidJwtToken, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, mock.Anything, true).Once().Return(nil, nil)

			resp, err := loginService.RefreshTokens(ctx, refreshToken)

			assert.Equal(t, auth.ErrUserNotFound, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, mock.Anything, true).Once().Return(nil, expectedErr)

			resp, err := loginService.RefreshTokens(ctx, refreshToken)

			assert.Equal(t, expectedErr, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"Success", func(t *testing.T) {
			userEntity := &model.UserEntity{
				ID:        uuid.New(),
				Email:     "test@example.com",
				IsEnabled: true,
			}
			userRepo.On("FindUserById", ctx, mock.Anything, true).Once().Return(userEntity, nil)

			resp, err := loginService.RefreshTokens(ctx, refreshToken)

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.AccessToken)
			assert.NotEmpty(t, resp.RefreshToken)
		},
	)
}
