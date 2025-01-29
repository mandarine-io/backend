package auth

import (
	"context"
	"github.com/go-resty/resty/v2"
	json2 "github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/tests/e2e"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/ozontech/cute"
	"net/http"
	"strings"
	"time"
)

type ResetPasswordSuite struct {
	suite.Suite
}

func (s *ResetPasswordSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('5d260929-4c31-4441-9793-7558e52b8720', 'user_for_reset_password', 'user_for_reset_password@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL), " +
			"('5d260929-4c31-4441-9793-7558e52b8721', 'user_for_reset_password_1', 'user_for_reset_password_1@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL) " +
			"ON CONFLICT DO NOTHING",
	)
	t.Require().NoError(tx.Error)
}

func (s *ResetPasswordSuite) AfterEach(t provider.T) {
	tx := db.Exec("DELETE FROM users WHERE id = '5d260929-4c31-4441-9793-7558e52b8720' AND  id = '5d260929-4c31-4441-9793-7558e52b8721'")
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

func (s *ResetPasswordSuite) Test_Success(t provider.T) {
	cute.NewTestBuilder().
		Title("RecoveryPassword returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("RecoveryPassword").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/recovery-password"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RecoveryPasswordInput{
						Email: "user_for_reset_password@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusAccepted).
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
		if message.Raw.To[0] == "user_for_reset_password@mail.ru" {
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
		currentKeys, cursor, err = rdb.Scan(context.Background(), cursor, "recovery_password.*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, currentKeys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	t.Require().Len(keys, 1)

	otp := strings.Replace(keys[0], "recovery_password.", "", 1)

	cute.NewTestBuilder().
		Title("VerifyRecoveryCode returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("VerifyRecoveryCode").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/recovery-password/verify"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.VerifyRecoveryCodeInput{
						OTP:   otp,
						Email: "user_for_reset_password@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)

	cute.NewTestBuilder().
		Title("ResetPassword returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("ResetPassword").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/reset-password"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.ResetPasswordInput{
						OTP:      otp,
						Email:    "user_for_reset_password@mail.ru",
						Password: uuid.New().String(),
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)
}

func (s *ResetPasswordSuite) Test_UserNotFound(t provider.T) {
	cute.NewTestBuilder().
		Title("RecoveryPassword returns UserNotFound").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("RecoveryPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/recovery-password"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RecoveryPasswordInput{
						Email: "user_for_reset_password_not_existent@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *ResetPasswordSuite) Test_IncorrectOTP(t provider.T) {
	cute.NewTestBuilder().
		Title("RecoveryPassword returns IncorrectOTP").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("RecoveryPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/recovery-password"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RecoveryPasswordInput{
						Email: "user_for_reset_password_1@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusAccepted).
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
		if message.Raw.To[0] == "user_for_reset_password_1@mail.ru" {
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
		currentKeys, cursor, err = rdb.Scan(context.Background(), cursor, "recovery_password.*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, currentKeys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	t.Require().Len(keys, 1)

	otp := strings.Replace(keys[0], "recovery_password.", "", 1)

	cute.NewTestBuilder().
		Title("VerifyRecoveryCode returns IncorrectOTP").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("VerifyRecoveryCode").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/recovery-password/verify"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.VerifyRecoveryCodeInput{
						OTP:   otp + "1",
						Email: "user_for_reset_password_1@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)

	cute.NewTestBuilder().
		Title("ResetPassword returns IncorrectOTP").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("ResetPassword").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/reset-password"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.ResetPasswordInput{
						OTP:      otp + "1",
						Email:    "user_for_reset_password_1@mail.ru",
						Password: uuid.New().String(),
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)
}
