package auth

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/mandarine-io/backend/third_party/oauth"
	"github.com/rs/zerolog"
)

const (
	registerCachePrefix       = "register"
	registerEmailDefaultTitle = "Activate your account"

	recoveryPasswordCachePrefix = "recovery_password"
	recoveryEmailDefaultTitle   = "Recovery password"
)

type svc struct {
	cfg            config.Config
	userRepo       repo.UserRepository
	oauthProviders map[string]oauth.Provider
	smtpSender     smtp.Sender
	templateEngine template.Engine
	jwtService     infra.JWTService
	otpService     infra.OTPService
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
	smtpSender smtp.Sender,
	templateEngine template.Engine,
	userRepo repo.UserRepository,
	jwtService infra.JWTService,
	otpService infra.OTPService,
	oauthProviders map[string]oauth.Provider,
	opts ...Option,
) domain.AuthService {
	s := &svc{
		userRepo:       userRepo,
		oauthProviders: oauthProviders,
		smtpSender:     smtpSender,
		templateEngine: templateEngine,
		jwtService:     jwtService,
		otpService:     otpService,
		cfg:            cfg,
		logger:         zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

//////////////////// Register ////////////////////

func (s *svc) Register(ctx context.Context, input v0.RegisterInput, localizer locale.Localizer) error {
	s.logger.Info().Msgf("register: %s", input.Username)

	// Check if user exists
	exists, err := s.userRepo.ExistsUserByUsernameOrEmail(ctx, input.Username, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check exist user")
		return err
	}
	if exists {
		s.logger.Error().Stack().Err(domain.ErrDuplicateUser).Msg("user already exists")
		return domain.ErrDuplicateUser
	}

	// Hashing password
	input.Password, err = security.HashPassword(input.Password)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to hash password")
		return err
	}

	// Create and save OTP
	otp, err := s.otpService.GenerateAndSaveWithCode(ctx, registerCachePrefix, input)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to create and save OTP")
		return err
	}

	// Localize email title
	emailTitle := registerEmailDefaultTitle
	if localizer != nil {
		emailTitle = localizer.Localize("email.register-confirm.title", nil, 0)
	}

	// Send mail
	args := v0.RegisterConfirmTemplateArgs{
		Email: input.Email,
		TTL:   s.cfg.Security.OTP.TTL / 60,
		OTP:   otp,
	}
	content, err := s.templateEngine.RenderHTML("register-confirm", args)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to register user")
		return err
	}

	err = s.smtpSender.SendHTMLMessage(emailTitle, content, s.cfg.SMTP.From, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to register user")
		return domain.ErrSendEmail
	}

	return nil
}

//////////////////// Register confirmation ////////////////////

func (s *svc) RegisterConfirm(ctx context.Context, input v0.RegisterConfirmInput) error {
	s.logger.Info().Msg("register confirm")

	// Get data by OTP
	var registerInput v0.RegisterInput
	err := s.otpService.GetDataByCode(ctx, registerCachePrefix, input.OTP, &registerInput)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get data by OTP")
		return err
	}

	// Check email
	if input.Email != registerInput.Email {
		s.logger.Error().Stack().Err(infra.ErrInvalidOrExpiredOTP).Msg("dont match emails in request and OTP data")
		return infra.ErrInvalidOrExpiredOTP
	}

	// Check if user exists
	exists, err := s.userRepo.ExistsUserByUsernameOrEmail(ctx, registerInput.Username, registerInput.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check exist user")
		return err
	}
	if exists {
		s.logger.Error().Stack().Err(domain.ErrDuplicateUser).Msg("user already exists")
		return domain.ErrDuplicateUser
	}

	// Create user in DB
	user := converter.MapRegisterInputToUserEntity(registerInput)

	_, err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to create user")

		if errors.Is(err, repo.ErrDuplicateUser) {
			return domain.ErrDuplicateUser
		}
		return err
	}

	// Delete cache
	err = s.otpService.DeleteDataByCode(ctx, registerCachePrefix, input.OTP)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to delete OTP data")
	}

	return nil
}

//////////////////// Login ////////////////////

func (s *svc) Login(ctx context.Context, input v0.LoginInput) (v0.JwtTokensOutput, error) {
	s.logger.Info().Msg("login")

	// Get user entity
	user, err := s.userRepo.FindUserByUsernameOrEmail(ctx, input.Login, s.userRepo.WithRolePreload())
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.JwtTokensOutput{}, err
	}
	if user == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return v0.JwtTokensOutput{}, domain.ErrUserNotFound
	}

	// Check password
	if !security.CheckPasswordHash(input.Password, user.Password) {
		s.logger.Error().Stack().Err(domain.ErrBadCredentials).Msg("failed to check hash password")
		return v0.JwtTokensOutput{}, domain.ErrBadCredentials
	}

	// Check if user is blocked
	if !user.IsEnabled {
		s.logger.Error().Stack().Err(domain.ErrUserIsBlocked).Msg("user is blocked")
		return v0.JwtTokensOutput{}, domain.ErrUserIsBlocked
	}

	// Create JWT tokens
	accessToken, refreshToken, err := s.jwtService.GenerateTokens(ctx, user)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to generate JWT tokens")
		return v0.JwtTokensOutput{}, err
	}

	return v0.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Refresh Tokens ////////////////////

