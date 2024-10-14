package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/repo"
	"mandarine/pkg/logging"
	"mandarine/pkg/rest/dto"
	"net/http"
	"strings"
	"time"
)

const (
	AuthUserKey = "authUser"
)

var (
	ErrJwtTokenIsMissing = dto.NewI18nError("JWT token is missing", "errors.session_invalid")
	ErrInvalidJwtToken   = dto.NewI18nError("invalid JWT token", "errors.session_invalid")
	ErrExpiredJwtToken   = dto.NewI18nError("expired JWT token", "errors.session_expired")
	ErrBannedJwtToken    = dto.NewI18nError("banned JWT token", "errors.session_banned")
	ErrUserNotFound      = dto.NewI18nError("user not found", "errors.user_not_found")
)

type AuthUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	JTI      string    `json:"jti"`
}

func JWTMiddleware(cfg config.JWTConfig, bannedTokenRepo repo.BannedTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		factoryErr := func(err error) {
			slog.Error("JWT authentication error", logging.ErrorAttr(err))
			_ = c.AbortWithError(http.StatusUnauthorized, err)
		}

		// Extract access token
		bearerHeader := c.Request.Header.Get("Authorization")
		if bearerHeader == "" {
			factoryErr(ErrJwtTokenIsMissing)
			return
		}

		accessToken, _ := strings.CutPrefix(bearerHeader, "Bearer ")
		token, err := security.DecodeAndValidateJwtToken(accessToken, cfg.Secret)
		if err != nil {
			factoryErr(err)
			return
		}

		// Check claims
		claims, err := security.GetClaimsFromJwtToken(token)
		if err != nil {
			factoryErr(err)
			return
		}
		sub, err := claims.GetSubject()
		if err != nil {
			factoryErr(ErrInvalidJwtToken)
			return
		}
		exp, err := claims.GetExpirationTime()
		if err != nil {
			factoryErr(ErrInvalidJwtToken)
			return
		}
		jti, ok := claims["jti"].(string)
		if !ok {
			factoryErr(ErrInvalidJwtToken)
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			factoryErr(ErrInvalidJwtToken)
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			factoryErr(ErrInvalidJwtToken)
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			factoryErr(ErrInvalidJwtToken)
			return
		}

		// Check if token has expired
		if exp.Unix() < time.Now().Unix() {
			factoryErr(ErrExpiredJwtToken)
			return
		}

		// Check if token is banned
		isBanned, err := bannedTokenRepo.ExistsBannedTokenByJTI(c, jti)
		if err != nil {
			factoryErr(err)
			return
		}
		if isBanned {
			factoryErr(ErrBannedJwtToken)
			return
		}

		// Save authenticated user
		userId, err := uuid.Parse(sub)
		if err != nil {
			factoryErr(ErrInvalidJwtToken)
			return
		}

		authUser := AuthUser{ID: userId, Username: username, Email: email, Role: role, JTI: jti}
		c.Set(AuthUserKey, authUser)
	}
}

func GetAuthUser(ctx *gin.Context) (AuthUser, error) {
	authUserAny, ok := ctx.Get(AuthUserKey)
	if !ok {
		slog.Error("Authenticated user getting error", logging.ErrorAttr(ErrUserNotFound))
		return AuthUser{}, ErrUserNotFound
	}

	authUser, ok := authUserAny.(AuthUser)
	if !ok {
		slog.Error("Authenticated user getting error", logging.ErrorAttr(ErrUserNotFound))
		return AuthUser{}, ErrUserNotFound
	}

	return authUser, nil
}
