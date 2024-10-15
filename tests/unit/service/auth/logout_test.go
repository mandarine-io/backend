package auth_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mandarine/internal/api/config"
	mock2 "mandarine/internal/api/persistence/repo/mock"
	"mandarine/internal/api/service/auth"
	"testing"
)

func Test_LogoutService_Login(t *testing.T) {
	bannedTokenRepo := new(mock2.BannedTokenRepositoryMock)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			JWT: config.JWTConfig{
				RefreshTokenTTL: 86400,
			},
		},
	}

	loginService := auth.NewLogoutService(bannedTokenRepo, cfg)

	ctx := context.Background()
	jti := uuid.New().String()

	t.Run(
		"Success", func(t *testing.T) {
			bannedTokenRepo.On("CreateOrUpdateBannedToken", ctx, mock.Anything).Once().Return(nil, nil)

			err := loginService.Logout(ctx, jti)

			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrCreateOrUpdateBannedToken", func(t *testing.T) {
			expectedErr := errors.New("database error")
			bannedTokenRepo.On("CreateOrUpdateBannedToken", ctx, mock.Anything).Once().Return(nil, expectedErr)

			err := loginService.Logout(ctx, jti)

			assert.Equal(t, expectedErr, err)
		},
	)
}
