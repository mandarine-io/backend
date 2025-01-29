package auth

import (
	"context"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/tests/e2e"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/ozontech/cute"
	"github.com/ozontech/cute/asserts/json"
	"net/http"
	"time"
)

type LoginSuite struct {
	suite.Suite
}

func (s *LoginSuite) BeforeEach(t provider.T) {
	tx := db.Exec(
		"INSERT INTO users(id, username, email, password, role_id, is_enabled, deleted_at) VALUES " +
			"('1d260929-4c31-4441-9793-7558e52b8720', 'user_for_login', 'user_for_login@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', true, NULL)," +
			"('1d260929-4c31-4441-9793-7558e52b8721', 'user_for_login_banned', 'user_for_login_banned@mail.ru', '$2a$12$7BsgTSO6Yg3FFS8dkVYcre4BLVOCp.8x8fyBAG7cDRxdkbIdOkgeS', '1', false, NULL) " +
			"ON CONFLICT DO NOTHING",
	)
	t.Require().NoError(tx.Error)
}

func (s *LoginSuite) AfterEach(t provider.T) {
	tx := db.Exec(
		"DELETE FROM users WHERE " +
			"id = '1d260929-4c31-4441-9793-7558e52b8720' OR " +
			"id = '1d260929-4c31-4441-9793-7558e52b8721'",
	)
	t.Require().NoError(tx.Error)
}

func (s *LoginSuite) Test_BadBody(t provider.T) {
	t.Title("Login returns BadBody error")
	t.Severity(allure.CRITICAL)
	t.Epic("Auth handler")
	t.Feature("Login")
	t.Tags("Negative")

	type testCase struct {
		name     string
		login    string
		password string
	}
	cases := []testCase{
		{"Empty login", "", "7676393c-0d28-4f68-807d-12aa6b88c039"},
		{"Empty password", "user", ""},
	}

	for _, c := range cases {
		t.Run(
			c.name, func(t provider.T) {
				cute.NewTestBuilder().
					Title(c.name).
					Create().
					RequestBuilder(
						cute.WithURI(serverURL+"/v0/auth/login"),
						cute.WithMethod(http.MethodPost),
						cute.WithBody(e2e.MustMarshal(t, v0.LoginInput{Login: c.login, Password: c.password})),
					).
					ExpectExecuteTimeout(10*time.Second).
					ExpectStatus(http.StatusBadRequest).
					ExecuteTest(context.Background(), t)
			},
		)
	}
}

func (s *LoginSuite) Test_UserNotFound(t provider.T) {
	cute.NewTestBuilder().
		Title("Login returns UserIsDeleted error").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Login").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/login"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.LoginInput{
						Login:    "user_for_login_non_exist@mail.ru",
						Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusNotFound).
		ExecuteTest(context.Background(), t)
}

func (s *LoginSuite) Test_BadCredentials(t provider.T) {
	cute.NewTestBuilder().
		Title("Login returns BadCredentials error").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Login").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/login"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.LoginInput{
						Login:    "user_for_login@mail.ru",
						Password: "bad_password",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusBadRequest).
		ExecuteTest(context.Background(), t)
}

func (s *LoginSuite) Test_UserIsBanned(t provider.T) {
	cute.NewTestBuilder().
		Title("Login returns UserIsBanned error").
		Severity(allure.CRITICAL).
		Epic("Auth handler").
		Feature("Login").
		Tags("Negative").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/login"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.LoginInput{
						Login:    "user_for_login_banned@mail.ru",
						Password: "password",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusForbidden).
		ExecuteTest(context.Background(), t)
}

func (s *LoginSuite) Test_SuccessWithEmail(t provider.T) {
	cute.NewTestBuilder().
		Title("Login returns success with email").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("Login").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/login"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.LoginInput{
						Login:    "user_for_login@mail.ru",
						Password: "password",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.NotEmpty("$.accessToken"),
			json.NotEmpty("$.refreshToken"),
		).
		ExecuteTest(context.Background(), t)
}

func (s *LoginSuite) Test_SuccessWithUsername(t provider.T) {
	cute.NewTestBuilder().
		Title("Login returns success with username").
		Severity(allure.NORMAL).
		Epic("Auth handler").
		Feature("Login").
		Tags("Positive").
		Create().
		RequestBuilder(
			cute.WithURI(serverURL+"/v0/auth/login"),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(
				e2e.MustMarshal(
					t, v0.LoginInput{
						Login:    "user_for_login",
						Password: "password",
					},
				),
			),
		).
		ExpectExecuteTimeout(10*time.Second).
		ExpectStatus(http.StatusOK).
		AssertBody(
			json.NotEmpty("$.accessToken"),
			json.NotEmpty("$.refreshToken"),
		).
		ExecuteTest(context.Background(), t)
}
