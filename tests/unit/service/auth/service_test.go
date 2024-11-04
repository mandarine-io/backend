package auth_test

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/internal/api/helper/cache"
	"github.com/mandarine-io/Backend/internal/api/helper/security"
	"github.com/mandarine-io/Backend/internal/api/persistence/model"
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	mock2 "github.com/mandarine-io/Backend/internal/api/persistence/repo/mock"
	"github.com/mandarine-io/Backend/internal/api/service/auth"
	"github.com/mandarine-io/Backend/internal/api/service/auth/dto"
	"github.com/mandarine-io/Backend/pkg/oauth"
	mock3 "github.com/mandarine-io/Backend/pkg/oauth/mock"
	mock5 "github.com/mandarine-io/Backend/pkg/smtp/mock"
	cache2 "github.com/mandarine-io/Backend/pkg/storage/cache"
	mock4 "github.com/mandarine-io/Backend/pkg/storage/cache/mock"
	mock6 "github.com/mandarine-io/Backend/pkg/template/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	ctx = context.Background()

	userRepo        *mock2.UserRepositoryMock
	bannedTokenRepo *mock2.BannedTokenRepositoryMock
	oauthProviders  map[string]oauth.Provider
	cacheManager    *mock4.ManagerMock
	smtpSender      *mock5.SenderMock
	templateEngine  *mock6.TemplateEngineMock
	cfg             *config.Config
	svc             *auth.Service
)

func TestMain(m *testing.M) {
	// Setup logger
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.Level(10000),
			},
		),
	)
	slog.SetDefault(logger)

	// Setup mocks
	userRepo = new(mock2.UserRepositoryMock)
	bannedTokenRepo = new(mock2.BannedTokenRepositoryMock)
	oauthProviders = make(map[string]oauth.Provider)
	oauthProviders["mock"] = new(mock3.OAuthProviderMock)
	cacheManager = new(mock4.ManagerMock)
	smtpSender = new(mock5.SenderMock)
	templateEngine = new(mock6.TemplateEngineMock)
	cfg = &config.Config{
		Server: config.ServerConfig{
			ExternalOrigin: "https://example.com",
		},
		Security: config.SecurityConfig{
			JWT: config.JWTConfig{
				Secret:          "LMRdskYdRNdXA0m1YK3stPFWAciSiwkvQVOZNebYvFI=",
				AccessTokenTTL:  3600,
				RefreshTokenTTL: 86400,
			},
			OTP: config.OTPConfig{
				Length: 6,
				TTL:    300,
			},
		},
	}
	svc = auth.NewService(userRepo, bannedTokenRepo, oauthProviders, cacheManager, smtpSender, templateEngine, cfg)

	os.Exit(m.Run())
}

func Test_AuthService_Login(t *testing.T) {
	req := dto.LoginInput{Login: "test@example.com", Password: "password123"}

	t.Run(
		"ErrUserNotFound", func(t *testing.T) {
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(nil, nil)

			resp, err := svc.Login(ctx, req)

			assert.Equal(t, err, auth.ErrUserNotFound)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserByUsernameOrEmail", ctx, req.Login, true).Once().Return(nil, expectedErr)

			resp, err := svc.Login(ctx, req)

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

			resp, err := svc.Login(ctx, req)

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

			resp, err := svc.Login(ctx, req)

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

			resp, err := svc.Login(ctx, req)

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.AccessToken)
			assert.NotEmpty(t, resp.RefreshToken)
		},
	)
}

func Test_AuthService_RefreshTokens(t *testing.T) {
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
			resp, err := svc.RefreshTokens(ctx, "invalid_refresh_token")

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
			resp, err := svc.RefreshTokens(ctx, refreshTokenSigned)

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
			resp, err := svc.RefreshTokens(ctx, refreshTokenSigned)

			assert.Equal(t, auth.ErrInvalidJwtToken, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, mock.Anything, true).Once().Return(nil, nil)

			resp, err := svc.RefreshTokens(ctx, refreshToken)

			assert.Equal(t, auth.ErrUserNotFound, err)
			assert.Equal(t, dto.JwtTokensOutput{}, resp)
		},
	)

	t.Run(
		"ErrFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, mock.Anything, true).Once().Return(nil, expectedErr)

			resp, err := svc.RefreshTokens(ctx, refreshToken)

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

			resp, err := svc.RefreshTokens(ctx, refreshToken)

			assert.NoError(t, err)
			assert.NotEmpty(t, resp.AccessToken)
			assert.NotEmpty(t, resp.RefreshToken)
		},
	)
}

