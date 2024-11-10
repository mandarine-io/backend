package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/auth/mapper"
	cachehelper "github.com/mandarine-io/Backend/internal/helper/cache"
	"github.com/mandarine-io/Backend/internal/helper/random"
	"github.com/mandarine-io/Backend/internal/helper/security"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/mandarine-io/Backend/pkg/oauth"
	"github.com/mandarine-io/Backend/pkg/smtp"
	"github.com/mandarine-io/Backend/pkg/storage/cache"
	"github.com/mandarine-io/Backend/pkg/template"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	registerCachePrefix       = "register"
	registerEmailDefaultTitle = "Activate your account"

	recoveryPasswordCachePrefix = "recovery_password"
	recoveryEmailDefaultTitle   = "Recovery password"
)

type svc struct {
	userRepo        repo.UserRepository
	bannedTokenRepo repo.BannedTokenRepository
	oauthProviders  map[string]oauth.Provider
	cacheManager    cache.Manager
	smtpSender      smtp.Sender
	templateEngine  template.Engine
	cfg             *config.Config
}

func NewService(
	userRepo repo.UserRepository,
	bannedTokenRepo repo.BannedTokenRepository,
	oauthProviders map[string]oauth.Provider,
	cacheManager cache.Manager,
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	cfg *config.Config,
) service.AuthService {
	return &svc{
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

func (s *svc) Register(ctx context.Context, input dto.RegisterInput, localizer *i18n.Localizer) error {
	log.Info().Msgf("register: %s", input.Username)

	// Check if user exists
	exists, err := s.userRepo.ExistsUserByUsernameOrEmail(ctx, input.Username, input.Email)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register user")
		return err
	}
	if exists {
		log.Error().Stack().Err(service.ErrDuplicateUser).Msg("failed to register user")
		return service.ErrDuplicateUser
	}

	// Hashing password
	input.Password, err = security.HashPassword(input.Password)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register user")
		return err
	}

	// Generate OTP
	otp, err := random.GenerateRandomNumber(s.cfg.Security.OTP.Length)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register user")
		return err
	}

	// Save in cache
	expiration := time.Duration(s.cfg.Security.OTP.TTL) * time.Second
	cacheEntry := dto.RegisterCache{
		User:      input,
		OTP:       otp,
		ExpiredAt: time.Now().Add(expiration),
	}
	err = s.cacheManager.SetWithExpiration(ctx, cachehelper.CreateCacheKey(registerCachePrefix, input.Email), cacheEntry, expiration)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register user")
		return err
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
		log.Error().Stack().Err(err).Msg("failed to register user")
		return err
	}

	err = s.smtpSender.SendHtmlMessage(emailTitle, content, input.Email)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register user")
		return service.ErrSendEmail
	}

	return nil
}

//////////////////// Register confirmation ////////////////////

func (s *svc) RegisterConfirm(ctx context.Context, input dto.RegisterConfirmInput) error {
	log.Info().Msg("register confirm")

	// Get entry from cache
	var cacheEntry dto.RegisterCache
	err := s.cacheManager.Get(ctx, cachehelper.CreateCacheKey(registerCachePrefix, input.Email), &cacheEntry)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register confirm user")

		if errors.Is(err, cache.ErrCacheEntryNotFound) {
			return service.ErrInvalidOrExpiredOtp
		}
		return err
	}

	// Check OTP
	if input.OTP != cacheEntry.OTP {
		log.Error().Stack().Err(service.ErrInvalidOrExpiredOtp).Msg("failed to register confirm user")
		return service.ErrInvalidOrExpiredOtp
	}

	// Check email
	if input.Email != cacheEntry.User.Email {
		log.Error().Stack().Err(service.ErrInvalidOrExpiredOtp).Msg("failed to register confirm user")
		return service.ErrInvalidOrExpiredOtp
	}

	// Check if user exists
	exists, err := s.userRepo.ExistsUserByUsernameOrEmail(ctx, cacheEntry.User.Username, cacheEntry.User.Email)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register confirm user")
		return err
	}
	if exists {
		log.Error().Stack().Err(service.ErrDuplicateUser).Msg("failed to register confirm user")
		return service.ErrDuplicateUser
	}

	// Create user in DB
	registerRequest := cacheEntry.User
	userEntity := mapper.MapRegisterRequestToUserEntity(registerRequest)
	_, err = s.userRepo.CreateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register confirm user")

		if errors.Is(err, repo.ErrDuplicateUser) {
			return service.ErrDuplicateUser
		}
		return err
	}

	// Delete cache
	err = s.cacheManager.Delete(ctx, cachehelper.CreateCacheKey(registerCachePrefix, input.Email))
	if err != nil {
		log.Warn().Err(err).Msg("failed to delete cache")
	}

	return nil
}

//////////////////// Login ////////////////////

