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

type RestoreAccountSuite struct {
	suite.Suite
}

func (s *RestoreAccountSuite) BeforeAll(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30', 'user_for_restoreAccount', 'user_for_restoreAccount@mail.ru', 'password', '1', true, NOW())," +
			"('2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_restoreAccount_banned', 'user_for_restoreAccount_banned@mail.ru', 'password', '1', false, NULL)," +
			"('2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32', 'user_for_restoreAccount_restored', 'user_for_restoreAccount_restored@mail.ru', 'password', '1', true, NULL)",
	)
	t.Require().NoError(tx.Error)
}

func (s *RestoreAccountSuite) AfterAll(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30' OR " +
			"id = '2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31' OR " +
			"id = '2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32'",
	)
	t.Require().NoError(tx.Error)
}

func (s *RestoreAccountSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_restoreAccount",
			"email":          "user_for_restoreAccount@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("RestoreAccount returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("RestoreAccount").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/restore"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		NextTest().
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.Equal("$.isDeleted", "false"),
		).
		ExecuteTest(context.Background(), t)
}

func (s *RestoreAccountSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_restoreAccount_banned",
			"email":          "user_for_restoreAccount_banned@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      false,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("RestoreAccount returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("RestoreAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/restore"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *RestoreAccountSuite) Test_UserIsNotDeleted(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_restoreAccount_restored",
			"email":          "user_for_restoreAccount_restored@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("RestoreAccount returns UserIsNotDeleted error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("RestoreAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/restore"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusConflict).
		ExecuteTest(context.Background(), t)
}

func (s *RestoreAccountSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "non_existent_user",
			"email":          "non_existent_user@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("RestoreAccount returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("RestoreAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/restore"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *RestoreAccountSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":        "mandarine",
			"sub":        "2cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
			"iat":        time.Now().Unix(),
			"exp":        time.Now().Add(-time.Hour).Unix(),
			"jti":        jti,
			"type":       "access",
			"username":   "non_existent_user",
			"email":      "non_existent_user@mail.ru",
			"role":       "user",
			"isEnabled":  true,
			"isRestored": false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("RestoreAccount returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("RestoreAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/restore"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *RestoreAccountSuite) Test_MissingJwtToken(t provider.T) {
	cute.NewTestBuilder().
		Title("RestoreAccount returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("RestoreAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/restore"),
			cute.WithMethod(http.MethodGet),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
