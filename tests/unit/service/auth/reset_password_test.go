package auth_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/cache"
	"mandarine/internal/api/persistence/model"
	mock2 "mandarine/internal/api/persistence/repo/mock"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	mock4 "mandarine/pkg/smtp/mock"
	"mandarine/pkg/storage/cache/manager"
	mock3 "mandarine/pkg/storage/cache/manager/mock"
	mock5 "mandarine/pkg/template/mock"
	"testing"
)

func Test_ResetPasswordService_RecoveryPassword(t *testing.T) {
	userRepo := new(mock2.UserRepositoryMock)
	cacheManager := new(mock3.CacheManagerMock)
	smtpSender := new(mock4.SenderMock)
	templateEngine := new(mock5.TemplateEngineMock)
	cfg := &config.Config{
		Security: config.SecurityConfig{
			OTP: config.OTPConfig{
				Length: 6,
				TTL:    300,
			},
		},
		Server: config.ServerConfig{ExternalOrigin: "https://example.com"},
	}
	service := auth.NewResetPasswordService(userRepo, cacheManager, smtpSender, templateEngine, cfg)

	input := dto.RecoveryPasswordInput{Email: "test@example.com"}

	t.Run("Success", func(t *testing.T) {
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, input.Email, []string(nil)).Return(nil).Once()
		templateEngine.On("Render", "recovery-password", mock.Anything).Return("email content", nil).Once()
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		err := service.RecoveryPassword(context.Background(), input)

		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(nil, nil).Once()

		err := service.RecoveryPassword(context.Background(), input)

		assert.Equal(t, auth.ErrUserNotFound, err)
	})

	t.Run("ErrorSettingCache", func(t *testing.T) {
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("cache error")).Once()

		err := service.RecoveryPassword(context.Background(), input)

		assert.Error(t, err)
	})

	t.Run("ErrorSendingEmail", func(t *testing.T) {
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		templateEngine.On("Render", "recovery-password", mock.Anything).Return("email content", nil).Once()
		smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, input.Email, []string(nil)).Return(errors.New("smtp error")).Once()

		err := service.RecoveryPassword(context.Background(), input)

		assert.Error(t, err)
	})
}

func Test_ResetPasswordService_VerifyRecoveryCode(t *testing.T) {
	userRepo := new(mock2.UserRepositoryMock)
	cacheManager := new(mock3.CacheManagerMock)
	cfg := &config.Config{}
	service := auth.NewResetPasswordService(userRepo, cacheManager, nil, nil, cfg)

	t.Run("Success", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()

		err := service.VerifyRecoveryCode(context.Background(), input)

		assert.NoError(t, err)
	})

	t.Run("InvalidOrExpiredOtp", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "wrong"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(manager.ErrCacheEntryNotFound).Once()

		err := service.VerifyRecoveryCode(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("CacheEntryNotFound", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(manager.ErrCacheEntryNotFound).Once()

		err := service.VerifyRecoveryCode(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("ErrorGettingCache", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(errors.New("cache error")).Once()

		err := service.VerifyRecoveryCode(context.Background(), input)

		assert.Error(t, err)
	})
}

func Test_ResetPasswordService_ResetPassword(t *testing.T) {
	userRepo := new(mock2.UserRepositoryMock)
	cacheManager := new(mock3.CacheManagerMock)
	smtpSender := new(mock4.SenderMock)
	templateEngine := new(mock5.TemplateEngineMock)
	cfg := &config.Config{}
	service := auth.NewResetPasswordService(userRepo, cacheManager, smtpSender, templateEngine, cfg)

	t.Run("Success", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		userRepo.On("UpdateUser", mock.Anything, userEntity).Return(userEntity, nil).Once()

		err := service.ResetPassword(context.Background(), input)

		assert.NoError(t, err)
	})

	t.Run("InvalidOrExpiredOtp", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "wrong", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()

		err := service.ResetPassword(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("CacheEntryNotFound", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(manager.ErrCacheEntryNotFound).Once()

		err := service.ResetPassword(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(nil, nil).Once()

		err := service.ResetPassword(context.Background(), input)

		assert.Equal(t, auth.ErrUserNotFound, err)
	})

	t.Run("ErrorUpdatingUser", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		userRepo.On("UpdateUser", mock.Anything, userEntity).Return(userEntity, errors.New("update error")).Once()

		err := service.ResetPassword(context.Background(), input)

		assert.Error(t, err)
	})
}
