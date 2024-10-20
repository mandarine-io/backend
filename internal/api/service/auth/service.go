package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/cache"
	"mandarine/internal/api/helper/random"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/persistence/repo"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/internal/api/service/auth/mapper"
	"mandarine/pkg/logging"
	"mandarine/pkg/oauth"
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

	recoveryPasswordCachePrefix = "recovery_password"
	recoveryEmailDefaultTitle   = "Recovery password"
)

var (
	ErrDuplicateUser       = dto2.NewI18nError("duplicate user", "errors.duplicate_user")
	ErrInvalidOrExpiredOtp = dto2.NewI18nError("OTP is invalid or has expired", "errors.invalid_or_expired_otp")
	ErrBadCredentials      = dto2.NewI18nError("bad credentials", "errors.bad_credentials")
	ErrUserIsBlocked       = dto2.NewI18nError("user is blocked", "errors.user_is_blocked")
	ErrUserNotFound        = dto2.NewI18nError("user not found", "errors.user_not_found")
	ErrInvalidJwtToken     = dto2.NewI18nError("invalid JWT token", "errors.invalid_jwt_token")
	ErrUserInfoNotReceived = dto2.NewI18nError("user info not received", "errors.userinfo_not_received")
	ErrInvalidProvider     = dto2.NewI18nError("invalid provider", "errors.invalid_provider")
	ErrSendEmail           = dto2.NewI18nError("failed to send email", "errors.failed_to_send_email")
)

type Service struct {
	userRepo        repo.UserRepository
	bannedTokenRepo repo.BannedTokenRepository
	oauthProviders  map[string]oauth.Provider
	cacheManager    manager.CacheManager
	smtpSender      smtp.Sender
	templateEngine  template.Engine
	cfg             *config.Config
}

func NewService(
	userRepo repo.UserRepository,
	bannedTokenRepo repo.BannedTokenRepository,
	oauthProviders map[string]oauth.Provider,
	cacheManager manager.CacheManager,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	cfg *config.Config,
) *Service {
	return &Service{
		userRepo:        userRepo,
		bannedTokenRepo: bannedTokenRepo,
		oauthProviders:  oauthProviders,
		cacheManager:    cacheManager,
		smtpSender:      smtpSender,
		templateEngine:  templateEngine,
		cfg:             cfg,
	}
}

//////////////////// Register ////////////////////

func (s *Service) Register(ctx context.Context, input dto.RegisterInput) error {
	slog.Info("Register")
	factoryErr := func(err error) error {
		slog.Error("Register error", logging.ErrorAttr(err))
		return err
	}
	factoryChildErr := func(err error, childErr error) error {
		slog.Error("Register error", logging.ErrorAttr(childErr))
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
	otp, err := random.GenerateRandomNumber(s.cfg.Security.OTP.Length)
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
		return factoryChildErr(ErrSendEmail, err)
	}

	return nil
}

//////////////////// Register confirmation ////////////////////

