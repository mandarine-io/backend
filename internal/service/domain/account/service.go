package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/converter"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/internal/infrastructure/smtp"
	"github.com/mandarine-io/backend/internal/infrastructure/template"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/service/domain"
	infra "github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/internal/util/security"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"time"
)

const (
	emailVerifyCachePrefix = "email_verify"
	emailDefaultTitle      = "Verify email"
)

type svc struct {
	userRepo       repo.UserRepository
	smtpSender     smtp.Sender
	templateEngine template.Engine
	otpService     infra.OTPService
	cfg            config.Config
	logger         zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(
	cfg config.Config,
	userRepo repo.UserRepository,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	otpService infra.OTPService,
	opts ...Option,
) domain.AccountService {
	s := &svc{
		userRepo:       userRepo,
		smtpSender:     smtpSender,
		templateEngine: templateEngine,
		otpService:     otpService,
		cfg:            cfg,
		logger:         zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

//////////////////// Get account ////////////////////

func (s *svc) GetAccount(ctx context.Context, id uuid.UUID) (v0.AccountOutput, error) {
	s.logger.Info().Msgf("get account: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.AccountOutput{}, err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return v0.AccountOutput{}, domain.ErrUserNotFound
	}

	return converter.MapUserEntityToAccountOutput(userEntity), nil
}

//////////////////// Update username ////////////////////

func (s *svc) UpdateUsername(
	ctx context.Context, id uuid.UUID, input v0.UpdateUsernameInput,
) (v0.AccountOutput, error) {
	s.logger.Info().Msgf("update username: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.AccountOutput{}, err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return v0.AccountOutput{}, domain.ErrUserNotFound
	}

	// Check if username not changed
	if input.Username == userEntity.Username {
		return converter.MapUserEntityToAccountOutput(userEntity), nil
	}

	// Check if username is already in use
	exists, err := s.userRepo.ExistsUserByUsername(ctx, input.Username)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check exist user by username")
		return v0.AccountOutput{}, err
	}
	if exists {
		s.logger.Error().Stack().Err(domain.ErrDuplicateUsername).Msg("user with such username already exists")
		return v0.AccountOutput{}, domain.ErrDuplicateUsername
	}

	// Update username
	userEntity.Username = input.Username

	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return v0.AccountOutput{}, err
	}

	return converter.MapUserEntityToAccountOutput(userEntity), nil
}

//////////////////// Update email ////////////////////

func (s *svc) UpdateEmail(
	ctx context.Context, id uuid.UUID, input v0.UpdateEmailInput, localizer locale.Localizer,
) (v0.AccountOutput, error) {
	s.logger.Info().Msgf("update email: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.AccountOutput{}, err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return v0.AccountOutput{}, domain.ErrUserNotFound
	}

	// Check if email not changed
	if input.Email == userEntity.Email {
		return converter.MapUserEntityToAccountOutput(userEntity), nil
	}

	// Check if email is already in use
	exists, err := s.userRepo.ExistsUserByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check exist user by email")
		return v0.AccountOutput{}, err
	}
	if exists {
		s.logger.Error().Stack().Err(domain.ErrDuplicateEmail).Msg("user with such email already exists")
		return v0.AccountOutput{}, domain.ErrDuplicateEmail
	}

	// Create and save OTP
	otp, err := s.otpService.GenerateAndSaveWithCode(ctx, emailVerifyCachePrefix, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to create and save OTP code")
		return v0.AccountOutput{}, err
	}

	// Localize email title
	emailTitle := emailDefaultTitle
	if localizer != nil {
		emailTitle = localizer.Localize("email.email-verify.title", nil, 0)
	}

	// Send mail
	args := v0.EmailVerifyTemplateArgs{
		Email: input.Email,
		TTL:   s.cfg.Security.OTP.TTL / 60,
		OTP:   otp,
	}
	content, err := s.templateEngine.RenderHTML("email-verify", args)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to render HTML email")
		return v0.AccountOutput{}, err
	}

	err = s.smtpSender.SendHTMLMessage(emailTitle, content, s.cfg.SMTP.From, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to send HTML message")
		return v0.AccountOutput{}, domain.ErrSendEmail
	}

	// Update email
	userEntity.Email = input.Email
	userEntity.IsEmailVerified = false

	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return v0.AccountOutput{}, err
	}

	return converter.MapUserEntityToAccountOutput(userEntity), nil
}

