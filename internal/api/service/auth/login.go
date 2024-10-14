package auth

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/repo"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/pkg/logging"
	dto2 "mandarine/pkg/rest/dto"
)

var (
	ErrBadCredentials  = dto2.NewI18nError("bad credentials", "errors.bad_credentials")
	ErrUserIsBlocked   = dto2.NewI18nError("user is blocked", "errors.user_is_blocked")
	ErrUserNotFound    = dto2.NewI18nError("user not found", "errors.user_not_found")
	ErrInvalidJwtToken = dto2.NewI18nError("invalid JWT token", "errors.invalid_jwt_token")
)

type LoginService struct {
	userRepo repo.UserRepository
	cfg      *config.Config
}

func NewLoginService(userRepo repo.UserRepository, cfg *config.Config) *LoginService {
	return &LoginService{userRepo: userRepo, cfg: cfg}
}

//////////////////// Login ////////////////////

func (s *LoginService) Login(ctx context.Context, input dto.LoginInput) (dto.JwtTokensOutput, error) {
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

func (s *LoginService) RefreshTokens(ctx context.Context, refreshToken string) (dto.JwtTokensOutput, error) {
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