func (s *svc) Login(ctx context.Context, input dto.LoginInput) (dto.JwtTokensOutput, error) {
	log.Info().Msg("login")

	// Get user entity
	userEntity, err := s.userRepo.FindUserByUsernameOrEmail(ctx, input.Login, true)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to login")
		return dto.JwtTokensOutput{}, err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to login")
		return dto.JwtTokensOutput{}, service.ErrUserNotFound
	}

	// Check password
	if !security.CheckPasswordHash(input.Password, userEntity.Password) {
		log.Error().Stack().Err(service.ErrBadCredentials).Msg("failed to login")
		return dto.JwtTokensOutput{}, service.ErrBadCredentials
	}

	// Check if user is blocked
	if !userEntity.IsEnabled {
		log.Error().Stack().Err(service.ErrUserIsBlocked).Msg("failed to login")
		return dto.JwtTokensOutput{}, service.ErrUserIsBlocked
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to login")
		return dto.JwtTokensOutput{}, err
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Refresh Tokens ////////////////////

func (s *svc) RefreshTokens(ctx context.Context, refreshToken string) (dto.JwtTokensOutput, error) {
	log.Info().Msg("refresh tokens")

	// Check token
	token, err := security.DecodeAndValidateJwtToken(refreshToken, s.cfg.Security.JWT.Secret)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, service.ErrInvalidJwtToken
	}

	// Get user id from token
	claims, err := security.GetClaimsFromJwtToken(token)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, service.ErrInvalidJwtToken
	}

	sub, err := claims.GetSubject()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, service.ErrInvalidJwtToken
	}

	userUUID, err := uuid.Parse(sub)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, service.ErrInvalidJwtToken
	}

	// Get user entity
	userEntity, err := s.userRepo.FindUserById(ctx, userUUID, true)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, service.ErrUserNotFound
	}

	// Check if user is blocked
	if !userEntity.IsEnabled {
		log.Error().Stack().Err(service.ErrUserIsBlocked).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, service.ErrUserIsBlocked
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to refresh tokens")
		return dto.JwtTokensOutput{}, err
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Logout ////////////////////

func (s *svc) Logout(ctx context.Context, jti string) error {
	log.Info().Msg("logout")

	// Create banned token
	bannedToken := &model.BannedTokenEntity{
		JTI:       jti,
		ExpiredAt: time.Now().Add(time.Duration(s.cfg.Security.JWT.RefreshTokenTTL) * time.Second).Unix(),
	}
	_, err := s.bannedTokenRepo.CreateOrUpdateBannedToken(ctx, bannedToken)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to logout")
		return err
	}

	return nil
}

//////////////////// Recovery password ////////////////////

func (s *svc) RecoveryPassword(ctx context.Context, input dto.RecoveryPasswordInput, localizer *i18n.Localizer) error {
	log.Info().Msg("recovery password")

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, input.Email, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to recovery password")
		return err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to recovery password")
		return service.ErrUserNotFound
	}

	// Generate OTP
	otp, err := random.GenerateRandomNumber(s.cfg.Security.OTP.Length)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to recovery password")
		return err
	}

	// Save in cache
	expiration := time.Duration(s.cfg.Security.OTP.TTL) * time.Second
	cacheEntry := dto.RecoveryPasswordCache{
		Email:     input.Email,
		OTP:       otp,
		ExpiredAt: time.Now().Add(expiration),
	}

	err = s.cacheManager.SetWithExpiration(ctx, cachehelper.CreateCacheKey(recoveryPasswordCachePrefix, input.Email), cacheEntry, expiration)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to recovery password")
		return err
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
		log.Error().Stack().Err(err).Msg("failed to recovery password")
		return err
	}

	err = s.smtpSender.SendHtmlMessage(emailTitle, content, input.Email)
	if err != nil {
		log.Error().Stack().Err(err).Err(err).Msg("failed to recovery password")
		return service.ErrSendEmail
	}

	return nil
}

//////////////////// Verify recovery password ////////////////////

func (s *svc) VerifyRecoveryCode(ctx context.Context, input dto.VerifyRecoveryCodeInput) error {
	log.Info().Msg("verify recovery password")

	// Get entry from cache
	var cacheEntry dto.RecoveryPasswordCache
	err := s.cacheManager.Get(ctx, cachehelper.CreateCacheKey(recoveryPasswordCachePrefix, input.Email), &cacheEntry)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to verify recovery password")

		if errors.Is(err, cache.ErrCacheEntryNotFound) {
			return service.ErrInvalidOrExpiredOtp
		}
		return err
	}

	// Check OTP
	if input.OTP != cacheEntry.OTP {
		log.Error().Stack().Err(service.ErrInvalidOrExpiredOtp).Msg("failed to verify recovery password")
		return service.ErrInvalidOrExpiredOtp
	}

	return nil
}

