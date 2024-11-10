package account_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/config"
	accountDto "github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/account"
	"github.com/mandarine-io/Backend/internal/helper/random"
	"github.com/mandarine-io/Backend/internal/helper/security"
	model2 "github.com/mandarine-io/Backend/internal/persistence/model"
	mock2 "github.com/mandarine-io/Backend/internal/persistence/repo/mock"
	mock4 "github.com/mandarine-io/Backend/pkg/smtp/mock"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	mock3 "github.com/mandarine-io/Backend/pkg/storage/cache/mock"
	mock5 "github.com/mandarine-io/Backend/pkg/template/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
	"time"
)

var (
	userRepo       = new(mock2.UserRepositoryMock)
	cacheManager   = new(mock3.ManagerMock)
	smtpSender     = new(mock4.SenderMock)
	templateEngine = new(mock5.TemplateEngineMock)
	cfg            = &config.Config{}
	svc            = account.NewService(userRepo, cacheManager, smtpSender, templateEngine, cfg)
	ctx            = context.Background()
)

func Test_AccountService_GetAccount(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			email := "test@example.com"
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           email,
				IsEnabled:       true,
				IsEmailVerified: true,
				IsPasswordTemp:  false,
				DeletedAt:       nil,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			resp, err := svc.GetAccount(ctx, userID)

			assert.NoError(t, err)
			assert.Equal(t, userEntity.Email, resp.Email)
			assert.Equal(t, userEntity.IsEnabled, resp.IsEnabled)
			assert.Equal(t, userEntity.IsEmailVerified, resp.IsEmailVerified)
			assert.Equal(t, userEntity.IsPasswordTemp, resp.IsPasswordTemp)
			assert.Equal(t, userEntity.DeletedAt != nil, resp.IsDeleted)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			resp, err := svc.GetAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserNotFound, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, expectedErr)

			resp, err := svc.GetAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)
}

