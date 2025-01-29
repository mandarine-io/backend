package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/ozontech/cute"
	"net/http"
	"time"
)

type LogoutSuite struct {
	suite.Suite
}

func (s *LogoutSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('2d260929-4c31-4441-9793-7558e52b8720', 'user_for_logout', 'user_for_logout@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL) " +
			"ON CONFLICT DO NOTHING",
	)
	t.Require().NoError(tx.Error)
}

func (s *LogoutSuite) AfterEach(t provider.T) {
	tx := db.Exec("DELETE FROM users WHERE id = '2d260929-4c31-4441-9793-7558e52b8720'")
	t.Require().NoError(tx.Error)
}

func (s *LogoutSuite) Test_Unauthorized(t provider.T) {
	cute.NewTestBuilder().
		Title("Logout returns Unauthorized error").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Logout").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/logout"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *LogoutSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "2d260929-4c31-4441-9793-7558e52b8720",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_logout",
			"email":          "user_for_logout@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("Logout returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("Logout").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/logout"),
			cute.WithHeadersKV("Authorization", accessTokenSigned),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/logout"),
			cute.WithHeadersKV("Authorization", accessTokenSigned),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
