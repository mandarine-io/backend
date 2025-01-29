package account

import (
	"context"
	"github.com/go-resty/resty/v2"
	json2 "github.com/goccy/go-json"
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
	"strings"
	"time"
)

type UpdateEmailSuite struct {
	suite.Suite
}

func (s *UpdateEmailSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_email_verified, is_enabled, deleted_at) VALUES " +
			"('3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30', 'user_for_updateEmail', 'user_for_updateEmail@mail.ru', 'password', '1', true, true, NULL)," +
			"('3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39', 'user_for_updateEmail_existent', 'user_for_updateEmail_existent@mail.ru', 'password', '1', true, true, NULL)," +
			"('3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31', 'user_for_updateEmail_banned', 'user_for_updateEmail_banned@mail.ru', 'password', '1', true, false, NULL)," +
			"('3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32', 'user_for_updateEmail_deleted', 'user_for_updateEmail_deleted@mail.ru', 'password', '1', true, true, NOW())",
	)
	t.Require().NoError(tx.Error)
}

func (s *UpdateEmailSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30' OR " +
			"id = '3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39' OR " +
			"id = '3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31' OR " +
			"id = '3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32'",
	)
	t.Require().NoError(tx.Error)

	// Invalidate cache
	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = rdb.Scan(context.Background(), cursor, "*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		err := rdb.Del(context.Background(), keys...).Err()
		t.Require().NoError(err)
	}
}

func (s *UpdateEmailSuite) Test_Success(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateEmail",
			"email":          "user_for_updateEmail@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateEmail returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusAccepted).
		AssertBody(
			json.Equal("$.email", "user_for_updateEmail_updated@mail.ru"),
			json.Equal("$.isEmailVerified", "false"),
		).
		ExecuteTest(context.Background(), t)

	// Check mail
	client := resty.New()

	resp, err := client.R().Get(mailhogURL + "/api/v2/messages")
	t.Require().NoError(err)
	t.Require().Equal(200, resp.StatusCode())

	messages := GetMessagesResponse{}
	err = json2.Unmarshal(resp.Body(), &messages)
	t.Require().NoError(err)

	existsMail := false
	for _, message := range messages.Items {
		if message.Raw.To[0] == "user_for_updateEmail_updated@mail.ru" {
			existsMail = true
			break
		}
	}

	t.Require().True(existsMail)

	// Check cache
	var (
		cursor uint64
		keys   []string
	)
	for {
		var currentKeys []string
		currentKeys, cursor, err = rdb.Scan(context.Background(), cursor, "email_verify.*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, currentKeys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	t.Require().Len(keys, 1)

	otp := strings.Replace(keys[0], "email_verify.", "", 1)

	cute.NewTestBuilder().
		Title("UpdateEmail returns successfully").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email/verify"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.VerifyEmailInput{
						OTP:   otp,
						Email: "user_for_updateEmail_updated@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusOK).
		NextTest().
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account"),
			cute.WithMethod(http.MethodGet),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.Equal("$.email", "user_for_updateEmail_updated@mail.ru"),
			json.Equal("$.isEmailVerified", "true"),
		).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_SuccessNotUpdate(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a39",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateEmail_existent",
			"email":          "user_for_updateEmail_existent@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateEmail returns successfully, email not changed").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_existent@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.Equal("$.email", "user_for_updateEmail_existent@mail.ru"),
			json.Equal("$.isEmailVerified", "true"),
		).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_IncorrectEnteredOTP(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateEmail",
			"email":          "user_for_updateEmail@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateEmail returns IncorrectEnteredOTP error").
		Severity(allure.NORMAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated_1@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusAccepted).
		AssertBody(
			json.Equal("$.email", "user_for_updateEmail_updated_1@mail.ru"),
			json.Equal("$.isEmailVerified", "false"),
		).
		NextTest().
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email/verify"),
			cute.WithMethod(http.MethodPost),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(
				e2e.MustMarshal(
					t,
					v0.VerifyEmailInput{OTP: "incorrect", Email: "user_for_updateEmail_updated_1@mail.ru"},
				),
			),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_DuplicateEmail(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a30",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateEmail",
			"email":          "user_for_updateEmail@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateEmail returns DuplicateEmail error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_existent@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusConflict).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_UserIsBanned(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a31",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateEmail_banned",
			"email":          "user_for_updateEmail_banned@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      false,
			"isDeleted":      false,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateEmail returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_UserIsDeleted(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a32",
			"iat":            time.Now().Unix(),
			"exp":            time.Now().Add(time.Hour).Unix(),
			"jti":            jti,
			"type":           "access",
			"username":       "user_for_updateEmail_deleted",
			"email":          "user_for_updateEmail_deleted@mail.ru",
			"role":           "user",
			"IsPasswordTemp": false,
			"isEnabled":      true,
			"isDeleted":      true,
		},
	)
	accessTokenSigned, err := accessToken.SignedString([]byte(secret))
	t.Require().NoError(err)

	cute.NewTestBuilder().
		Title("UpdateEmail returns UserIsDeleted error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_UserNotFound(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("UpdateEmail returns UserNotFound error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_IncorrectBody_UpdateEmail(t provider.T) {
	t.Title("UpdateEmail returns IncorrectBody error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account handler")
	t.Feature("UpdateEmail")
	t.Tags("Negative")

	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		name  string
		email string
	}
	cases := []testCase{
		{"Empty email", ""},
		{"Invalid email", "not email"},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t provider.T) {
				cute.NewTestBuilder().
					Title(c.name).
					Create().
					RequestBuilder(
						cute.WithURI(serverURL+"/v0/account/email"),
						cute.WithMethod(http.MethodPatch),
						cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
						cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: c.email})),
					).
					ExpectExecuteTimeout(1*time.Minute).
					ExpectStatus(http.StatusBadRequest).
					ExecuteTest(context.Background(), t)
			},
		)
	}
}

