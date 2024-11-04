package account

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/api/config"
	cache2 "github.com/mandarine-io/Backend/internal/api/helper/cache"
	"github.com/mandarine-io/Backend/internal/api/helper/random"
	"github.com/mandarine-io/Backend/internal/api/helper/security"
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	"github.com/mandarine-io/Backend/internal/api/service/account/dto"
	"github.com/mandarine-io/Backend/internal/api/service/account/mapper"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/mandarine-io/Backend/pkg/template"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	emailVerifyCachePrefix = "email-verify"
	emailDefaultTitle      = "Verify email"
)

var (
	ErrInvalidOrExpiredOtp  = dto2.NewI18nError("invalid or expired otp", "errors.invalid_or_expired_otp")
	ErrUserNotFound         = dto2.NewI18nError("user not found", "errors.user_not_found")
	ErrDuplicateUsername    = dto2.NewI18nError("username already in use", "errors.duplicate_username")
	ErrDuplicateEmail       = dto2.NewI18nError("email already in use", "errors.duplicate_email")
	ErrPasswordIsSet        = dto2.NewI18nError("password is already set", "errors.password_is_set")
	ErrIncorrectOldPassword = dto2.NewI18nError("incorrect old password", "errors.incorrect_old_password")
	ErrUserNotDeleted       = dto2.NewI18nError("user not deleted", "errors.user_not_deleted")
	ErrUserAlreadyDeleted   = dto2.NewI18nError("user already deleted", "errors.user_already_deleted")
	ErrSendEmail            = dto2.NewI18nError("failed to send email", "errors.failed_to_send_email")
)

type Service struct {
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
) *Service {
	return &Service{
		userRepo:       userRepo,
		cacheManager:   cacheManager,
		smtpSender:     smtpSender,
		templateEngine: templateEngine,
		cfg:            cfg,
	}
}

//////////////////// Get account ////////////////////