//////////////////// Reset password ////////////////////

func (s *svc) ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error {
	log.Info().Msg("reset password")

	// Get entry from cache
	var cacheEntry dto.RecoveryPasswordCache
	err := s.cacheManager.Get(ctx, cachehelper.CreateCacheKey(recoveryPasswordCachePrefix, input.Email), &cacheEntry)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to reset password")

		if errors.Is(err, cache.ErrCacheEntryNotFound) {
			return service.ErrInvalidOrExpiredOtp
		}
		return err
	}

	// Check OTP
	if input.OTP != cacheEntry.OTP {
		log.Error().Stack().Err(service.ErrInvalidOrExpiredOtp).Msg("failed to reset password")
		return service.ErrInvalidOrExpiredOtp
	}

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, input.Email, false)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to reset password")
		return err
	}
	if userEntity == nil {
		log.Error().Stack().Err(service.ErrUserNotFound).Msg("failed to reset password")
		return service.ErrUserNotFound
	}

	// Hashing password
	userEntity.Password, err = security.HashPassword(input.Password)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to reset password")
		return err
	}

	// Update password
	_, err = s.userRepo.UpdateUser(ctx, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to reset password")
		return err
	}

	return nil
}

//////////////////// Get consent page url ////////////////////

func (s *svc) GetConsentPageUrl(_ context.Context, provider string, redirectUrl string) (dto.GetConsentPageUrlOutput, error) {
	log.Info().Msgf("get consent page url: provider=%s", provider)

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		log.Error().Stack().Err(service.ErrInvalidProvider).Msg("failed to get consent page url")
		return dto.GetConsentPageUrlOutput{}, service.ErrInvalidProvider
	}

	consentPageUrl, oauthState := oauthProvider.GetConsentPageUrl(redirectUrl)
	return dto.GetConsentPageUrlOutput{ConsentPageUrl: consentPageUrl, OauthState: oauthState}, nil
}

//////////////////// Fetch user info ////////////////////

func (s *svc) FetchUserInfo(ctx context.Context, provider string, input dto.FetchUserInfoInput) (
	oauth.UserInfo, error,
) {
	log.Info().Msgf("fetch user info: provider=%s", provider)

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return oauth.UserInfo{}, service.ErrInvalidProvider
	}

	// Exchange code to token
	log.Info().Msg("exchange code to token")
	socialLoginCallbackUrl := fmt.Sprintf("%s/auth/social/%s/callback/", s.cfg.Server.ExternalOrigin, provider)
	token, err := oauthProvider.ExchangeCodeToToken(ctx, input.Code, socialLoginCallbackUrl)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to fetch user info")
		return oauth.UserInfo{}, err
	}

	// Get user info
	log.Info().Msg("get user info")
	userInfo, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to fetch user info")

		if errors.Is(err, oauth.ErrUserInfoNotReceived) {
			return oauth.UserInfo{}, service.ErrUserInfoNotReceived
		}
		return oauth.UserInfo{}, err
	}

	return userInfo, nil
}

//////////////////// Register or login ////////////////////

func (s *svc) RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (dto.JwtTokensOutput, error) {
	log.Info().Msg("register or login")

	// Get user by email
	userEntity, err := s.userRepo.FindUserByEmail(ctx, userInfo.Email, true)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register or login")
		return dto.JwtTokensOutput{}, err
	}

	// Save user
	if userEntity == nil {
		log.Info().Msg("create new user")
		userInfo.Username, err = s.searchUniqueUsername(ctx, userInfo.Username)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to register or login")
			return dto.JwtTokensOutput{}, err
		}

		userEntity = mapper.MapUserInfoToUserEntity(userInfo)
		userEntity, err = s.userRepo.CreateUser(ctx, userEntity)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to register or login")
			return dto.JwtTokensOutput{}, err
		}
	} else if !userEntity.IsEnabled {
		log.Error().Stack().Err(service.ErrUserIsBlocked).Msg("failed to register or login")
		return dto.JwtTokensOutput{}, service.ErrUserIsBlocked
	}

	// Create JWT tokens
	accessToken, refreshToken, err := security.GenerateTokens(s.cfg.Security.JWT, userEntity)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register or login")
		return dto.JwtTokensOutput{}, err
	}

	return dto.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Additional helpers ////////////////////

func (s *svc) searchUniqueUsername(ctx context.Context, defaultUsername string) (string, error) {
	log.Debug().Msg("search unique username")
	exists, err := s.userRepo.ExistsUserByUsername(ctx, defaultUsername)
	if err != nil {
		return "", err
	}
	if !exists {
		return defaultUsername, nil
	}

	log.Debug().Msgf("username %s already exists", defaultUsername)

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

		log.Debug().Msgf("username %s already exists", username)
	}
}