func (s *svc) RefreshTokens(ctx context.Context, input v0.RefreshTokensInput) (v0.JwtTokensOutput, error) {
	s.logger.Info().Msg("refresh tokens")

	claims, err := s.jwtService.GetRefreshTokenClaims(ctx, input.RefreshToken)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to getting refresh token claims")
		return v0.JwtTokensOutput{}, err
	}

	// Get user entity
	user, err := s.userRepo.FindUserByID(ctx, claims.UserID, s.userRepo.WithRolePreload())
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.JwtTokensOutput{}, err
	}
	if user == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return v0.JwtTokensOutput{}, domain.ErrUserNotFound
	}

	// Check if user is blocked
	if !user.IsEnabled {
		s.logger.Error().Stack().Err(domain.ErrUserIsBlocked).Msg("user is blocked")
		return v0.JwtTokensOutput{}, domain.ErrUserIsBlocked
	}

	// Create JWT tokens
	accessToken, refreshToken, err := s.jwtService.GenerateTokens(ctx, user)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to generate JWT tokens")
		return v0.JwtTokensOutput{}, err
	}

	return v0.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Logout ////////////////////

func (s *svc) Logout(ctx context.Context, jti string) error {
	s.logger.Info().Msg("logout")

	err := s.jwtService.BanToken(ctx, jti)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to ban token")
		return err
	}

	return nil
}

//////////////////// Recovery password ////////////////////

func (s *svc) RecoveryPassword(
	ctx context.Context,
	input v0.RecoveryPasswordInput,
	localizer locale.Localizer,
) error {
	s.logger.Info().Msg("recovery password")

	// Get user by email
	exists, err := s.userRepo.ExistsUserByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check exist user")
		return err
	}
	if !exists {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return domain.ErrUserNotFound
	}

	// Create and save OTP
	otp, err := s.otpService.GenerateAndSaveWithCode(ctx, recoveryPasswordCachePrefix, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to create and save OTP code")
		return err
	}

	// Localize email title
	emailTitle := recoveryEmailDefaultTitle
	if localizer != nil {
		emailTitle = localizer.Localize("email.recovery-password.title", nil, 0)
	}

	// Send email
	args := v0.RegisterConfirmTemplateArgs{
		Email: input.Email,
		TTL:   s.cfg.Security.OTP.TTL / 60,
		OTP:   otp,
	}
	content, err := s.templateEngine.RenderHTML("recovery-password", args)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to render HTML email template")
		return err
	}

	err = s.smtpSender.SendHTMLMessage(emailTitle, content, s.cfg.SMTP.From, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Err(err).Msg("failed to send HTML message")
		return domain.ErrSendEmail
	}

	return nil
}

//////////////////// Verify recovery password ////////////////////

func (s *svc) VerifyRecoveryCode(ctx context.Context, input v0.VerifyRecoveryCodeInput) error {
	s.logger.Info().Msg("verify recovery password")

	// Get data by OTP
	var email string
	err := s.otpService.GetDataByCode(ctx, recoveryPasswordCachePrefix, input.OTP, &email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get data by OTP")
		return err
	}

	// Check OTP
	if input.Email != email {
		s.logger.Error().Stack().Err(infra.ErrInvalidOrExpiredOTP).Msg("dont match emails in request and OTP data")
		return infra.ErrInvalidOrExpiredOTP
	}

	return nil
}

//////////////////// Reset password ////////////////////

func (s *svc) ResetPassword(ctx context.Context, input v0.ResetPasswordInput) error {
	s.logger.Info().Msg("reset password")

	// Get entry from cache
	var email string
	err := s.otpService.GetDataByCode(ctx, recoveryPasswordCachePrefix, input.OTP, &email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get data by OTP")
		return err
	}

	// Check OTP
	if input.Email != email {
		s.logger.Error().Stack().Err(infra.ErrInvalidOrExpiredOTP).Msg("dont match emails in request and OTP data")
		return infra.ErrInvalidOrExpiredOTP
	}

	// Get user by email
	user, err := s.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return err
	}
	if user == nil {
		s.logger.Error().Stack().Err(domain.ErrUserNotFound).Msg("user not found")
		return domain.ErrUserNotFound
	}

	// Hashing password
	user.Password, err = security.HashPassword(input.Password)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to hash password")
		return err
	}

	// Update password
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to update user")
		return err
	}

	// Delete cache
	err = s.otpService.DeleteDataByCode(ctx, recoveryPasswordCachePrefix, input.OTP)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to delete OTP data")
	}

	return nil
}

