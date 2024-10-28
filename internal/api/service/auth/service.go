package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/internal/api/helper/cache"
	"github.com/mandarine-io/Backend/internal/api/helper/random"
	"github.com/mandarine-io/Backend/internal/api/helper/security"
	"github.com/mandarine-io/Backend/internal/api/persistence/model"
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	"github.com/mandarine-io/Backend/internal/api/service/auth/dto"
	"github.com/mandarine-io/Backend/internal/api/service/auth/mapper"
	"github.com/mandarine-io/Backend/pkg/oauth"
	dto2 "github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/mandarine-io/Backend/pkg/storage/cache/manager"
	"github.com/mandarine-io/Backend/pkg/template"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog/log"
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

func (s *Service) Register(ctx context.Context, input dto.RegisterInput, localizer *i18n.Localizer) error {
	log.Info().Msgf("register: %s", input.Username)
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to register user")
		return err
	}
	factoryChildErr := func(err error, childErr error) error {
		log.Error().Stack().Err(childErr).Msg("failed to register user")
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

	// Localize email title
	emailTitle := registerEmailDefaultTitle
	if localizer != nil {
		emailTitle = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email.register-confirm.title"})
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
	log.Info().Msg("register confirm")
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to register confirm user")
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
		log.Warn().Stack().Err(err).Msg("failed to delete cache")
	}

	return nil
}

//////////////////// Login ////////////////////

func (s *Service) Login(ctx context.Context, input dto.LoginInput) (dto.JwtTokensOutput, error) {
	log.Info().Msg("login")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to login")
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
	log.Info().Msg("refresh tokens")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, err
	}

	// Check token
	token, err := security.DecodeAndValidateJwtToken(refreshToken, s.cfg.Security.JWT.Secret)
	if err != nil {
		return factoryErr(ErrInvalidJwtToken)
	}

	// Get user id from token
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
	log.Info().Msg("logout")
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to logout")
		return err
	}

	// Create banned token
	bannedToken := &model.BannedTokenEntity{
		JTI:       jti,
		ExpiredAt: time.Now().Add(time.Duration(s.cfg.Security.JWT.RefreshTokenTTL) * time.Second).Unix(),
	}
	_, err := s.bannedTokenRepo.CreateOrUpdateBannedToken(ctx, bannedToken)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}

//////////////////// Recovery password ////////////////////

func (s *Service) RecoveryPassword(ctx context.Context, input dto.RecoveryPasswordInput, localizer *i18n.Localizer) error {
	log.Info().Msg("recovery password")
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to recovery password")
		return err
	}
	factoryChildErr := func(err error, childErr error) error {
		log.Error().Stack().Err(childErr).Err(err).Msg("failed to recovery password")
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

	// Localize email title
	emailTitle := recoveryEmailDefaultTitle
	if localizer != nil {
		emailTitle = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "email.recovery-password.title"})
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
	log.Info().Msg("verify recovery password")
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to verify recovery password")
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
	log.Info().Msg("reset password")
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to reset password")
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
	log.Info().Msgf("get consent page url: provider=%s", provider)
	factoryErr := func(err error) (dto.GetConsentPageUrlOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to get consent page url")
		return dto.GetConsentPageUrlOutput{}, err
	}

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return factoryErr(ErrInvalidProvider)
	}

	consentPageUrl, oauthState := oauthProvider.GetConsentPageUrl(redirectUrl)
	return dto.GetConsentPageUrlOutput{ConsentPageUrl: consentPageUrl, OauthState: oauthState}, nil
}

//////////////////// Fetch user info ////////////////////

func (s *Service) FetchUserInfo(ctx context.Context, provider string, input dto.FetchUserInfoInput) (
	oauth.UserInfo, error,
) {
	log.Info().Msgf("fetch user info: provider=%s", provider)
	factoryErr := func(err error) (oauth.UserInfo, error) {
		log.Error().Stack().Err(err).Msg("failed to fetch user info")
		return oauth.UserInfo{}, err
	}

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return oauth.UserInfo{}, ErrInvalidProvider
	}

	// Exchange code to token
	log.Info().Msg("exchange code to token")
	socialLoginCallbackUrl := fmt.Sprintf("%s/auth/social/%s/callback/", s.cfg.Server.ExternalOrigin, provider)
	token, err := oauthProvider.ExchangeCodeToToken(ctx, input.Code, socialLoginCallbackUrl)
	if err != nil {
		return factoryErr(err)
	}

	// Get user info
	log.Info().Msg("get user info")
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
	log.Info().Msg("register or login")
	factoryErr := func(err error) (dto.JwtTokensOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to register or login")
		return dto.JwtTokensOutput{}, err
	}

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, userInfo.Email, true)
	if err != nil {
		return factoryErr(err)
	}

	// Save user
	if userEntity == nil {
		log.Info().Msg("create new user")
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
	log.Debug().Msg("search unique username")
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

		log.Debug().Msg("unique username not found")
	}
}
