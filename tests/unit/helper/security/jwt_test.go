package security_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mandarine/internal/api/config"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/model"
	"testing"
	"time"
)

func Test_SecurityUtil_DecodeAndValidateJwtToken(t *testing.T) {
	secret := "testsecret"
	validToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "1234567890",
			"iss": "mandarine",
			"exp": time.Now().Add(time.Hour).Unix(),
		},
	)
	validTokenString, _ := validToken.SignedString([]byte(secret))

	t.Run(
		"valid token", func(t *testing.T) {
			token, err := security.DecodeAndValidateJwtToken(validTokenString, secret)
			assert.NoError(t, err)
			assert.True(t, token.Valid)
		},
	)

	t.Run(
		"invalid token - unexpected signing method", func(t *testing.T) {
			key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

			invalidToken := jwt.NewWithClaims(
				jwt.SigningMethodES256, jwt.MapClaims{
					"sub": "1234567890",
					"iss": "mandarine",
				},
			)
			invalidTokenString, _ := invalidToken.SignedString(key)

			_, err := security.DecodeAndValidateJwtToken(invalidTokenString, secret)
			assert.ErrorIs(t, err, security.ErrInvalidJwtToken)
		},
	)

	t.Run(
		"invalid token - wrong signature", func(t *testing.T) {
			_, err := security.DecodeAndValidateJwtToken(validTokenString, "wrongsecret")
			assert.ErrorIs(t, err, security.ErrInvalidJwtToken)
		},
	)

	t.Run(
		"invalid token - wrong issuer", func(t *testing.T) {
			invalidToken := jwt.NewWithClaims(
				jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": "1234567890",
					"iss": "wrong issuer",
					"exp": time.Now().Add(time.Hour).Unix(),
				},
			)
			invalidTokenString, _ := invalidToken.SignedString([]byte(secret))

			_, err := security.DecodeAndValidateJwtToken(invalidTokenString, secret)
			assert.ErrorIs(t, err, security.ErrInvalidJwtToken)
		},
	)

	t.Run(
		"malformed token", func(t *testing.T) {
			_, err := security.DecodeAndValidateJwtToken("malformed.token.string", secret)
			assert.ErrorIs(t, err, security.ErrInvalidJwtToken)
		},
	)
}

func Test_SecurityUtil_GetClaimsFromJwtToken(t *testing.T) {
	claims := jwt.MapClaims{
		"sub": "1234567890",
		"iss": "mandarine",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t.Run(
		"valid claims", func(t *testing.T) {
			resultClaims, err := security.GetClaimsFromJwtToken(token)
			assert.NoError(t, err)
			assert.Equal(t, claims, resultClaims)
		},
	)
}

func Test_SecurityUtil_GenerateTokens(t *testing.T) {
	cfg := config.JWTConfig{
		Secret:          "testsecret",
		AccessTokenTTL:  3600,  // 1 hour
		RefreshTokenTTL: 86400, // 1 day
	}
	userEntity := &model.UserEntity{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "user@example.com",
		Role:     model.RoleEntity{Name: "user"},
	}

	t.Run(
		"generate valid tokens", func(t *testing.T) {
			accessToken, refreshToken, err := security.GenerateTokens(cfg, userEntity)
			assert.NoError(t, err)
			assert.NotEmpty(t, accessToken)
			assert.NotEmpty(t, refreshToken)

			// Validate the access token
			token, err := security.DecodeAndValidateJwtToken(accessToken, cfg.Secret)
			assert.NoError(t, err)
			assert.True(t, token.Valid)

			claims, ok := token.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			assert.Equal(t, userEntity.ID.String(), claims["sub"])
			assert.Equal(t, userEntity.Username, claims["username"])
			assert.Equal(t, userEntity.Email, claims["email"])
			assert.Equal(t, userEntity.Role.Name, claims["role"])
		},
	)
}
