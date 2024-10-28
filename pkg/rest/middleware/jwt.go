package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/internal/api/helper/security"
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	"github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/rs/zerolog/log"
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
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsEnabled bool      `json:"isEnabled"`
	IsDeleted bool      `json:"isDeleted"`
	JTI       string    `json:"jti"`
}

func JWTMiddleware(cfg config.JWTConfig, bannedTokenRepo repo.BannedTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		factoryErr := func(err error) {
			log.Error().Stack().Err(err).Msg("failed to authenticate user")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
		}
		factoryChildErr := func(err error, childErr error) {
			log.Error().Stack().Err(childErr).Msg("failed to authenticate user")
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
			factoryChildErr(ErrInvalidJwtToken, err)
			return
		}

		// Check claims
		claims, err := security.GetClaimsFromJwtToken(token)
		if err != nil {
			factoryChildErr(ErrInvalidJwtToken, err)
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
		isEnabled, ok := claims["isEnabled"].(bool)
		if !ok {
			factoryErr(ErrInvalidJwtToken)
			return
		}
		isDeleted, ok := claims["isDeleted"].(bool)
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

		authUser := AuthUser{
			ID:        userId,
			Username:  username,
			Email:     email,
			Role:      role,
			IsEnabled: isEnabled,
			IsDeleted: isDeleted,
			JTI:       jti}
		c.Set(AuthUserKey, authUser)
	}
}

func GetAuthUser(ctx *gin.Context) (AuthUser, error) {
	log.Debug().Msg("get authenticated user from gin context")
	authUserAny, ok := ctx.Get(AuthUserKey)
	if !ok {
		return AuthUser{}, ErrUserNotFound
	}

	log.Debug().Msgf("check authenticated user type: %T", authUserAny)
	authUser, ok := authUserAny.(AuthUser)
	if !ok {
		return AuthUser{}, ErrUserNotFound
	}

	return authUser, nil
}
