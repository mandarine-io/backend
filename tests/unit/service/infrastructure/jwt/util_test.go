package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"time"
)

func createRefreshToken(t provider.T) (string, string) {
	jti := uuid.New().String()
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  "mandarine",
			"sub":  uuid.New().String(),
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
			"jti":  jti,
			"type": "refresh",
		},
	)

	refreshTokenSigned, err := refreshToken.SignedString([]byte(cfg.Secret))
	t.Require().NoError(err)

	return refreshTokenSigned, jti
}

func createAccessToken(t provider.T) (string, string) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            uuid.New().String(),
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Duration(cfg.AccessTokenTTL) * time.Second).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "username",
			"email":          "email",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)

	accessTokenSigned, err := accessToken.SignedString([]byte(cfg.Secret))
	t.Require().NoError(err)

	return accessTokenSigned, jti
}

func createTokenWithClaims(t provider.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenSigned, err := token.SignedString([]byte(cfg.Secret))
	t.Require().NoError(err)

	return tokenSigned
}

func createTokenWithIncorrectSigningMethod(t provider.T) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":  "mandarine",
			"sub":  uuid.New().String(),
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
			"jti":  uuid.New().String(),
			"type": "refresh",
		},
	)

	tokenSigned, err := token.SignedString([]byte(cfg.Secret))
	t.Require().NoError(err)

	return tokenSigned
}
