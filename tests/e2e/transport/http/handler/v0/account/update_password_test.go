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

type UpdatePasswordSuite struct {
	suite.Suite
}

func (s *UpdatePasswordSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30', 'user_for_updatePassword', 'user_for_updatePassword@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL)," +
			"('4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_updatePassword_banned', 'user_for_updatePassword_banned@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', false, NULL)," +
			"('4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32', 'user_for_updatePassword_deleted', 'user_for_updatePassword_deleted@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NOW())",
	)
	t.Require().NoError(tx.Error)
}

func (s *UpdatePasswordSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30' OR " +
			"id = '4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31' OR " +
			"id = '4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32'",
	)
	t.Require().NoError(tx.Error)
}

func (s *UpdatePasswordSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updatePassword",
			"email":          "user_for_updatePassword@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdatePassword returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{
						OldPassword: "password",
						NewPassword: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)
}

func (s *UpdatePasswordSuite) Test_IncorrectOldPassword(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updatePassword",
			"email":          "user_for_updatePassword@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdatePassword returns IncorrectOldPassword error").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{
						OldPassword: "incorrect_password",
						NewPassword: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)
}

func (s *UpdatePasswordSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updatePassword_banned",
			"email":          "user_for_updatePassword_banned@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      false,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdatePassword returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{
						OldPassword: "password",
						NewPassword: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *UpdatePasswordSuite) Test_UserIsDeleted(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updatePassword_deleted",
			"email":          "user_for_updatePassword_deleted@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdatePassword returns UserIsDeleted error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{
						OldPassword: "password",
						NewPassword: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *UpdatePasswordSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("UpdatePassword returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{
						OldPassword: "password",
						NewPassword: "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *UpdatePasswordSuite) Test_IncorrectBody(t provider.T) {
	t.Title("UpdatePassword returns IncorrectBody error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account handler")
	t.Feature("UpdatePassword")
	t.Tags("Negative")

	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		name        string
		oldPassword string
		newPassword string
	}
	cases := []testCase{
		{"Empty old password", "", "593c5c35017e38864e266ddcee6e7aeca31288abf5e5e135f25f0ebf57a3fa63"},
		{"Empty new password", "password", ""},
		{"Weak new password", "password", "weak"},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t provider.T) {
				cute.NewTestBuilder().
					Title(c.name).
					Create().
					RequestBuilder(
						cute.WithURI(serverURL+"/v0/account/password"),
						cute.WithMethod(http.MethodPatch),
						cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
						cute.WithBody(
							e2e.MustMarshal(
								t,
								v0.UpdatePasswordInput{OldPassword: c.oldPassword, NewPassword: c.newPassword},
							),
						),
					).
					ExpectExecuteTimeout(10*time.Second).
					ExpectStatus(http.StatusBadRequest).
					ExecuteTest(context.Background(), t)
			},
		)
	}
}

func (s *UpdatePasswordSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "4cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("UpdatePassword returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{OldPassword: "password", NewPassword: "password"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *UpdatePasswordSuite) Test_MissingJwtToken(t provider.T) {

	cute.NewTestBuilder().
		Title("UpdatePassword returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdatePassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/password"),
			cute.WithMethod(http.MethodPatch),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.UpdatePasswordInput{OldPassword: "password", NewPassword: "password"},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