//////////////////// Get consent page url ////////////////////

func (s *svc) GetConsentPageURL(_ context.Context, provider string, redirectURL string) (
	v0.GetConsentPageURLOutput,
	error,
) {
	s.logger.Info().Msgf("get consent page url: provider=%s", provider)

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		s.logger.Error().Stack().Err(domain.ErrInvalidProvider).Msg("provider not found")
		return v0.GetConsentPageURLOutput{}, domain.ErrInvalidProvider
	}

	consentPageURL, oauthState := oauthProvider.GetConsentPageURL(redirectURL)
	return v0.GetConsentPageURLOutput{ConsentPageURL: consentPageURL, OauthState: oauthState}, nil
}

//////////////////// Fetch user info ////////////////////

func (s *svc) FetchUserInfo(ctx context.Context, provider string, input v0.FetchUserInfoInput) (
	oauth.UserInfo, error,
) {
	s.logger.Info().Msgf("fetch user info: provider=%s", provider)

	// Get oauth provider
	oauthProvider, ok := s.oauthProviders[provider]
	if !ok {
		return oauth.UserInfo{}, domain.ErrInvalidProvider
	}

	// Exchange code to token
	s.logger.Info().Msg("exchange code to token")
	socialLoginCallbackURL := fmt.Sprintf("%s/auth/social/%s/callback/", s.cfg.Server.ExternalURL, provider)
	token, err := oauthProvider.ExchangeCodeToToken(ctx, input.Code, socialLoginCallbackURL)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to exchange code to token")
		return oauth.UserInfo{}, err
	}

	// Get user info
	s.logger.Info().Msg("get user info")
	userInfo, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to get user info")

		if errors.Is(err, oauth.ErrUserInfoNotReceived) {
			return oauth.UserInfo{}, domain.ErrUserInfoNotReceived
		}
		return oauth.UserInfo{}, err
	}

	return userInfo, nil
}

//////////////////// Register or login ////////////////////

func (s *svc) RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (v0.JwtTokensOutput, error) {
	s.logger.Info().Msg("register or login")

	// Get user by email
	user, err := s.userRepo.FindUserByEmail(ctx, userInfo.Email, s.userRepo.WithRolePreload())
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to find user")
		return v0.JwtTokensOutput{}, err
	}

	// Save user
	if user == nil {
		s.logger.Info().Msg("create new user")
		userInfo.Username, err = s.searchUniqueUsername(ctx, userInfo.Username)
		if err != nil {
			return v0.JwtTokensOutput{}, err
		}

		user = converter.MapUserInfoToUserEntity(userInfo)
		user, err = s.userRepo.CreateUser(ctx, user)
		if err != nil {
			s.logger.Error().Stack().Err(err).Msg("failed to create user")
			return v0.JwtTokensOutput{}, err
		}
	} else if !user.IsEnabled {
		s.logger.Error().Stack().Err(domain.ErrUserIsBlocked).Msg("user is banned")
		return v0.JwtTokensOutput{}, domain.ErrUserIsBlocked
	}

	// Create JWT tokens
	accessToken, refreshToken, err := s.jwtService.GenerateTokens(ctx, user)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to generate tokens")
		return v0.JwtTokensOutput{}, err
	}

	return v0.JwtTokensOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

//////////////////// Additional helpers ////////////////////

func (s *svc) searchUniqueUsername(ctx context.Context, defaultUsername string) (string, error) {
	s.logger.Debug().Msg("search unique username")

	exists, err := s.userRepo.ExistsUserByUsername(ctx, defaultUsername)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to check exist user")
		return "", err
	}
	if !exists {
		s.logger.Debug().Msgf("searched unique username: %s", defaultUsername)
		return defaultUsername, nil
	}

	s.logger.Debug().Msgf("username %s already exists", defaultUsername)

	for {
		username := fmt.Sprintf("%s_%s", defaultUsername, uuid.New().String())
		exists, err := s.userRepo.ExistsUserByUsername(ctx, username)
		if err != nil {
			s.logger.Error().Stack().Err(err).Msg("failed to check exist user")
			return "", err
		}
		if !exists {
			s.logger.Debug().Msgf("searched unique username: %s", username)
			return username, nil
		}

		s.logger.Debug().Msgf("username %s already exists", username)
	}
}