func Test_AccountService_UpdateUsername(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Username:        "old-username",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateUsernameInput{
				Username: "username",
			}

			userRepo.On("ExistsUserByUsername", ctx, "username").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

			resp, err := svc.UpdateUsername(ctx, userID, req)

			assert.NoError(t, err)
			assert.Equal(t, "username", resp.Username)
		},
	)

	t.Run(
		"ErrUserNotFound", func(t *testing.T) {
			req := accountDto.UpdateUsernameInput{
				Username: "username",
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			resp, err := svc.UpdateUsername(ctx, userID, req)

			assert.Equal(t, service.ErrUserNotFound, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrFindUserById", func(t *testing.T) {
			req := accountDto.UpdateUsernameInput{
				Username: "username",
			}
			err := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, err)

			resp, err1 := svc.UpdateUsername(ctx, userID, req)

			assert.Equal(t, err, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run("UsernameNotChanged", func(t *testing.T) {
		req := accountDto.UpdateUsernameInput{
			Username: "old-username",
		}
		userEntity := &model2.UserEntity{
			ID:              userID,
			Username:        "old-username",
			IsEmailVerified: true,
		}

		userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

		resp, err := svc.UpdateUsername(ctx, userID, req)

		assert.NoError(t, err)
		assert.Equal(t, "old-username", resp.Username)
	})

	t.Run(
		"ErrDuplicateUsername", func(t *testing.T) {
			req := accountDto.UpdateUsernameInput{
				Username: "username",
			}
			userEntity := &model2.UserEntity{
				ID:              userID,
				Username:        "old-username",
				IsEmailVerified: true,
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("ExistsUserByUsername", ctx, "username").Once().Return(true, nil)

			resp, err := svc.UpdateUsername(ctx, userID, req)

			assert.Equal(t, accountDto.AccountOutput{}, resp)
			assert.Equal(t, service.ErrDuplicateUsername, err)
		},
	)

	t.Run(
		"ErrExistsUserByUsername", func(t *testing.T) {
			req := accountDto.UpdateUsernameInput{
				Username: "username",
			}
			err := errors.New("database error")

			userEntity := &model2.UserEntity{
				ID:              userID,
				Username:        "old-username",
				IsEmailVerified: true,
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("ExistsUserByUsername", ctx, "username").Once().Return(true, err)

			resp, err1 := svc.UpdateUsername(ctx, userID, req)

			assert.Equal(t, accountDto.AccountOutput{}, resp)
			assert.Equal(t, err, err1)
		},
	)

	t.Run(
		"ErrUpdateUser", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Username:        "old-username",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateUsernameInput{
				Username: "username",
			}

			err := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("ExistsUserByUsername", ctx, "username").Once().Return(false, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, err)

			resp, err1 := svc.UpdateUsername(ctx, userID, req)

			assert.Equal(t, err, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)
}

func Test_AccountService_UpdateEmail(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			cacheManager.On(
				"SetWithExpiration",
				ctx,
				strings.Join([]string{"email-verify", req.Email}, "."),
				mock.Anything,
				time.Duration(cfg.Security.OTP.TTL)*time.Second).
				Once().Return(nil)
			templateEngine.On("Render", "email-verify", mock.Anything).Once().Return("content", nil)
			smtpSender.On("SendHtmlMessage", mock.Anything, mock.Anything, req.Email).Once().Return(nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

			resp, err := svc.UpdateEmail(ctx, userID, req, nil)

			assert.NoError(t, err)
			assert.Equal(t, "new@example.com", resp.Email)
			assert.False(t, resp.IsEmailVerified)
		},
	)

	t.Run(
		"ErrUserNotFound", func(t *testing.T) {
			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			resp, err := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, service.ErrUserNotFound, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrFindUserById", func(t *testing.T) {
			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}
			err := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, err)

			resp, err1 := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, err, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"EmailNotChanged", func(t *testing.T) {
			req := accountDto.UpdateEmailInput{
				Email: "test@example.com",
			}
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			resp, err := svc.UpdateEmail(ctx, userID, req, nil)

			assert.NoError(t, err)
			assert.Equal(t, req.Email, resp.Email)
		},
	)

	t.Run(
		"ErrDuplicateEmail", func(t *testing.T) {
			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(true, nil)

			resp, err := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, accountDto.AccountOutput{}, resp)
			assert.Equal(t, service.ErrDuplicateEmail, err)
		},
	)

	t.Run(
		"ErrExistsUserByEmail", func(t *testing.T) {
			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}
			err := errors.New("database error")
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(true, err)

			resp, err1 := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, accountDto.AccountOutput{}, resp)
			assert.Equal(t, err, err1)
		},
	)

	t.Run(
		"ErrGenerateOTP", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}
			cfg.Security.OTP.Length = -1

			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			resp, err := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, random.ErrInvalidOtpLength, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)

			cfg.Security.OTP.Length = 6
		},
	)

	t.Run(
		"ErrSetInCache", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			err := errors.New("cache error")
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			cacheManager.On(
				"SetWithExpiration",
				ctx,
				strings.Join([]string{"email-verify", req.Email}, "."),
				mock.Anything,
				time.Duration(cfg.Security.OTP.TTL)*time.Second).
				Once().Return(err)

			resp, err1 := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, err, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrRenderContent", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			err := errors.New("template error")
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			cacheManager.On(
				"SetWithExpiration",
				ctx,
				strings.Join([]string{"email-verify", req.Email}, "."),
				mock.Anything,
				time.Duration(cfg.Security.OTP.TTL)*time.Second).
				Once().Return(nil)
			templateEngine.On("Render", "email-verify", mock.Anything).Once().Return("", err)

			resp, err1 := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, err, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrSendHtmlMessage", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			err := errors.New("smtp error")
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			cacheManager.On(
				"SetWithExpiration",
				ctx,
				strings.Join([]string{"email-verify", req.Email}, "."),
				mock.Anything,
				time.Duration(cfg.Security.OTP.TTL)*time.Second).
				Once().Return(nil)
			templateEngine.On("Render", "email-verify", mock.Anything).Once().Return("content", nil)
			smtpSender.On("SendHtmlMessage", mock.Anything, "content", req.Email).Once().Return(err)

			resp, err1 := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, service.ErrSendEmail, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrUpdateUser", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.UpdateEmailInput{
				Email: "new@example.com",
			}

			err := errors.New("database error")
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			cacheManager.On(
				"SetWithExpiration",
				ctx,
				strings.Join([]string{"email-verify", req.Email}, "."),
				mock.Anything,
				time.Duration(cfg.Security.OTP.TTL)*time.Second).
				Once().Return(nil)
			templateEngine.On("Render", "email-verify", mock.Anything).Once().Return("content", nil)
			smtpSender.On("SendHtmlMessage", mock.Anything, "content", req.Email).Once().Return(nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, err)

			resp, err1 := svc.UpdateEmail(ctx, userID, req, nil)

			assert.Equal(t, err, err1)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)
}