func (s *Service) GetAccount(ctx context.Context, id uuid.UUID) (dto.AccountOutput, error) {
	log.Info().Msgf("get account: %s", id.String())
	factoryErr := func(err error) (dto.AccountOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to get account")
		return dto.AccountOutput{}, err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Update username ////////////////////

func (s *Service) UpdateUsername(
	ctx context.Context, id uuid.UUID, input dto.UpdateUsernameInput,
) (dto.AccountOutput, error) {
	log.Info().Msgf("update username: %s", id.String())
	factoryErr := func(err error) (dto.AccountOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to update account")
		return dto.AccountOutput{}, err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check if username not changed
	if input.Username == userEntity.Username {
		return mapper.MapUserEntityToAccountResponse(userEntity), nil
	}

	// Check if username is already in use
	exists, err := s.userRepo.ExistsUserByUsername(ctx, input.Username)
	if err != nil {
		return factoryErr(err)
	}
	if exists {
		return factoryErr(ErrDuplicateUsername)
	}

	// Update username
	userEntity.Username = input.Username
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Update email ////////////////////

func (s *Service) UpdateEmail(
	ctx context.Context, id uuid.UUID, input dto.UpdateEmailInput, localizer *i18n.Localizer,
) (dto.AccountOutput, error) {
	log.Info().Msgf("update email: %s", id.String())
	factoryErr := func(err error) (dto.AccountOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to update account")
		return dto.AccountOutput{}, err
	}
	factoryChildErr := func(err error, childErr error) (dto.AccountOutput, error) {
		log.Error().Stack().Err(childErr).Msg("failed to update account")
		return dto.AccountOutput{}, err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check if email not changed
	if input.Email == userEntity.Email {
		return mapper.MapUserEntityToAccountResponse(userEntity), nil
	}

	// Check if email is already in use
	exists, err := s.userRepo.ExistsUserByEmail(ctx, input.Email)
	if err != nil {
		return factoryErr(err)
	}
	if exists {
		return factoryErr(ErrDuplicateEmail)
	}

	// Generate OTP
	otp, err := random.GenerateRandomNumber(s.cfg.Security.OTP.Length)
	if err != nil {
		return factoryErr(err)
	}

	// Save in cache
	expiration := time.Duration(s.cfg.Security.OTP.TTL) * time.Second
	cacheEntry := dto.EmailVerifyCache{
		Email:     input.Email,
		OTP:       otp,
		ExpiredAt: time.Now().Add(expiration),
	}
	err = s.cacheManager.SetWithExpiration(
		ctx, cache2.CreateCacheKey(emailVerifyCachePrefix, input.Email), cacheEntry, expiration,
	)
	if err != nil {
		return factoryErr(err)
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
		return factoryErr(err)
	}

	err = s.smtpSender.SendHtmlMessage(emailTitle, content, input.Email)
	if err != nil {
		return factoryChildErr(ErrSendEmail, err)
	}

	// Update email
	userEntity.Email = input.Email
	userEntity.IsEmailVerified = false
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Verify email ////////////////////

func (s *Service) VerifyEmail(ctx context.Context, id uuid.UUID, req dto.VerifyEmailInput) error {
	log.Info().Msgf("verify email: %s", id.String())
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to verify email")
		return err
	}

	// Get entry from cache
	var cacheEntry dto.EmailVerifyCache
	err := s.cacheManager.Get(ctx, cache2.CreateCacheKey(emailVerifyCachePrefix, req.Email), &cacheEntry)
	if err != nil {
		if errors.Is(err, cache.ErrCacheEntryNotFound) {
			return factoryErr(ErrInvalidOrExpiredOtp)
		}
		return factoryErr(err)
	}

	// Check OTP
	if req.OTP != cacheEntry.OTP {
		return factoryErr(ErrInvalidOrExpiredOtp)
	}

	// Get user entity by id
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check email
	if userEntity.Email != cacheEntry.Email {
		return factoryErr(ErrInvalidOrExpiredOtp)
	}

	// Verify email
	userEntity.IsEmailVerified = true
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	// Delete cache entry
	err = s.cacheManager.Delete(ctx, cache2.CreateCacheKey(emailVerifyCachePrefix, req.Email))
	if err != nil {
		log.Warn().Stack().Err(err).Msg("failed to delete cache entry")
	}

	return nil
}

//////////////////// Set password ////////////////////

func (s *Service) SetPassword(ctx context.Context, id uuid.UUID, input dto.SetPasswordInput) error {
	log.Info().Msgf("set password: %s", id.String())
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to set password")
		return err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check if password is empty
	if !userEntity.IsPasswordTemp {
		return factoryErr(ErrPasswordIsSet)
	}

	// Hash new password
	userEntity.Password, err = security.HashPassword(input.Password)
	if err != nil {
		return factoryErr(err)
	}

	// Update password
	userEntity.IsPasswordTemp = false
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}

//////////////////// Update password ////////////////////

func (s *Service) UpdatePassword(ctx context.Context, id uuid.UUID, input dto.UpdatePasswordInput) error {
	log.Info().Msgf("update password: %s", id.String())
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to update password")
		return err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check old password
	if !security.CheckPasswordHash(input.OldPassword, userEntity.Password) {
		return factoryErr(ErrIncorrectOldPassword)
	}

	// Hash new password
	userEntity.Password, err = security.HashPassword(input.NewPassword)
	if err != nil {
		return factoryErr(err)
	}

	// Update password
	userEntity.IsPasswordTemp = false
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}

//////////////////// Restore account ////////////////////

func (s *Service) RestoreAccount(ctx context.Context, id uuid.UUID) (dto.AccountOutput, error) {
	log.Info().Msgf("restore account: %s", id.String())
	factoryErr := func(err error) (dto.AccountOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to restore account")
		return dto.AccountOutput{}, err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check if user is not deleted
	if userEntity.DeletedAt == nil {
		return factoryErr(ErrUserNotDeleted)
	}

	// Restore
	userEntity.DeletedAt = nil
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return mapper.MapUserEntityToAccountResponse(userEntity), nil
}

//////////////////// Delete account ////////////////////

func (s *Service) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	log.Info().Msgf("delete account: %s", id.String())
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to delete account")
		return err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, id, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check if user already is deleted
	if userEntity.DeletedAt != nil {
		return factoryErr(ErrUserAlreadyDeleted)
	}

	// Delete
	now := time.Now()
	userEntity.DeletedAt = &now
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}
