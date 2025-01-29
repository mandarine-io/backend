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

type RegisterSuite struct {
	suite.Suite
}

func (s *RegisterSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('4d260929-4c31-4441-9793-7558e52b8720', 'user_for_register_exists', 'user_for_register_exists@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL) " +
			"ON CONFLICT DO NOTHING",
	)
	t.Require().NoError(tx.Error)
}

func (s *RegisterSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE id = '4d260929-4c31-4441-9793-7558e52b8720' OR " +
			"username = 'user_for_register'",
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

func (s *RegisterSuite) Test_Success(t provider.T) {
	password := uuid.New().String()

	cute.NewTestBuilder().
		Title("Register returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("Register").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/register"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RegisterInput{
						Username: "user_for_register",
						Email:    "user_for_register@mail.ru",
						Password: password,
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
		if message.Raw.To[0] == "user_for_register@mail.ru" {
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
		currentKeys, cursor, err = rdb.Scan(context.Background(), cursor, "register.*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, currentKeys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	t.Require().Len(keys, 1)

	otp := strings.Replace(keys[0], "register.", "", 1)

	cute.NewTestBuilder().
		Title("Register confirm returns successfully").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("RegisterConfirm").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/register/confirm"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RegisterConfirmInput{
						OTP:   otp,
						Email: "user_for_register@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		ExecuteTest(context.Background(), t)
}

func (s *RegisterSuite) Test_IncorrectOTP(t provider.T) {
	password := uuid.New().String()

	cute.NewTestBuilder().
		Title("Register returns IncorrectOTP").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Register").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/register"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RegisterInput{
						Username: "user_for_register",
						Email:    "user_for_register@mail.ru",
						Password: password,
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
		if message.Raw.To[0] == "user_for_register@mail.ru" {
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
		currentKeys, cursor, err = rdb.Scan(context.Background(), cursor, "register.*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, currentKeys...)

		if cursor == 0 { // no more keys
			break
		}
	}

	t.Require().Len(keys, 1)

	otp := strings.Replace(keys[0], "register.", "", 1)

	cute.NewTestBuilder().
		Title("Register confirm returns IncorrectOTP").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("RegisterConfirm").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/register/confirm"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RegisterConfirmInput{
						OTP:   otp + "1",
						Email: "user_for_register@mail.ru",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)
}

func (s *RegisterSuite) Test_DuplicateUser(t provider.T) {
	password := uuid.New().String()

	cute.NewTestBuilder().
		Title("Register returns DuplicateUser").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Register").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/register"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.RegisterInput{
						Username: "user_for_register_exists",
						Email:    "user_for_register_exists@mail.ru",
						Password: password,
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusConflict).
		ExecuteTest(context.Background(), t)
}