func Test_AccountService_VerifyEmail(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)
			cacheManager.On("Delete", ctx, strings.Join([]string{"email-verify", req.Email}, ".")).
				Once().Return(nil)

			err := svc.VerifyEmail(ctx, userID, req)

			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrCacheEntryNotFound", func(t *testing.T) {
			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Once().Return(cache.ErrCacheEntryNotFound)

			err := svc.VerifyEmail(ctx, userID, req)

			assert.Equal(t, service.ErrInvalidOrExpiredOtp, err)
		},
	)

	t.Run(
		"ErrGetCache", func(t *testing.T) {
			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "654321",
			}

			err := errors.New("cache error")
			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Once().Return(err)

			err1 := svc.VerifyEmail(ctx, userID, req)

			assert.Error(t, err, err1)
		},
	)

	t.Run(
		"ErrMismatchOtp", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "654321",
			}

			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)

			err := svc.VerifyEmail(ctx, userID, req)

			assert.Error(t, service.ErrInvalidOrExpiredOtp, err)
		},
	)

	t.Run(
		"ErrUserNotFound", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			err := svc.VerifyEmail(ctx, userID, req)

			assert.Error(t, service.ErrUserNotFound, err)
		},
	)

	t.Run(
		"ErrFindUserById", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			err := errors.New("database error")
			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, err)

			err1 := svc.VerifyEmail(ctx, userID, req)

			assert.Equal(t, err, err1)
		},
	)

	t.Run(
		"ErrCheckEmail", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			userEntity := &model2.UserEntity{
				ID:    userID,
				Email: "another@example.com",
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)
			userRepo.On("ExistsUserByEmail", ctx, "new@example.com").Once().Return(false, nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			err := svc.VerifyEmail(ctx, userID, req)

			assert.Equal(t, service.ErrInvalidOrExpiredOtp, err)
		},
	)

	t.Run(
		"ErrUpdateUser", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			err := errors.New("database error")
			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(nil, err)

			err1 := svc.VerifyEmail(ctx, userID, req)

			assert.Equal(t, err, err1)
		},
	)

	t.Run(
		"ErrDeleteInCache", func(t *testing.T) {
			cacheEntry := accountDto.EmailVerifyCache{
				Email: "test@example.com",
				OTP:   "123456",
			}

			userEntity := &model2.UserEntity{
				ID:              userID,
				Email:           "test@example.com",
				IsEmailVerified: true,
			}

			req := accountDto.VerifyEmailInput{
				Email: "test@example.com",
				OTP:   "123456",
			}

			err := errors.New("cache error")
			cacheManager.On("Get", ctx, strings.Join([]string{"email-verify", req.Email}, "."), mock.Anything).
				Run(
					func(args mock.Arguments) {
						*args.Get(2).(*accountDto.EmailVerifyCache) = cacheEntry
					},
				).
				Once().Return(nil)
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)
			cacheManager.On("Delete", ctx, strings.Join([]string{"email-verify", req.Email}, ".")).
				Once().Return(err)

			err1 := svc.VerifyEmail(ctx, userID, req)

			assert.NoError(t, err1)
		},
	)
}

func Test_AccountService_SetPassword(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: true,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

			req := accountDto.SetPasswordInput{
				Password: "newpassword",
			}

			err := svc.SetPassword(ctx, userID, req)

			assert.NoError(t, err)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			req := accountDto.SetPasswordInput{
				Password: "newpassword",
			}

			err := svc.SetPassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserNotFound, err)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, expectedErr)

			req := accountDto.SetPasswordInput{
				Password: "newpassword",
			}

			err := svc.SetPassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)

	t.Run(
		"PasswordAlreadySet", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: false,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			req := accountDto.SetPasswordInput{
				Password: "newpassword",
			}

			err := svc.SetPassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, service.ErrPasswordIsSet, err)
		},
	)

	t.Run(
		"ErrorHashPassword", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: true,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			req := accountDto.SetPasswordInput{
				Password: strings.Repeat("1", 1000),
			}

			err := svc.SetPassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, bcrypt.ErrPasswordTooLong, err)
		},
	)

	t.Run(
		"ErrorUpdateUser", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: true,
			}
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

			req := accountDto.SetPasswordInput{
				Password: "newpassword",
			}

			err := svc.SetPassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)
}

