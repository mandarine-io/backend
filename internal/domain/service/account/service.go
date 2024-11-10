package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/account/mapper"
	cachehelper "github.com/mandarine-io/Backend/internal/helper/cache"
	"github.com/mandarine-io/Backend/internal/helper/random"
	"github.com/mandarine-io/Backend/internal/helper/security"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/mandarine-io/Backend/pkg/template"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	emailVerifyCachePrefix = "email-verify"
	emailDefaultTitle      = "Verify email"
)

type svc struct {
	userRepo       repo.UserRepository
	cacheManager   cache.Manager
	smtpSender     smtp.Sender
	templateEngine template.Engine
	cfg            *config.Config
}

func NewService(
	userRepo repo.UserRepository,
	cacheManager cache.Manager,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	cfg *config.Config,
) service.AccountService {
	return &svc{
		userRepo:       userRepo,
		cacheManager:   cacheManager,
		smtpSender:     smtpSender,
		templateEngine: templateEngine,
		cfg:            cfg,
	}
}

//////////////////// Get account ////////////////////

func (s *svc) GetAccount(ctx context.Context, id uuid.UUID) (dto.AccountOutput, error) {
	log.Info().Msgf("get account: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get account")
		return dto.AccountOutput{}, err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to get account")
		return dto.AccountOutput{}, service.ErrUserNotFound
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Update username ////////////////////

func (s *svc) UpdateUsername(
	ctx context.Context, id uuid.UUID, input dto.UpdateUsernameInput,
) (dto.AccountOutput, error) {
	log.Info().Msgf("update username: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update username")
		return dto.AccountOutput{}, err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to update username")
		return dto.AccountOutput{}, service.ErrUserNotFound
	}

	// Check if username not changed
	if input.Username == userEntity.Username {
		return mapper.MapUserEntityToAccountResponse(userEntity), nil
	}

	// Check if username is already in use
	exists, err := s.userRepo.ExistsUserByUsername(ctx, input.Username)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update username")
		return dto.AccountOutput{}, err
	}
	if exists {
		log.Error().Stack().Err(service.ErrDuplicateUsername).Msg("failed to update username")
		return dto.AccountOutput{}, service.ErrDuplicateUsername
	}

	// Update username
	userEntity.Username = input.Username
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update username")
		return dto.AccountOutput{}, err
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Update email ////////////////////

func (s *svc) UpdateEmail(
	ctx context.Context, id uuid.UUID, input dto.UpdateEmailInput, localizer *i18n.Localizer,
) (dto.AccountOutput, error) {
	log.Info().Msgf("update email: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update email")
		return dto.AccountOutput{}, err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to update email")
		return dto.AccountOutput{}, service.ErrUserNotFound
	}

	// Check if email not changed
	if input.Email == userEntity.Email {
		return mapper.MapUserEntityToAccountResponse(userEntity), nil
	}

	// Check if email is already in use
	exists, err := s.userRepo.ExistsUserByEmail(ctx, input.Email)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update email")
		return dto.AccountOutput{}, err
	}
	if exists {
		log.Error().Stack().Err(service.ErrDuplicateEmail).Msg("failed to update email")
		return dto.AccountOutput{}, service.ErrDuplicateEmail
	}

	// Generate OTP
	otp, err := random.GenerateRandomNumber(s.cfg.Security.OTP.Length)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update email")
		return dto.AccountOutput{}, err
	}

	// Save in cache
	expiration := time.Duration(s.cfg.Security.OTP.TTL) * time.Second
	cacheEntry := dto.EmailVerifyCache{
		Email:     input.Email,
		OTP:       otp,
		ExpiredAt: time.Now().Add(expiration),
	}
	err = s.cacheManager.SetWithExpiration(
		ctx, cachehelper.CreateCacheKey(emailVerifyCachePrefix, input.Email), cacheEntry, expiration,
	)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update email")
		return dto.AccountOutput{}, err
	}

	// Localize email title
	emailTitle := emailDefaultTitle
	if localizer != nil {
		emailTitle = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email.email-verify.title"})
	}

	// Send mail
	args := dto.EmailVerifyTemplateArgs{
		Email: input.Email,
		TTL:   s.cfg.Security.OTP.TTL / 60,
		OTP:   otp,
	}
	content, err := s.templateEngine.Render("email-verify", args)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get account")
		return dto.AccountOutput{}, err
	}

	err = s.smtpSender.SendHtmlMessage(emailTitle, content, input.Email)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get account")
		return dto.AccountOutput{}, service.ErrSendEmail
	}

	// Update email
	userEntity.Email = input.Email
	userEntity.IsEmailVerified = false
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to get account")
		return dto.AccountOutput{}, err
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Verify email ////////////////////

