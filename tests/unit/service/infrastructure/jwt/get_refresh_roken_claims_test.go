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

type GetRefreshTokenClaimsSuite struct {
	suite.Suite
}

func (suite *GetRefreshTokenClaimsSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("JWT service")
	t.Feature("GetRefreshTokenClaimsSuite")
	t.Tags("Positive")

	refreshToken, jti := createRefreshToken(t)

	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Return(cache.ErrCacheEntryNotFound)

	claims, err := svc.GetRefreshTokenClaims(ctx, refreshToken)

	t.Require().NoError(err)
	t.Require().Equal(jti, claims.JTI)
	t.Require().Less(time.Now().Unix(), claims.Exp)
}

func (suite *GetRefreshTokenClaimsSuite) Test_ErrInvalidJWTToken(t provider.T) {
	t.Title("Returns invalid JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetRefreshTokenClaimsSuite")
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
					"type": "access",
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
				_, err := svc.GetRefreshTokenClaims(ctx, tc.token)

				t.Require().Error(err)
				t.Require().Equal(infrastructure.ErrInvalidJWTToken, err)
			},
		)
	}
}

func (suite *GetRefreshTokenClaimsSuite) Test_ErrExpiredJWTToken(t provider.T) {
	t.Title("Returns expired JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetRefreshTokenClaims")
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

	_, err := svc.GetRefreshTokenClaims(ctx, expiredToken)

	t.Require().Error(err)
	t.Require().Equal(infrastructure.ErrExpiredJWTToken, err)
}

func (suite *GetRefreshTokenClaimsSuite) Test_ErrBannedJWTToken(t provider.T) {
	t.Title("Returns banned JWT token error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetRefreshTokenClaimsSuite")
	t.Tags("Negative")

	refreshToken, jti := createRefreshToken(t)

	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Run(
		func(args mock.Arguments) {
			jtiPtr := args.Get(2).(*string)
			*jtiPtr = jti
		},
	).Return(nil)

	_, err := svc.GetRefreshTokenClaims(ctx, refreshToken)

	t.Require().Error(err)
	t.Require().ErrorIs(err, infrastructure.ErrBannedJWTToken)
}

func (suite *GetRefreshTokenClaimsSuite) Test_ErrGetCache(t provider.T) {
	t.Title("Returns getting cache error")
	t.Severity(allure.CRITICAL)
	t.Epic("JWT service")
	t.Feature("GetRefreshTokenClaimsSuite")
	t.Tags("Negative")

	refreshToken, _ := createRefreshToken(t)

	cacheErr := errors.New("cache error")
	managerMock.On("Get", ctx, mock.Anything, mock.Anything).Once().Return(cacheErr)

	_, err := svc.GetRefreshTokenClaims(ctx, refreshToken)

	t.Require().Error(err)
	t.Require().ErrorIs(err, cacheErr)
}