func Test_AccountService_UpdatePassword(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			hashPassword, _ := security.HashPassword("oldpassword")
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: false,
				Password:       hashPassword,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

			req := accountDto.UpdatePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			}

			err := svc.UpdatePassword(ctx, userID, req)

			assert.NoError(t, err)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			req := accountDto.UpdatePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			}

			err := svc.UpdatePassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserNotFound, err)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, expectedErr)

			req := accountDto.UpdatePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			}

			err := svc.UpdatePassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)

	t.Run(
		"IncorrectOldPassword", func(t *testing.T) {
			hashPassword, _ := security.HashPassword("oldpassword")
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: false,
				Password:       hashPassword,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			req := accountDto.UpdatePasswordInput{
				OldPassword: "wrongoldpassword",
				NewPassword: "newpassword",
			}

			err := svc.UpdatePassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, service.ErrIncorrectOldPassword, err)
		},
	)

	t.Run(
		"ErrorHashPassword", func(t *testing.T) {
			hashPassword, _ := security.HashPassword("oldpassword")
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: false,
				Password:       hashPassword,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			req := accountDto.UpdatePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: strings.Repeat("1", 1000),
			}

			err := svc.UpdatePassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, bcrypt.ErrPasswordTooLong, err)
		},
	)

	t.Run(
		"ErrorUpdateUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			hashPassword, _ := security.HashPassword("oldpassword")
			userEntity := &model2.UserEntity{
				ID:             userID,
				IsPasswordTemp: false,
				Password:       hashPassword,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

			req := accountDto.UpdatePasswordInput{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			}

			err := svc.UpdatePassword(ctx, userID, req)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)
}

func Test_AccountService_RestoreAccount(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			deletedAt := time.Now()
			userEntity := &model2.UserEntity{
				ID:        userID,
				Email:     "test@example.com",
				IsEnabled: true,
				DeletedAt: &deletedAt,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

			resp, err := svc.RestoreAccount(ctx, userID)

			assert.NoError(t, err)
			assert.Equal(t, userEntity.Email, resp.Email)
			assert.Equal(t, userEntity.IsEnabled, resp.IsEnabled)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			resp, err := svc.RestoreAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserNotFound, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, expectedErr)

			resp, err := svc.RestoreAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"UserNotDeleted", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID:        userID,
				DeletedAt: nil,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			resp, err := svc.RestoreAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserNotDeleted, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)

	t.Run(
		"ErrorUpdateUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			deletedAt := time.Now()
			userEntity := &model2.UserEntity{
				ID:        userID,
				DeletedAt: &deletedAt,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

			resp, err := svc.RestoreAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, accountDto.AccountOutput{}, resp)
		},
	)
}

func Test_AccountService_DeleteAccount(t *testing.T) {
	userID := uuid.New()

	t.Run(
		"Success", func(t *testing.T) {
			userEntity := &model2.UserEntity{
				ID: userID,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(userEntity, nil)

			err := svc.DeleteAccount(ctx, userID)

			assert.NoError(t, err)
		},
	)

	t.Run(
		"UserNotFound", func(t *testing.T) {
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, nil)

			err := svc.DeleteAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserNotFound, err)
		},
	)

	t.Run(
		"ErrorFindingUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(nil, expectedErr)

			err := svc.DeleteAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)

	t.Run(
		"UserAlreadyDeleted", func(t *testing.T) {
			deletedAt := time.Now()
			userEntity := &model2.UserEntity{
				ID:        userID,
				DeletedAt: &deletedAt,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)

			err := svc.DeleteAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, service.ErrUserAlreadyDeleted, err)
		},
	)

	t.Run(
		"ErrorUpdateUser", func(t *testing.T) {
			expectedErr := errors.New("database error")
			userEntity := &model2.UserEntity{
				ID:        userID,
				DeletedAt: nil,
			}
			userRepo.On("FindUserById", ctx, userID, false).Once().Return(userEntity, nil)
			userRepo.On("UpdateUser", ctx, userEntity).Once().Return(nil, expectedErr)

			err := svc.DeleteAccount(ctx, userID)

			assert.Error(t, err)
			assert.Equal(t, expectedErr, err)
		},
	)
}