func (s *UpdateEmailSuite) Test_IncorrectBody_VerifyEmail(t provider.T) {
	t.Title("VerifyEmail returns IncorrectBody error")
	t.Severity(allure.CRITICAL)
	t.Epic("Account handler")
	t.Feature("VerifyEmail")
	t.Tags("Negative")

	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		name  string
		email string
		otp   string
	}
	cases := []testCase{
		{"Empty email", "", "111111"},
		{"Invalid email", "not email", "111111"},
		{"Empty OTP", "email@mail.ru", ""},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t provider.T) {
				cute.NewTestBuilder().
					Title(c.name).
					Create().
					RequestBuilder(
						cute.WithURI(serverURL+"/v0/account/email/verify"),
						cute.WithMethod(http.MethodPost),
						cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
						cute.WithBody(e2e.MustMarshal(t, v0.VerifyEmailInput{Email: c.email, OTP: c.otp})),
					).
					ExpectExecuteTimeout(1*time.Minute).
					ExpectStatus(http.StatusBadRequest).
					ExecuteTest(context.Background(), t)
			},
		)
	}
}

func (s *UpdateEmailSuite) Test_ExpiredToken(t provider.T) {
	jti := uuid.New().String()
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":            "mandarine",
			"sub":            "3cf34f8c-1bd2-4b2a-8e4f-50dfe5c94a33",
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
		Title("UpdateEmail returns ExpiredToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithHeadersKV("Authorization", "Bearer "+accessTokenSigned),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}

func (s *UpdateEmailSuite) Test_MissingJwtToken(t provider.T) {

	cute.NewTestBuilder().
		Title("UpdateEmail returns MissingJwtToken error").
		Severity(allure.CRITICAL).
		Epic("Account handler").
		Feature("UpdateEmail").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/account/email"),
			cute.WithMethod(http.MethodPatch),
			cute.WithBody(e2e.MustMarshal(t, v0.UpdateEmailInput{Email: "user_for_updateEmail_updated@mail.ru"})),
		).
		ExpectExecuteTimeout(1*time.Minute).
		ExpectStatus(http.StatusUnauthorized).
		ExecuteTest(context.Background(), t)
}
