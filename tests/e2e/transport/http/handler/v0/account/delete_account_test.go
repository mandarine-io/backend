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

type DeleteAccountSuite struct {
	suite.Suite
}

func (s *DeleteAccountSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30', 'user_for_deleteAccount', 'user_for_deleteAccount@mail.ru', 'password', '1', true, NULL)," +
			"('1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_deleteAccount_banned', 'user_for_deleteAccount_banned@mail.ru', 'password', '1', false, NULL)," +
			"('1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32', 'user_for_deleteAccount_deleted', 'user_for_deleteAccount_deleted@mail.ru', 'password', '1', true, NOW())",
	)
	t.Require().NoError(tx.Error)
}

func (s *DeleteAccountSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30' OR " +
			"id = '1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31' OR " +
			"id = '1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32'",
	)
	t.Require().NoError(tx.Error)
}

func (s *DeleteAccountSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_deleteAccount",
			"email":          "user_for_deleteAccount@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("DeleteAccount returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("DeleteAccount").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNoContent).
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
			json.Equal("$.isDeleted", "true"),
		).
		ExecuteTest(context.Background(), t)
}

func (s *DeleteAccountSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_deleteAccount_banned",
			"email":          "user_for_deleteAccount_banned@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      false,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("DeleteAccount returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("DeleteAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *DeleteAccountSuite) Test_UserIsDeleted(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_deleteAccount_deleted",
			"email":          "user_for_deleteAccount_deleted@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("DeleteAccount returns UserIsDeleted error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("DeleteAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusConflict).
		ExecuteTest(context.Background(), t)
}

func (s *DeleteAccountSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("DeleteAccount returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("DeleteAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *DeleteAccountSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "1cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("DeleteAccount returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("DeleteAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodDelete),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *DeleteAccountSuite) Test_MissingJwtToken(t provider.T) {
	cute.NewTestBuilder().
		Title("DeleteAccount returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("DeleteAccount").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodDelete),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
