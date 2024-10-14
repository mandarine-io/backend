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
	"mandarine/internal/api/service/auth/mapper"
	"mandarine/pkg/logging"
	dto2 "mandarine/pkg/rest/dto"
	"mandarine/pkg/rest/middleware"
	"mandarine/pkg/smtp"
	"mandarine/pkg/storage/cache/manager"
	"mandarine/pkg/template"
	"time"
)

const (
	registerCachePrefix       = "register"
	registerEmailDefaultTitle = "Activate your account"
)

var (
	ErrDuplicateUser       = dto2.NewI18nError("duplicate user", "errors.duplicate_user")
	ErrInvalidOrExpiredOtp = dto2.NewI18nError("OTP is invalid or has expired", "errors.invalid_or_expired_otp")
)

type RegisterService struct {
	userRepo       repo.UserRepository
	cacheManager   manager.CacheManager
	smtpSender     smtp.Sender
	templateEngine template.Engine
	cfg            *config.Config
}

func NewRegisterService(
	userRepo repo.UserRepository,
	cacheManager manager.CacheManager,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	cfg *config.Config,
) *RegisterService {
	return &RegisterService{
		userRepo:       userRepo,
		cacheManager:   cacheManager,
		smtpSender:     smtpSender,
		templateEngine: templateEngine,
		cfg:            cfg,
	}
}

//////////////////// Register ////////////////////

func (s *RegisterService) Register(ctx context.Context, input dto.RegisterInput) error {
	slog.Info("Register")
	factoryErr := func(err error) error {
		slog.Error("Register error", logging.ErrorAttr(err))
		return err
	}

	// Check if user exists
	exists, err := s.userRepo.ExistsUserByUsernameOrEmail(ctx, input.Username, input.Email)
	if err != nil {
		return factoryErr(err)
	}
	if exists {
		return factoryErr(ErrDuplicateUser)
	}

	// Hashing password
	input.Password, err = security.HashPassword(input.Password)
	if err != nil {
		return factoryErr(err)
	}

	// Generate OTP
	otp, err := security.GenerateOTP(s.cfg.Security.OTP.Length)
	if err != nil {
		return factoryErr(err)
	}

	// Save in cache
	expiration := time.Duration(s.cfg.Security.OTP.TTL) * time.Second
	cacheEntry := dto.RegisterCache{
		User:      input,
		OTP:       otp,
		ExpiredAt: time.Now().Add(expiration),
	}
	err = s.cacheManager.SetWithExpiration(ctx, cache.CreateCacheKey(registerCachePrefix, input.Email), cacheEntry, expiration)
	if err != nil {
		return factoryErr(err)
	}

	// Get localizer
	localizer := ctx.Value(middleware.LocalizerKey)
	emailTitle := registerEmailDefaultTitle
	if localizer != nil {
		switch localizer := localizer.(type) {
		case *i18n.Localizer:
			emailTitle = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email.register-confirm.title"})
		}
	}

	// Send mail
	args := dto.RegisterConfirmTemplateArgs{
		Email: input.Email,
		TTL:   s.cfg.Security.OTP.TTL / 60,
		OTP:   otp,
	}
	content, err := s.templateEngine.Render("register-confirm", args)
	if err != nil {
		return factoryErr(err)
	}

	err = s.smtpSender.SendHtmlMessage(emailTitle, content, input.Email)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}

//////////////////// Register confirmation ////////////////////

func (s *RegisterService) RegisterConfirm(ctx context.Context, input dto.RegisterConfirmInput) error {
	slog.Info("Register confirm")
	factoryErr := func(err error) error {
		slog.Error("Register confirm error", logging.ErrorAttr(err))
		return err
	}

	// Get entry from cache
	var cacheEntry dto.RegisterCache
	err := s.cacheManager.Get(ctx, cache.CreateCacheKey(registerCachePrefix, input.Email), &cacheEntry)
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

	// Check email
	if input.Email != cacheEntry.User.Email {
		return factoryErr(ErrInvalidOrExpiredOtp)
	}

	// Check if user exists
	exists, err := s.userRepo.ExistsUserByUsernameOrEmail(ctx, cacheEntry.User.Username, cacheEntry.User.Email)
	if err != nil {
		return factoryErr(err)
	}
	if exists {
		return factoryErr(ErrDuplicateUser)
	}

	// Create user in DB
	registerRequest := cacheEntry.User
	userEntity := mapper.MapRegisterRequestToUserEntity(registerRequest)
	_, err = s.userRepo.CreateUser(ctx, userEntity)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicateUser) {
			return factoryErr(ErrDuplicateUser)
		}
		return factoryErr(err)
	}

	// Delete cache
	err = s.cacheManager.Delete(ctx, cache.CreateCacheKey(registerCachePrefix, input.Email))
	if err != nil {
		slog.Warn("Register confirm error", logging.ErrorAttr(err))
	}

	return nil
}
