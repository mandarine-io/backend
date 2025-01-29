package account

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/json"
	"net/http"
	"time"
)

type GetAccountSuite struct {
	suite.Suite
}

func (s *GetAccountSuite) BeforeAll(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id) VALUES " +
			"('dcf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_getAccount', 'user_for_getAccount@mail.ru', 'password', '1')",
	)
	t.Require().NoError(tx.Error)
}

func (s *GetAccountSuite) AfterAll(t provider.T) {
	tx := db.Exec("DELETE FROM users WHERE id = 'dcf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31'")
	t.Require().NoError(tx.Error)
}

func (s *GetAccountSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "dcf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_getAccount",
			"email":          "user_for_getAccount@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("GetAccount returns account successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("GetAccount").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.Equal("$.username", "user_for_getAccount"),
			json.Equal("$.email", "user_for_getAccount@mail.ru"),
		).
		ExecuteTest(context.Background(), t)
}

func (s *GetAccountSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "dcf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "non_existent_user",
			"email":          "non_existent_user@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("GetAccount returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("GetAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *GetAccountSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "dcf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(-time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "non_existent_user",
			"email":          "non_existent_user@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("GetAccount returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("GetAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *GetAccountSuite) Test_MissingJwtToken(t provider.T) {
	cute.NewTestBuilder().
		Title("GetAccount returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("GetAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
