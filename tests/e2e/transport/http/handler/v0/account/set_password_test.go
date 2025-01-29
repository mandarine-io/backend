package account

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
	"net/http"
	"time"
)

type SetPasswordSuite struct {
	suite.Suite
}

func (s *SetPasswordSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, is_password_temp, role_id, is_enabled, deleted_at) VALUES " +
			"('5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30', 'user_for_setPassword', 'user_for_setPassword@mail.ru', 'password', true, '1', true, NULL)," +
			"('5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39', 'user_for_setPassword_set', 'user_for_setPassword_set@mail.ru', 'password', false, '1', true, NULL)," +
			"('5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_setPassword_banned', 'user_for_setPassword_banned@mail.ru', 'password', true, '1', false, NULL)," +
			"('5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32', 'user_for_setPassword_deleted', 'user_for_setPassword_deleted@mail.ru', 'password', true, '1', true, NOW())",
	)
	t.Require().NoError(tx.Error)
}

func (s *SetPasswordSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30' OR " +
			"id = '5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39' OR " +
			"id = '5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31' OR " +
			"id = '5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32'",
	)
	t.Require().NoError(tx.Error)
}

func (s *SetPasswordSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_setPassword",
			"email":          "user_for_setPassword@mail.ru",
			"role":           "user",
			"IsPasswordTemp": true,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("SetPassword returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.SetPasswordInput{Password: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)
}

func (s *SetPasswordSuite) Test_PasswordAlreadySet(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_setPassword_set",
			"email":          "user_for_setPassword_set@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("SetPassword returns PasswordAlreadySet error").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.SetPasswordInput{Password: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusConflict).
		ExecuteTest(context.Background(), t)
}

func (s *SetPasswordSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_setPassword_banned",
			"email":          "user_for_setPassword_banned@mail.ru",
			"role":           "user",
			"IsPasswordTemp": true,
			"isEnabled":      false,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("SetPassword returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.SetPasswordInput{Password: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *SetPasswordSuite) Test_UserIsDeleted(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_setPassword_deleted",
			"email":          "user_for_setPassword_deleted@mail.ru",
			"role":           "user",
			"IsPasswordTemp": true,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("SetPassword returns UserIsDeleted error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.SetPasswordInput{Password: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *SetPasswordSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("SetPassword returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.SetPasswordInput{Password: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *SetPasswordSuite) Test_IncorrectBody(t provider.T) {
	t.Title("SetPassword returns IncorrectBody error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account handler")
	t.Feature("SetPassword")
	t.Tags("Negative")

	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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

	type testCase struct {
		name     string
		password string
	}
	cases := []testCase{
		{"Empty new password", ""},
		{"Weak new password", "weak"},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t provider.T) {
				cute.NewTestBuilder().
					Title(c.name).
					Create().
					RequestBuilder(
						cute.WithURI(serverURL+"/v0/account/password"),
						cute.WithMethod(http.MethodPost),
						cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
						cute.WithBody(e2e.MustMarshal(t, v0.SetPasswordInput{Password: c.password})),
					).
					ExpectExecuteTimeout(10*time.Second).
					ExpectStatus(http.StatusBadRequest).
					ExecuteTest(context.Background(), t)
			},
		)
	}
}

func (s *SetPasswordSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "5cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("SetPassword returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.SetPasswordInput{Password: "password"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *SetPasswordSuite) Test_MissingJwtToken(t provider.T) {

	cute.NewTestBuilder().
		Title("SetPassword returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("SetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(e2e.MustMarshal(t, v0.SetPasswordInput{Password: "password"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