func Test_AuthService_Logout(t *testing.T) {
	jti := uuid.New().String()

	t.Run(
		"Success", func(t *testing.T) {
			bannedTokenRepo.On("CreateOrUpdateBannedToken", ctx, mock.Anything).Once().Return(nil, nil)

			err := svc.Logout(ctx, jti)

			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrCreateOrUpdateBannedToken", func(t *testing.T) {
			expectedErr := errors.New("database error")
			bannedTokenRepo.On("CreateOrUpdateBannedToken", ctx, mock.Anything).Once().Return(nil, expectedErr)

			err := svc.Logout(ctx, jti)

			assert.Equal(t, expectedErr, err)
		},
	)
}

func Test_AuthService_Register(t *testing.T) {
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

		err := svc.Register(ctx, req, nil)

		assert.NoError(t, err)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		req := dto.RegisterInput{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password",
		}

		userRepo.On("ExistsUserByUsernameOrEmail", mock.Anything, req.Username, req.Email).Once().Return(true, nil)

		err := svc.Register(ctx, req, nil)

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

		err := svc.Register(ctx, req, nil)

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

		err := svc.Register(ctx, req, nil)

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

		err := svc.Register(ctx, req, nil)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrSendEmail, err)
	})
}

func Test_AuthService_RegisterConfirm(t *testing.T) {
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

		err := svc.RegisterConfirm(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("OTPNotFound", func(t *testing.T) {
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Once().Return(cache2.ErrCacheEntryNotFound)

		err := svc.RegisterConfirm(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("ErrorGettingCache", func(t *testing.T) {
		expectedErr := errors.New("cache error")
		cacheManager.On("Get", mock.Anything, mock.Anything, mock.Anything).Once().Return(expectedErr)

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

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

		err := svc.RegisterConfirm(ctx, req)

		assert.NoError(t, err)
	})
}

func Test_AuthService_RecoveryPassword(t *testing.T) {
	input := dto.RecoveryPasswordInput{Email: "test@example.com"}

	t.Run("Success", func(t *testing.T) {
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, input.Email, []string(nil)).Return(nil).Once()
		templateEngine.On("Render", "recovery-password", mock.Anything).Return("email content", nil).Once()
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		err := svc.RecoveryPassword(context.Background(), input, nil)

		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(nil, nil).Once()

		err := svc.RecoveryPassword(context.Background(), input, nil)

		assert.Equal(t, auth.ErrUserNotFound, err)
	})

	t.Run("ErrorSettingCache", func(t *testing.T) {
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("cache error")).Once()

		err := svc.RecoveryPassword(context.Background(), input, nil)

		assert.Error(t, err)
	})

	t.Run("ErrorSendingEmail", func(t *testing.T) {
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		cacheManager.On("SetWithExpiration", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		templateEngine.On("Render", "recovery-password", mock.Anything).Return("email content", nil).Once()
		smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, input.Email, []string(nil)).Return(errors.New("smtp error")).Once()

		err := svc.RecoveryPassword(context.Background(), input, nil)

		assert.Error(t, err)
	})
}

func Test_AuthService_VerifyRecoveryCode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()

		err := svc.VerifyRecoveryCode(context.Background(), input)

		assert.NoError(t, err)
	})

	t.Run("InvalidOrExpiredOtp", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "wrong"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(cache2.ErrCacheEntryNotFound).Once()

		err := svc.VerifyRecoveryCode(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("CacheEntryNotFound", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(cache2.ErrCacheEntryNotFound).Once()

		err := svc.VerifyRecoveryCode(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("ErrorGettingCache", func(t *testing.T) {
		input := dto.VerifyRecoveryCodeInput{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(errors.New("cache error")).Once()

		err := svc.VerifyRecoveryCode(context.Background(), input)

		assert.Error(t, err)
	})
}

func Test_AuthService_ResetPassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(userEntity, nil).Once()
		userRepo.On("UpdateUser", mock.Anything, userEntity).Return(userEntity, nil).Once()

		err := svc.ResetPassword(context.Background(), input)

		assert.NoError(t, err)
	})

	t.Run("InvalidOrExpiredOtp", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "wrong", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()

		err := svc.ResetPassword(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("CacheEntryNotFound", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Return(cache2.ErrCacheEntryNotFound).Once()

		err := svc.ResetPassword(context.Background(), input)

		assert.Equal(t, auth.ErrInvalidOrExpiredOtp, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		input := dto.ResetPasswordInput{Email: "test@example.com", OTP: "123456", Password: "newpassword"}
		cacheEntry := dto.RecoveryPasswordCache{Email: "test@example.com", OTP: "123456"}

		cacheManager.On("Get", mock.Anything, cache.CreateCacheKey("recovery_password", input.Email), mock.Anything).Run(func(args mock.Arguments) {
			*args.Get(2).(*dto.RecoveryPasswordCache) = cacheEntry
		}).Return(nil).Once()
		userRepo.On("FindUserByEmail", mock.Anything, input.Email, false).Return(nil, nil).Once()

		err := svc.ResetPassword(context.Background(), input)

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

		err := svc.ResetPassword(context.Background(), input)

		assert.Error(t, err)
	})
}

func Test_AuthService_GetConsentPageUrl(t *testing.T) {
	redirectUrl := "https://example.com/callback"
	oauthProvider := oauthProviders["mock"].(*mock3.OAuthProviderMock)

	t.Run("NotSupportedProvider", func(t *testing.T) {
		_, err := svc.GetConsentPageUrl(context.Background(), "unsupported", redirectUrl)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrInvalidProvider, err)
	})

	t.Run("Success", func(t *testing.T) {
		oauthProvider.On("GetConsentPageUrl", redirectUrl).Return("consentUrl", "oauthState").Once()

		result, err := svc.GetConsentPageUrl(context.Background(), "mock", redirectUrl)

		assert.NoError(t, err)
		assert.Equal(t, "consentUrl", result.ConsentPageUrl)
		assert.Equal(t, "oauthState", result.OauthState)
	})
}

func Test_AuthService_FetchUserInfo(t *testing.T) {
	oauthProvider := oauthProviders["mock"].(*mock3.OAuthProviderMock)

	t.Run("NotSupportedProvider", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		_, err := svc.FetchUserInfo(context.Background(), "unsupported", input)

		assert.Error(t, err)
		assert.Equal(t, auth.ErrInvalidProvider, err)
	})

	t.Run("Success", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		expectedUserInfo := oauth.UserInfo{Email: "test@example.com"}

		oauthProvider.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(&oauth2.Token{}, nil).Once()
		oauthProvider.On("GetUserInfo", mock.Anything, mock.Anything).Return(expectedUserInfo, nil).Once()

		userInfo, err := svc.FetchUserInfo(context.Background(), "mock", input)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserInfo, userInfo)
	})

	t.Run("ErrorExchangingCodeToToken", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		expectedError := errors.New("exchange error")

		oauthProvider.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(nil, expectedError).Once()

		_, err := svc.FetchUserInfo(context.Background(), "mock", input)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ErrorGettingUserInfo", func(t *testing.T) {
		input := dto.FetchUserInfoInput{Code: "someCode"}
		token := &oauth2.Token{}
		expectedError := errors.New("user info error")

		oauthProvider.On("ExchangeCodeToToken", mock.Anything, input.Code, mock.Anything).Return(token, nil).Once()
		oauthProvider.On("GetUserInfo", mock.Anything, token).Return(oauth.UserInfo{}, expectedError).Once()

		_, err := svc.FetchUserInfo(context.Background(), "mock", input)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func Test_AuthService_RegisterOrLogin(t *testing.T) {
	t.Run("Success_NewUser_UniqueUsername", func(t *testing.T) {
		userInfo := oauth.UserInfo{Username: "test", Email: "test@example.com"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(userEntity, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, userInfo.Username).Return(false, nil).Once()

		result, err := svc.RegisterOrLogin(context.Background(), userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
	})

	t.Run("Success_NewUser_NotUniqueUsername", func(t *testing.T) {
		userInfo := oauth.UserInfo{Username: "test", Email: "test@example.com"}
		userEntity := &model.UserEntity{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(userEntity, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, userInfo.Username).Return(true, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, nil).Once()

		result, err := svc.RegisterOrLogin(context.Background(), userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
	})

	t.Run("Success_ExistingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}
		userEntity := &model.UserEntity{Email: "test@example.com", IsEnabled: true}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(userEntity, nil).Once()

		result, err := svc.RegisterOrLogin(context.Background(), userInfo)

		assert.NoError(t, err)
		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
	})

	t.Run("ErrorFindingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}
		expectedError := errors.New("repo error")

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, expectedError).Once()

		_, err := svc.RegisterOrLogin(context.Background(), userInfo)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ErrorExistingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}
		expectedError := errors.New("repo error")

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, expectedError).Once()

		_, err := svc.RegisterOrLogin(context.Background(), userInfo)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ErrorCreatingUser", func(t *testing.T) {
		userInfo := oauth.UserInfo{Email: "test@example.com"}

		userRepo.On("FindUserByEmail", mock.Anything, userInfo.Email, true).Return(nil, nil).Once()
		userRepo.On("ExistsUserByUsername", mock.Anything, mock.Anything).Return(false, nil).Once()
		userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("create error")).Once()

		_, err := svc.RegisterOrLogin(context.Background(), userInfo)

		assert.Error(t, err)
		assert.Equal(t, "create error", err.Error())
	})
}
