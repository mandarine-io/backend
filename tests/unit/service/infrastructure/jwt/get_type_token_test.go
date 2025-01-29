package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type GetTypeTokenSuite struct {
	suite.Suite
}

func (suite *GetTypeTokenSuite) Test_Success(t provider.T) {
	t.Title("Returns right token type")
	t.Severity(allure.NORMAL)
	t.Epic("JWT service")
	t.Feature("GetTypeToken")
	t.Tags("Positive")

	t.Run(
		"Refresh token", func(t provider.T) {
			refreshToken, _ := createRefreshToken(t)

			tokenType, err := svc.GetTypeToken(ctx, refreshToken)

			t.Require().NoError(err)
			t.Require().Equal("refresh", tokenType)
		},
	)

	t.Run(
		"Access token", func(t provider.T) {
			accessToken, _ := createAccessToken(t)

			tokenType, err := svc.GetTypeToken(ctx, accessToken)

			t.Require().NoError(err)
			t.Require().Equal("access", tokenType)
		},
	)
}

func (suite *GetTypeTokenSuite) Test_ErrInvalidJWTToken(t provider.T) {
	t.Title("Returns invalid JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetTypeToken")
	t.Tags("Negative")

	type testCase struct {
		name  string
		token string
	}

	testCases := []testCase{
		{
			name:  "Not JWT token",
			token: "invalid_token",
		},
		{
			name:  "Empty token",
			token: "",
		},
		{
			name: "Without token type",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss": "mandarine",
					"sub": uuid.New().String(),
					"iat": time.Now().Unix(),
					"exp": time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti": uuid.New().String(),
				},
			),
		},
		{
			name: "With incorrect issuer",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine1",
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "refresh",
				},
			),
		},
		{
			name: "Without issuer",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "refresh",
				},
			),
		},
		{
			name: "Without issued at",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"sub":  uuid.New().String(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "refresh",
				},
			),
		},
		{
			name: "Without expiration",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"jti":  uuid.New().String(),
					"type": "refresh",
				},
			),
		},
		{
			name: "Without subject",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "refresh",
				},
			),
		},
		{
			name: "Without JTI",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"type": "refresh",
				},
			),
		},
		{
			name:  "Invalid signing method",
			token: createTokenWithIncorrectSigningMethod(t),
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t provider.T) {
				_, err := svc.GetTypeToken(ctx, tc.token)

				t.Require().Error(err)
				t.Require().Equal(infrastructure.ErrInvalidJWTToken, err)
			},
		)
	}
}

func (suite *GetTypeTokenSuite) Test_ErrExpiredJWTToken(t provider.T) {
	t.Title("Returns expired JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetTypeToken")
	t.Tags("Negative")

	expiredToken := createTokenWithClaims(
		t, jwt.MapClaims{
			"iss":  "mandarine",
			"sub":  uuid.New().String(),
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(-1 * time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
			"jti":  uuid.New().String(),
			"type": "refresh",
		},
	)

	_, err := svc.GetTypeToken(ctx, expiredToken)

	t.Require().Error(err)
	t.Require().Equal(infrastructure.ErrExpiredJWTToken, err)
}
