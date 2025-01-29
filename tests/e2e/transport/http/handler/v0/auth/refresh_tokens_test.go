package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/tests/e2e"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/json"
	"net/http"
	"time"
)

type RefreshTokensSuite struct {
	suite.Suite
}

func (s *RefreshTokensSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('3d260929-4c31-4441-9793-7558e52b8720', 'user_for_refresh', 'user_for_refresh@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL), " +
			"('3d260929-4c31-4441-9793-7558e52b8721', 'user_for_refresh_banned', 'user_for_refresh_banned@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', false, NULL) " +
			"ON CONFLICT DO NOTHING",
	)
	t.Require().NoError(tx.Error)
}

func (s *RefreshTokensSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE id = '3d260929-4c31-4441-9793-7558e52b8720' AND " +
			"id = '3d260929-4c31-4441-9793-7558e52b8721'",
	)
	t.Require().NoError(tx.Error)
}

func (s *RefreshTokensSuite) Test_InvalidToken(t provider.T) {
	cute.NewTestBuilder().
		Title("Refresh returns Invalid token error").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Refresh").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/refresh"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RefreshTokensInput{
						RefreshToken: "refreshToken",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)
}

func (s *RefreshTokensSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  "mandarine",
			"sub":  "3d260929-4c31-4441-9793-7558e52b8722",
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Hour).Unix(),
			"jti":  jti,
			"type": "refresh",
		},
	)
	refreshTokenSigned, err := refreshToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("Refresh returns UserNotFound").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Refresh").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/refresh"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RefreshTokensInput{
						RefreshToken: refreshTokenSigned,
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *RefreshTokensSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  "mandarine",
			"sub":  "3d260929-4c31-4441-9793-7558e52b8721",
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Hour).Unix(),
			"jti":  jti,
			"type": "refresh",
		},
	)
	refreshTokenSigned, err := refreshToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("Refresh returns UserIsBanned").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Refresh").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/refresh"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RefreshTokensInput{
						RefreshToken: refreshTokenSigned,
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *RefreshTokensSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  "mandarine",
			"sub":  "3d260929-4c31-4441-9793-7558e52b8720",
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Hour).Unix(),
			"jti":  jti,
			"type": "refresh",
		},
	)
	refreshTokenSigned, err := refreshToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("Refresh returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("Refresh").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/refresh"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RefreshTokensInput{
						RefreshToken: refreshTokenSigned,
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.NotEmpty("$.accessToken"),
			json.NotEmpty("$.refreshToken"),
		).
		ExecuteTest(context.Background(), t)
}