//////////////////// Verify email ////////////////////

func (s *svc) VerifyEmail(ctx context.Context, id uuid.UUID, input v0.VerifyEmailInput) error {
	s.logger.Info().Msgf("verify email: %s", id.String())

	// Get entry from cache
	var email string
	err := s.otpService.GetDataByCode(ctx, emailVerifyCachePrefix, input.OTP, &email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get data by OTP")
		return err
	}

	// Check OTP
	if input.Email != email {
		s.logger.Error().Stack().Err(infra.ErrInvalidOrExpiredOTP).Msg("dont match emails in request and OTP data")
		return infra.ErrInvalidOrExpiredOTP
	}

	// Get user entity by id
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return domain.ErrUserNotFound
	}

	// Check email
	if userEntity.Email != email {
		s.logger.Error().Stack().Err(infra.ErrInvalidOrExpiredOTP).Msg("dont match emails in request and user data")
		return infra.ErrInvalidOrExpiredOTP
	}

	// Verify email
	userEntity.IsEmailVerified = true

	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return err
	}

	// Delete cache entry
	err = s.otpService.DeleteDataByCode(ctx, emailVerifyCachePrefix, input.OTP)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to delete OTP data")
	}

	return nil
}

//////////////////// Set password ////////////////////

func (s *svc) SetPassword(ctx context.Context, id uuid.UUID, input v0.SetPasswordInput) error {
	s.logger.Info().Msgf("set password: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return domain.ErrUserNotFound
	}

	// Check if password is empty
	if !userEntity.IsPasswordTemp {
		s.logger.Error().Stack().Err(domain.ErrPasswordIsSet).Msg("password is temporary")
		return domain.ErrPasswordIsSet
	}

	// Hash new password
	userEntity.Password, err = security.HashPassword(input.Password)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to hash password")
		return err
	}

	// Update password
	userEntity.IsPasswordTemp = false

	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return err
	}

	return nil
}

//////////////////// Update password ////////////////////

func (s *svc) UpdatePassword(ctx context.Context, id uuid.UUID, input v0.UpdatePasswordInput) error {
	s.logger.Info().Msgf("update password: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return domain.ErrUserNotFound
	}

	// Check old password
	if !security.CheckPasswordHash(input.OldPassword, userEntity.Password) {
		s.logger.Error().Stack().Err(domain.ErrIncorrectOldPassword).Msg("failed to check password")
		return domain.ErrIncorrectOldPassword
	}

	// Hash new password
	userEntity.Password, err = security.HashPassword(input.NewPassword)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to hash password")
		return err
	}

	// Update password
	userEntity.IsPasswordTemp = false
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return err
	}

	return nil
}

//////////////////// Restore account ////////////////////

func (s *svc) RestoreAccount(ctx context.Context, id uuid.UUID) (v0.AccountOutput, error) {
	s.logger.Info().Msgf("restore account: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.AccountOutput{}, err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return v0.AccountOutput{}, domain.ErrUserNotFound
	}

	// Check if user is not deleted
	if userEntity.DeletedAt == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotDeleted).Msg("user not deleted")
		return v0.AccountOutput{}, domain.ErrUserNotDeleted
	}

	// Restore
	userEntity.DeletedAt = nil
	userEntity, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return v0.AccountOutput{}, err
	}

	return converter.MapUserEntityToAccountOutput(userEntity), nil
}

//////////////////// Delete account ////////////////////

func (s *svc) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	s.logger.Info().Msgf("delete account: %s", id.String())

	// Get user entity
	userEntity, err := s.userRepo.FindUserByID(ctx, id)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return err
	}
	if userEntity == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return domain.ErrUserNotFound
	}

	// Check if user already is deleted
	if userEntity.DeletedAt != nil {
		s.logger.Error().Stack().Err(domain.ErrUserAlreadyDeleted).Msg("user already deleted")
		return domain.ErrUserAlreadyDeleted
	}

	// Delete
	now := time.Now()
	userEntity.DeletedAt = &now
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to delete account")
		return err
	}

	return nil
}
