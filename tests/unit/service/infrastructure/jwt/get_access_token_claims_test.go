package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"time"
)

type GetAccessTokenClaimsSuite struct {
	suite.Suite
}

func (suite *GetAccessTokenClaimsSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("JWT service")
	t.Feature("GetAccessTokenClaimsSuite")
	t.Tags("Positive")

	accessToken, jti := createAccessToken(t)

	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Return(cache.ErrCacheEntryNotFound)

	claims, err := svc.GetAccessTokenClaims(ctx, accessToken)

	t.Require().NoError(err)
	t.Require().Equal("username", claims.Username)
	t.Require().Equal("email", claims.Email)
	t.Require().Equal("user", claims.Role)
	t.Require().Equal(false, claims.IsPasswordTemp)
	t.Require().Equal(true, claims.IsEnabled)
	t.Require().Equal(false, claims.IsDeleted)
	t.Require().Equal(jti, claims.JTI)
	t.Require().Less(time.Now().Unix(), claims.Exp)
}

func (suite *GetAccessTokenClaimsSuite) Test_ErrInvalidJWTToken(t provider.T) {
	t.Title("Returns invalid JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetAccessTokenClaimsSuite")
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
			name: "Incorrect token type",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "refresh",
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
					"type": "access",
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
					"type": "access",
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
					"type": "access",
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
					"type": "access",
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
					"type": "access",
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
					"type": "access",
				},
			),
		},
		{
			name:  "Invalid signing method",
			token: createTokenWithIncorrectSigningMethod(t),
		},
		{
			name: "Without username",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "access",
				},
			),
		},
		{
			name: "Without username",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":  "mandarine",
					"sub":  uuid.New().String(),
					"iat":  time.Now().Unix(),
					"exp":  time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":  uuid.New().String(),
					"type": "access",
				},
			),
		},
		{
			name: "Without email",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":      "mandarine",
					"sub":      uuid.New().String(),
					"iat":      time.Now().Unix(),
					"exp":      time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":      uuid.New().String(),
					"type":     "access",
					"username": "username",
				},
			),
		},
		{
			name: "Without role",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":      "mandarine",
					"sub":      uuid.New().String(),
					"iat":      time.Now().Unix(),
					"exp":      time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":      uuid.New().String(),
					"type":     "access",
					"username": "username",
					"email":    "email",
				},
			),
		},
		{
			name: "Without IsPasswordTemp",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":      "mandarine",
					"sub":      uuid.New().String(),
					"iat":      time.Now().Unix(),
					"exp":      time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":      uuid.New().String(),
					"type":     "access",
					"username": "username",
					"email":    "email",
					"role":     "user",
				},
			),
		},
		{
			name: "Incorrect IsPasswordTemp",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":            "mandarine",
					"sub":            uuid.New().String(),
					"iat":            time.Now().Unix(),
					"exp":            time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":            uuid.New().String(),
					"type":           "access",
					"username":       "username",
					"email":          "email",
					"role":           "user",
					"IsPasswordTemp": "incorrect",
				},
			),
		},
		{
			name: "Without isEnabled",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":            "mandarine",
					"sub":            uuid.New().String(),
					"iat":            time.Now().Unix(),
					"exp":            time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":            uuid.New().String(),
					"type":           "access",
					"username":       "username",
					"email":          "email",
					"role":           "user",
					"IsPasswordTemp": false,
				},
			),
		},
		{
			name: "Incorrect isEnabled",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":            "mandarine",
					"sub":            uuid.New().String(),
					"iat":            time.Now().Unix(),
					"exp":            time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":            uuid.New().String(),
					"type":           "access",
					"username":       "username",
					"email":          "email",
					"role":           "user",
					"IsPasswordTemp": false,
					"isEnabled":      "incorrect",
				},
			),
		},
		{
			name: "Without isDeleted",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":            "mandarine",
					"sub":            uuid.New().String(),
					"iat":            time.Now().Unix(),
					"exp":            time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":            uuid.New().String(),
					"type":           "access",
					"username":       "username",
					"email":          "email",
					"role":           "user",
					"IsPasswordTemp": false,
					"isEnabled":      true,
				},
			),
		},
		{
			name: "Incorrect isDeleted",
			token: createTokenWithClaims(
				t, jwt.MapClaims{
					"iss":            "mandarine",
					"sub":            uuid.New().String(),
					"iat":            time.Now().Unix(),
					"exp":            time.Now().Add(time.Duration(cfg.RefreshTokenTTL) * time.Second).Unix(),
					"jti":            uuid.New().String(),
					"type":           "access",
					"username":       "username",
					"email":          "email",
					"role":           "user",
					"IsPasswordTemp": false,
					"isEnabled":      true,
					"isDeleted":      "incorrect",
				},
			),
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t provider.T) {
				_, err := svc.GetAccessTokenClaims(ctx, tc.token)

				t.Require().Error(err)
				t.Require().Equal(infrastructure.ErrInvalidJWTToken, err)
			},
		)
	}
}

func (suite *GetAccessTokenClaimsSuite) Test_ErrExpiredJWTToken(t provider.T) {
	t.Title("Returns expired JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetAccessTokenClaims")
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

	_, err := svc.GetAccessTokenClaims(ctx, expiredToken)

	t.Require().Error(err)
	t.Require().Equal(infrastructure.ErrExpiredJWTToken, err)
}

func (suite *GetAccessTokenClaimsSuite) Test_ErrBannedJWTToken(t provider.T) {
	t.Title("Returns banned JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetAccessTokenClaimsSuite")
	t.Tags("Negative")

	accessToken, jti := createAccessToken(t)

	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Run(
		func(args mock.Arguments) {
			jtiPtr := args.Get(2).(*string)
			*jtiPtr = jti
		},
	).Return(nil)

	_, err := svc.GetAccessTokenClaims(ctx, accessToken)

	t.Require().Error(err)
	t.Require().ErrorIs(err, infrastructure.ErrBannedJWTToken)
}

func (suite *GetAccessTokenClaimsSuite) Test_ErrGetCache(t provider.T) {
	t.Title("Returns getting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetAccessTokenClaimsSuite")
	t.Tags("Negative")

	accessToken, _ := createAccessToken(t)

	cacheErr := errors.New("cache error")
	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Return(cacheErr)

	_, err := svc.GetAccessTokenClaims(ctx, accessToken)

	t.Require().Error(err)
	t.Require().ErrorIs(err, cacheErr)
}