func (s *Service) RegisterConfirm(ctx context.Context, input dto.RegisterConfirmInput) error {
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

//////////////////// Login ////////////////////

func (s *Service) Login(ctx context.Context, input dto.LoginInput) (dto.JwtTokensOutput, error) {
	slog.Info("Login")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		slog.Error("Login error", logging.ErrorAttr(err))
		return dto.JwtTokensOutput{}, err
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserByUsernameOrEmail(ctx, input.Login, true)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check password
	if !security.CheckPasswordHash(input.Password, userEntity.Password) {
		return factoryErr(ErrBadCredentials)
	}

	// Check if user is blocked
	if !userEntity.IsEnabled {
		return factoryErr(ErrUserIsBlocked)
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Refresh Tokens ////////////////////

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (dto.JwtTokensOutput, error) {
	slog.Info("RefreshTokens tokens")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		slog.Error("RefreshTokens tokens error", logging.ErrorAttr(err))
		return dto.JwtTokensOutput{}, err
	}

	// Check token
	token, err := security.DecodeAndValidateJwtToken(refreshToken, s.cfg.Security.JWT.Secret)
	if err != nil {
		return factoryErr(ErrInvalidJwtToken)
	}

	// Get user ID from token
	claims, err := security.GetClaimsFromJwtToken(token)
	if err != nil {
		return factoryErr(ErrInvalidJwtToken)
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return factoryErr(ErrInvalidJwtToken)
	}

	userUUID, err := uuid.Parse(sub)
	if err != nil {
		return factoryErr(ErrInvalidJwtToken)
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, userUUID, true)
	if err != nil {
		return factoryErr(err)
	}
	if userEntity == nil {
		return factoryErr(ErrUserNotFound)
	}

	// Check if user is blocked
	if !userEntity.IsEnabled {
		return factoryErr(ErrUserIsBlocked)
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Logout ////////////////////

func (s *Service) Logout(ctx context.Context, jti string) error {
	bannedToken := &model.BannedTokenEntity{
		JTI:       jti,
		ExpiredAt: time.Now().Add(time.Duration(s.cfg.Security.JWT.RefreshTokenTTL) * time.Second).Unix(),
	}
	_, err := s.bannedTokenRepo.CreateOrUpdateBannedToken(ctx, bannedToken)
	return err
}

//////////////////// Recovery password ////////////////////

func (s *Service) RecoveryPassword(ctx context.Context, input dto.RecoveryPasswordInput) error {
	slog.Info("Recovery password")
	factoryErr := func(err error) error {
		slog.Error("Recovery password error", logging.ErrorAttr(err))
		return err
	}
	factoryChildErr := func(err error, childErr error) error {
		slog.Error("Recovery password error", logging.ErrorAttr(childErr))
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
	otp, err := random.GenerateRandomNumber(s.cfg.Security.OTP.Length)
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
		return factoryChildErr(ErrSendEmail, err)
	}

	return nil
}

//////////////////// Verify recovery password ////////////////////

func (s *Service) VerifyRecoveryCode(ctx context.Context, input dto.VerifyRecoveryCodeInput) error {
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

func (s *Service) ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error {
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

//////////////////// Get consent page url ////////////////////

func (s *Service) GetConsentPageUrl(_ context.Context, provider string, redirectUrl string) (dto.GetConsentPageUrlOutput, error) {
	slog.Info(fmt.Sprintf("Get consent page url: provider=%s", provider))

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return dto.GetConsentPageUrlOutput{}, ErrInvalidProvider
	}

	consentPageUrl, oauthState := oauthProvider.GetConsentPageUrl(redirectUrl)
	return dto.GetConsentPageUrlOutput{ConsentPageUrl: consentPageUrl, OauthState: oauthState}, nil
}

//////////////////// Fetch user info ////////////////////

func (s *Service) FetchUserInfo(ctx context.Context, provider string, input dto.FetchUserInfoInput) (
	oauth.UserInfo, error,
) {
	factoryErr := func(err error) (oauth.UserInfo, error) {
		slog.Error("Get user info error", logging.ErrorAttr(err))
		return oauth.UserInfo{}, err
	}

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return oauth.UserInfo{}, ErrInvalidProvider
	}

	// Exchange code to token
	slog.Info("Exchange code to token")
	socialLoginCallbackUrl := fmt.Sprintf("%s/auth/social/%s/callback/", s.cfg.Server.ExternalOrigin, provider)
	token, err := oauthProvider.ExchangeCodeToToken(ctx, input.Code, socialLoginCallbackUrl)
	if err != nil {
		return factoryErr(err)
	}

	// Get user info
	slog.Info("Get user info")
	userInfo, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		if errors.Is(err, oauth.ErrUserInfoNotReceived) {
			return factoryErr(ErrUserInfoNotReceived)
		}
		return factoryErr(err)
	}

	return userInfo, nil
}

//////////////////// Register or login ////////////////////

func (s *Service) RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (dto.JwtTokensOutput, error) {
	slog.Info("Register and login")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		slog.Error("Register and login error", logging.ErrorAttr(err))
		return dto.JwtTokensOutput{}, err
	}

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, userInfo.Email, true)
	if err != nil {
		return factoryErr(err)
	}

	// Save user
	if userEntity == nil {
		slog.Info("Create new user")
		userInfo.Username, err = s.searchUniqueUsername(ctx, userInfo.Username)
		if err != nil {
			return factoryErr(err)
		}

		userEntity = mapper.MapUserInfoToUserEntity(userInfo)
		userEntity, err = s.userRepo.CreateUser(ctx, userEntity)
		if err != nil {
			return factoryErr(err)
		}
	} else if !userEntity.IsEnabled {
		return factoryErr(ErrUserIsBlocked)
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		return factoryErr(err)
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Additional helpers ////////////////////

func (s *Service) searchUniqueUsername(ctx context.Context, defaultUsername string) (string, error) {
	slog.Info("Search unique username")
	exists, err := s.userRepo.ExistsUserByUsername(ctx, defaultUsername)
	if err != nil {
		return "", err
	}
	if !exists {
		return defaultUsername, nil
	}

	for {
		suffix, err := random.GenerateRandomNumber(10)
		if err != nil {
			return "", err
		}

		username := fmt.Sprintf("%s_%s", defaultUsername, suffix)
		exists, err := s.userRepo.ExistsUserByUsername(ctx, username)
		if err != nil {
			return "", err
		}
		if !exists {
			return username, nil
		}

		slog.Debug("Unique username not found")
	}
}
