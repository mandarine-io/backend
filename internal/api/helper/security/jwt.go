package security

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/internal/api/persistence/model"
	"time"
)

const (
	jwtIssuer = "mandarine"
)

var (
	ErrInvalidJwtToken = errors.New("invalid JWT token")
)

func DecodeAndValidateJwtToken(token, secret string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(
		token, getJwtToken(secret), jwt.WithIssuer(jwtIssuer), jwt.WithIssuedAt(), jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(), jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, ErrInvalidJwtToken
	}
	return jwtToken, nil
}

func GetClaimsFromJwtToken(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidJwtToken
	}

	return claims, nil
}

func GenerateTokens(cfg config.JWTConfig, userEntity *model.UserEntity) (string, string, error) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":       jwtIssuer,
			"sub":       userEntity.ID.String(),
			"iat":       time.Now().Unix(),
			"exp":       time.Now().Add(time.Duration(cfg.AccessTokenTTL) * time.Second).Unix(),
			"jti":       jti,
			"username":  userEntity.Username,
			"email":     userEntity.Email,
			"role":      userEntity.Role.Name,
			"isEnabled": userEntity.IsEnabled,
			"isDeleted": userEntity.DeletedAt != nil,
		},
	)
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": jwtIssuer,
			"sub": userEntity.ID.String(),
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
			"jti": jti,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	refreshTokenSigned, err := refreshToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenSigned, refreshTokenSigned, nil
}

func getJwtToken(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJwtToken
		}
		return []byte(secret), nil
	}
}
