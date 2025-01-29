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

type UpdateUsernameSuite struct {
	suite.Suite
}

func (s *UpdateUsernameSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30', 'user_for_updateUsername', 'user_for_updateUsername@mail.ru', 'password', '1', true, NULL)," +
			"('6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39', 'user_for_updateusername_exist', 'user_for_updateusername_exist@mail.ru', 'password', '1', true, NULL)," +
			"('6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_updateUsername_banned', 'user_for_updateUsername_banned@mail.ru', 'password', '1', false, NULL)," +
			"('6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32', 'user_for_updateUsername_deleted', 'user_for_updateUsername_deleted@mail.ru', 'password', '1', true, NOW())",
	)
	t.Require().NoError(tx.Error)
}

func (s *UpdateUsernameSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30' OR " +
			"id = '6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39' OR " +
			"id = '6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31' OR " +
			"id = '6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32'",
	)
	t.Require().NoError(tx.Error)
}

func (s *UpdateUsernameSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateUsername",
			"email":          "user_for_updateUsername@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateUsername returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "username"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateUsernameSuite) Test_UsernameExists(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateUsername",
			"email":          "user_for_updateUsername@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateUsername returns UsernameExists error").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "user_for_updateusername_exist"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusConflict).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateUsernameSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateUsername_banned",
			"email":          "user_for_updateUsername_banned@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      false,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateUsername returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "username"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateUsernameSuite) Test_UserIsDeleted(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateUsername_deleted",
			"email":          "user_for_updateUsername_deleted@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateUsername returns UserIsDeleted error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "username"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateUsernameSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("UpdateUsername returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "username"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateUsernameSuite) Test_IncorrectBody(t provider.T) {
	t.Title("UpdateUsername returns IncorrectBody error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account handler")
	t.Feature("UpdateUsername")
	t.Tags("Negative")

	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		username string
	}
	cases := []testCase{
		{"Empty new username", ""},
		{"Invalid new username 1", "A"},
		{"Invalid new username 2", "_username"},
		{"Invalid new username 3", "username1!"},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t provider.T) {
				cute.NewTestBuilder().
					Title(c.name).
					Create().
					RequestBuilder(
						cute.WithURI(serverURL+"/v0/account/username"),
						cute.WithMethod(http.MethodPatch),
						cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
						cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: c.username})),
					).
					ExpectExecuteTimeout(10*time.Second).
					ExpectStatus(http.StatusBadRequest).
					ExecuteTest(context.Background(), t)
			},
		)
	}
}

func (s *UpdateUsernameSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "6cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("UpdateUsername returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "username"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateUsernameSuite) Test_MissingJwtToken(t provider.T) {

	cute.NewTestBuilder().
		Title("UpdateUsername returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateUsername").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/username"),
			cute.WithMethod(http.MethodPatch),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateUsernameInput{Username: "username"})),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
