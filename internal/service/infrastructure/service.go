package infrastructure

import (
	"context"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/pkg/model/v0"
)

var (
	// JWT error

	ErrInvalidJWTToken = v0.NewI18nError("invalid JWT token", "errors.session_invalid")
	ErrExpiredJWTToken = v0.NewI18nError("expired JWT token", "errors.session_expired")
	ErrBannedJWTToken  = v0.NewI18nError("banned JWT token", "errors.session_banned")

	// OTP error

	ErrInvalidOrExpiredOTP = v0.NewI18nError("invalid or expired otp", "errors.invalid_or_expired_otp")
)

type JWTService interface {
	GetTypeToken(ctx context.Context, token string) (string, error)
	GetAccessTokenClaims(ctx context.Context, token string) (AccessTokenClaims, error)
	GetRefreshTokenClaims(ctx context.Context, token string) (RefreshTokenClaims, error)
	BanToken(ctx context.Context, jti string) error
	GenerateTokens(ctx context.Context, userEntity *entity.User) (string, string, error)
}

type OTPService interface {
	GenerateCode(ctx context.Context) (string, error)
	SaveWithCode(ctx context.Context, prefix string, code string, data any) error
	GenerateAndSaveWithCode(ctx context.Context, prefix string, data any) (string, error)
	GetDataByCode(ctx context.Context, prefix string, code string, data any) error
	DeleteDataByCode(ctx context.Context, prefix string, code string) error
}