func (s *svc) VerifyEmail(ctx context.Context, id uuid.UUID, req dto.VerifyEmailInput) error {
	log.Info().Msgf("verify email: %s", id.String())

	// Get entry from cache
	var cacheEntry dto.EmailVerifyCache
	err := s.cacheManager.Get(ctx, cachehelper.CreateCacheKey(emailVerifyCachePrefix, req.Email), &cacheEntry)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to verify email")

		if errors.Is(err, cache.ErrCacheEntryNotFound) {
			return service.ErrInvalidOrExpiredOtp
		}
		return err
	}

	// Check OTP
	if req.OTP != cacheEntry.OTP {
		log.Error().Stack().Err(service.ErrInvalidOrExpiredOtp).Msg("failed to verify email")
		return service.ErrInvalidOrExpiredOtp
	}

	// Get user entity by id
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to verify email")
		return err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to verify email")
		return service.ErrUserNotFound
	}

	// Check email
	if userEntity.Email != cacheEntry.Email {
		log.Error().Stack().Err(service.ErrInvalidOrExpiredOtp).Msg("failed to verify email")
		return service.ErrInvalidOrExpiredOtp
	}

	// Verify email
	userEntity.IsEmailVerified = true
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to verify email")
		return err
	}

	// Delete cache entry
	err = s.cacheManager.Delete(ctx, cachehelper.CreateCacheKey(emailVerifyCachePrefix, req.Email))
	if err != nil {
		log.Warn().Err(err).Msg("failed to delete cache entry")
	}

	return nil
}

//////////////////// Set password ////////////////////

func (s *svc) SetPassword(ctx context.Context, id uuid.UUID, input dto.SetPasswordInput) error {
	log.Info().Msgf("set password: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to set password")
		return err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to set password")
		return service.ErrUserNotFound
	}

	// Check if password is empty
	if !userEntity.IsPasswordTemp {
		log.Error().Stack().Err(service.ErrPasswordIsSet).Msg("failed to set password")
		return service.ErrPasswordIsSet
	}

	// Hash new password
	userEntity.Password, err = security.HashPassword(input.Password)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to set password")
		return err
	}

	// Update password
	userEntity.IsPasswordTemp = false
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to set password")
		return err
	}

	return nil
}

//////////////////// Update password ////////////////////

func (s *svc) UpdatePassword(ctx context.Context, id uuid.UUID, input dto.UpdatePasswordInput) error {
	log.Info().Msgf("update password: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update password")
		return err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to update password")
		return service.ErrUserNotFound
	}

	// Check old password
	if !security.CheckPasswordHash(input.OldPassword, userEntity.Password) {
		log.Error().Stack().Err(service.ErrIncorrectOldPassword).Msg("failed to update password")
		return service.ErrIncorrectOldPassword
	}

	// Hash new password
	userEntity.Password, err = security.HashPassword(input.NewPassword)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update password")
		return err
	}

	// Update password
	userEntity.IsPasswordTemp = false
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to update password")
		return err
	}

	return nil
}

//////////////////// Restore account ////////////////////

func (s *svc) RestoreAccount(ctx context.Context, id uuid.UUID) (dto.AccountOutput, error) {
	log.Info().Msgf("restore account: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to restore account")
		return dto.AccountOutput{}, err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to restore account")
		return dto.AccountOutput{}, service.ErrUserNotFound
	}

	// Check if user is not deleted
	if userEntity.DeletedAt == nil {
		log.Error().Stack().Err(service.ErrUserNotDeleted).Msg("failed to restore account")
		return dto.AccountOutput{}, service.ErrUserNotDeleted
	}

	// Restore
	userEntity.DeletedAt = nil
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to restore account")
		return dto.AccountOutput{}, err
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Delete account ////////////////////

func (s *svc) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	log.Info().Msgf("delete account: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to delete account")
		return err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to delete account")
		return service.ErrUserNotFound
	}

	// Check if user already is deleted
	if userEntity.DeletedAt != nil {
		log.Error().Stack().Err(service.ErrUserAlreadyDeleted).Msg("failed to delete account")
		return service.ErrUserAlreadyDeleted
	}

	// Delete
	now := time.Now()
	userEntity.DeletedAt = &now
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to delete account")
		return err
	}

	return nil
}
