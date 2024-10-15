package auth

import (
	"context"
	"errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/cache"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/repo"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/pkg/logging"
	"mandarine/pkg/rest/middleware"
	"mandarine/pkg/smtp"
	"mandarine/pkg/storage/cache/manager"
	"mandarine/pkg/template"
	"time"
)

const (
	recoveryPasswordCachePrefix = "recovery_password"
	recoveryEmailDefaultTitle   = "Recovery password"
)

type ResetPasswordService struct {
	userRepo       repo.UserRepository
	cacheManager   manager.CacheManager
	smtpSender     smtp.Sender
	templateEngine template.Engine
	cfg            *config.Config
}

func NewResetPasswordService(
	userRepo repo.UserRepository,
	cacheManager manager.CacheManager,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	cfg *config.Config,
) *ResetPasswordService {
	return &ResetPasswordService{
		userRepo:       userRepo,
		cacheManager:   cacheManager,
		smtpSender:     smtpSender,
		templateEngine: templateEngine,
		cfg:            cfg,
	}
}

//////////////////// Recovery password ////////////////////

func (s *ResetPasswordService) RecoveryPassword(ctx context.Context, input dto.RecoveryPasswordInput) error {
	slog.Info("Recovery password")
	factoryErr := func(err error) error {
		slog.Error("Recovery password error", logging.ErrorAttr(err))
		return err
	}

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, input.Email, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Generate OTP
	otp, err := security.GenerateOTP(s.cfg.Security.OTP.Length)
	if err != nil {
		return factoryErr(err)
	}

	// Save in cache
	expiration := time.Duration(s.cfg.Security.OTP.TTL) * time.Second
	cacheEntry := dto.RecoveryPasswordCache{
		Email:     input.Email,
		OTP:       otp,
		ExpiredAt: time.Now().Add(expiration),
	}

	err = s.cacheManager.SetWithExpiration(ctx, cache.CreateCacheKey(recoveryPasswordCachePrefix, input.Email), cacheEntry, expiration)
	if err != nil {
		return factoryErr(err)
	}

	// Get localizer
	localizer := ctx.Value(middleware.LocalizerKey)
	emailTitle := recoveryEmailDefaultTitle
	if localizer != nil {
		switch localizer := localizer.(type) {
		case *i18n.Localizer:
			emailTitle = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email.recovery-password.title"})
		}
	}

	// Send email
	args := dto.RegisterConfirmTemplateArgs{
		Email: input.Email,
		TTL:   s.cfg.Security.OTP.TTL / 60,
		OTP:   otp,
	}
	content, err := s.templateEngine.Render("recovery-password", args)
	if err != nil {
		return factoryErr(err)
	}

	err = s.smtpSender.SendHtmlMessage(emailTitle, content, input.Email)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}

//////////////////// Verify recovery password ////////////////////

func (s *ResetPasswordService) VerifyRecoveryCode(ctx context.Context, input dto.VerifyRecoveryCodeInput) error {
	slog.Info("Verify recovery password")
	factoryErr := func(err error) error {
		slog.Error("Verify recovery password error", logging.ErrorAttr(err))
		return err
	}

	// Get entry from cache
	var cacheEntry dto.RecoveryPasswordCache
	err := s.cacheManager.Get(ctx, cache.CreateCacheKey(recoveryPasswordCachePrefix, input.Email), &cacheEntry)
	if err != nil {
		if errors.Is(err, manager.ErrCacheEntryNotFound) {
			return factoryErr(ErrInvalidOrExpiredOtp)
		}
		return factoryErr(err)
	}

	// Check OTP
	if input.OTP != cacheEntry.OTP {
		return factoryErr(ErrInvalidOrExpiredOtp)
	}

	return nil
}

//////////////////// Reset password ////////////////////

func (s *ResetPasswordService) ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error {
	slog.Info("Reset password")
	factoryErr := func(err error) error {
		slog.Error("Reset password error", logging.ErrorAttr(err))
		return err
	}

	// Get entry from cache
	var cacheEntry dto.RecoveryPasswordCache
	err := s.cacheManager.Get(ctx, cache.CreateCacheKey(recoveryPasswordCachePrefix, input.Email), &cacheEntry)
	if err != nil {
		if errors.Is(err, manager.ErrCacheEntryNotFound) {
			return factoryErr(ErrInvalidOrExpiredOtp)
		}
		return factoryErr(err)
	}

	// Check OTP
	if input.OTP != cacheEntry.OTP {
		return factoryErr(ErrInvalidOrExpiredOtp)
	}

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, input.Email, false)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Hashing password
	userEntity.Password, err = security.HashPassword(input.Password)
	if err != nil {
		return factoryErr(err)
	}

	// Update password
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}
