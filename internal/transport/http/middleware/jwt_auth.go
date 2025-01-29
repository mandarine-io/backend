package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	AuthUserKey = "authUser"
)

var (
	ErrJWTTokenIsMissing = v0.NewI18nError("JWT token is missing", "errors.session_invalid")
	ErrUserNotFound      = v0.NewI18nError("user not found", "errors.user_not_found")
)

type AuthUser struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	IsPasswordTemp bool      `json:"isPasswordTemp"`
	IsEnabled      bool      `json:"isEnabled"`
	IsDeleted      bool      `json:"isDeleted"`
	JTI            string    `json:"jti"`
}

func JWTAuthMiddleware(jwtService infrastructure.JWTService) gin.HandlerFunc {
	log.Debug().Msg("setup JWT auth middleware")
	logger := log.With().Str("middleware", "jwt-auth").Logger()

	return func(c *gin.Context) {
		bearerHeader := c.Request.Header.Get("Authorization")
		if bearerHeader == "" {
			logger.Error().Stack().Msg("not found Authorization header")
			_ = c.AbortWithError(http.StatusUnauthorized, ErrJWTTokenIsMissing)
			return
		}

		accessToken, _ := strings.CutPrefix(bearerHeader, "Bearer ")

		claims, err := jwtService.GetAccessTokenClaims(c, accessToken)
		if err != nil {
			logger.Error().Err(err).Stack().Msg("failed to parse JWT token")
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		authUser := AuthUser{
			ID:             claims.UserID,
			Username:       claims.Username,
			Email:          claims.Email,
			Role:           claims.Role,
			IsPasswordTemp: claims.IsPasswordTemp,
			IsEnabled:      claims.IsEnabled,
			IsDeleted:      claims.IsDeleted,
			JTI:            claims.JTI,
		}

		c.Set(AuthUserKey, authUser)
	}
}

func GetAuthUser(ctx *gin.Context) (AuthUser, error) {
	authUserAny, ok := ctx.Get(AuthUserKey)
	if !ok {
		return AuthUser{}, ErrUserNotFound
	}

	authUser, ok := authUserAny.(AuthUser)
	if !ok {
		return AuthUser{}, ErrUserNotFound
	}

	return authUser, nil
}
