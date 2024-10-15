package auth_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/cache"
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/persistence/repo"
	mock2 "mandarine/internal/api/persistence/repo/mock"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	mock4 "mandarine/pkg/smtp/mock"
	"mandarine/pkg/storage/cache/manager"
	mock3 "mandarine/pkg/storage/cache/manager/mock"
	mock5 "mandarine/pkg/template/mock"
	"strings"
	"testing"
	"time"
)

func Test_RegisterService_Register(t *testing.T) {
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
		Server: config.ServerConfig{
			ExternalOrigin: "https://example.com",
		},
	}
	service := auth.NewRegisterService(userRepo, cacheManager, smtpSender, templateEngine, cfg)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := dto.RegisterInput{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password",
		}

		userRepo.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
		templateEngine.On("Render", mock.Anything, mock.Anything).Once().Return("email content", nil)
		smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, req.Email, []string(nil)).Once().Return(nil)

		err := service.Register(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		req := dto.RegisterInput{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password",
		}

		userRepo.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(true, nil)

		err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrDuplicateUser, err)
	})

	t.Run("ErrorHashingPassword", func(t *testing.T) {
		req := dto.RegisterInput{
			Email:    "test@example.com",
			Username: "testuser",
			Password: strings.Repeat("1", 1000),
		}

		userRepo.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)

		err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, bcrypt.ErrPasswordTooLong, err)
	})

	t.Run("ErrorSavingCache", func(t *testing.T) {
		req := dto.RegisterInput{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password",
		}

		userRepo.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(errors.New("cache error"))

		err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "cache error", err.Error())
	})

	t.Run("ErrorSendingEmail", func(t *testing.T) {
		req := dto.RegisterInput{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password",
		}

		userRepo.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(false, nil)
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
		templateEngine.On("Render", mock.Anything, mock.Anything).Once().Return("email content", nil)
		smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, req.Email, []string(nil)).Once().Return(errors.New("smtp error"))

		err := service.Register(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "smtp error", err.Error())
	})
}

func Test_RegisterService_RegisterConfirm(t *testing.T) {
	userRepo := new(mock2.UserRepositoryMock)
	cacheManager := new(mock3.CacheManagerMock)
	service := auth.NewRegisterService(userRepo, cacheManager, nil, nil, nil)

	ctx := context.Background()
	req := dto.RegisterConfirmInput{
		OTP:   "123456",
		Email: "test@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("register", "test@example.com"), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)
		userEntity := &model.UserEntity{Email: "test@example.com", Username: "testuser"}
		userRepo.On("ExistsUserByUsernameOrEmail", ctx, "testuser", "test@example.com").Once().Return(false, nil)
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Once().Return(userEntity, nil)
		cacheManager.On("Delete", mock.Anything, []string{cache.CreateCacheKey("register", "test@example.com")}).Once().Return(nil)

		err := service.RegisterConfirm(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("OTPNotFound", func(t *testing.T) {
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Once().Return(manager.ErrCacheEntryNotFound)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("ErrorGettingCache", func(t *testing.T) {
		expectedErr := errors.New("cache error")
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Once().Return(expectedErr)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("InvalidOTP", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "654321",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "another@example.com",
				Username: "anotheruser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("ExistsUser", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)
		userRepo.On("ExistsUserByUsernameOrEmail", ctx, "testuser", "test@example.com").Once().Return(true, nil)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrDuplicateUser, err)
	})

	t.Run("ErrorExistsUser", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheError := errors.New("cache error")
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)
		userRepo.On("ExistsUserByUsernameOrEmail", ctx, "testuser", "test@example.com").Once().Return(false, cacheError)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, cacheError, err)
	})

	t.Run("DuplicateUser", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)
		userRepo.On("ExistsUserByUsernameOrEmail", ctx, "testuser", "test@example.com").Once().Return(false, nil)
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Once().Return(nil, repo.ErrDuplicateUser)

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrDuplicateUser, err)
	})

	t.Run("ErrorSavingUser", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)
		userRepo.On("ExistsUserByUsernameOrEmail", ctx, "testuser", "test@example.com").Once().Return(false, nil)
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Once().Return(nil, errors.New("db error"))

		err := service.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("ErrDeleteCache", func(t *testing.T) {
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Email:    "test@example.com",
				Username: "testuser",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(10 * time.Minute),
		}
		cacheError := errors.New("cache error")
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RegisterCache) = cacheEntry
		}).Once().Return(nil)
		userEntity := &model.UserEntity{Email: "test@example.com", Username: "testuser"}
		userRepo.On("ExistsUserByUsernameOrEmail", ctx, "testuser", "test@example.com").Once().Return(false, nil)
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Once().Return(userEntity, nil)
		cacheManager.On("Delete", mock.Anything, []string{cache.CreateCacheKey("register", "test@example.com")}).Once().Return(cacheError)

		err := service.RegisterConfirm(ctx, req)

		assert.NoError(t, err)
	})
}
