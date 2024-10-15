package auth_e2e_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"mandarine/internal/api/service/auth/dto"
	dto2 "mandarine/pkg/rest/dto"
	"mandarine/tests/e2e"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mailhogMessagesResponse struct {
	Total    int           `json:"total"`
	Start    int           `json:"start"`
	Count    int           `json:"count"`
	Messages []interface{} `json:"messages"`
}

func Test_RegisterHandler_Register(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/register"

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.RegisterInput{
			nil,

			// Bad username
			{
				Username: "User",
				Email:    "user_for_register@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			{
				Username: "1user",
				Email:    "user_for_register@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			{
				Username: "user!",
				Email:    "user_for_register@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			{
				Username: "",
				Email:    "user_for_register@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			{
				Username: strings.Repeat("user", 256),
				Email:    "user_for_register@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},

			// Bad email
			{
				Username: "user_for_register",
				Email:    "user_for_registerexample.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			{
				Username: "user_for_register",
				Email:    "",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},

			// Weak password
			{
				Username: "user_for_register",
				Email:    "user_for_register@example.com",
				Password: "weak",
			},
			{
				Username: "user_for_register",
				Email:    "user_for_register@example.com",
				Password: "",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				// Send request
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(*body)
				}

				req, _ := http.NewRequest("POST", url, reqBodyReader)
				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				// Check response
				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("User already exists", func(t *testing.T) {
		// Send request
		reqBody := dto.RegisterInput{
			Username: "user_for_register",
			Email:    "user_for_register@example.com",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		// Send request
		reqBody := dto.RegisterInput{
			Username: "new_user_for_register",
			Email:    "new_user_for_register@example.com",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Check cache
		code, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "get", "register.new_user_for_register@example.com"})
		assert.NoError(t, err)
		assert.Equal(t, 0, code)

		// Check email
		apiPort, err := testEnvironment.SmtpC.MappedPort(context.Background(), "8025")
		assert.NoError(t, err)

		mailhogResp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v2/messages?start=0&limit=10", apiPort.Port()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, mailhogResp.StatusCode)

		var mailhogMessages mailhogMessagesResponse
		err = e2e.ReadResponseBody(mailhogResp, &mailhogMessages)
		assert.NoError(t, err)
		assert.True(t, mailhogMessages.Total > 0)
	})
}

func Test_RegisterHandler_RegisterConfirm(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/register/confirm"

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.RegisterConfirmInput{
			nil,

			// Bad OTP
			{
				OTP:   "",
				Email: "user_for_register@example.com",
			},

			// Bad email
			{
				OTP:   "123456",
				Email: "",
			},
			{
				OTP:   "123456",
				Email: "user_for_register_confirmexample.com",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				// Send request
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(*body)
				}

				req, _ := http.NewRequest("POST", url, reqBodyReader)
				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				// Check response
				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("Not found cache entry", func(t *testing.T) {
		// Send request
		reqBody := dto.RegisterConfirmInput{
			OTP:   "123456",
			Email: "user_for_register_confirm@example.com",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Incorrect OTP", func(t *testing.T) {
		// Set cache
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Username: "user_for_register_confirm",
				Email:    "user_for_register_confirm@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(time.Hour),
		}
		err := testEnvironment.Container.CacheManager.Set(
			context.Background(),
			"register.user_for_register_confirm@example.com",
			cacheEntry,
		)
		assert.NoError(t, err)

		// Send request
		reqBody := dto.RegisterConfirmInput{
			OTP:   "654321",
			Email: "user_for_register_confirme@xample.com",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("User already exists", func(t *testing.T) {
		// Set cache
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Username: "user_for_register_confirm_duplicate",
				Email:    "user_for_register_confirm_duplicate@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(time.Hour),
		}
		err := testEnvironment.Container.CacheManager.Set(
			context.Background(),
			"register.user_for_register_confirm_duplicate@example.com",
			cacheEntry,
		)
		assert.NoError(t, err)

		// Send request
		reqBody := dto.RegisterConfirmInput{
			OTP:   "123456",
			Email: "user_for_register_confirm_duplicate@example.com",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		// Set cache
		cacheEntry := dto.RegisterCache{
			User: dto.RegisterInput{
				Username: "user_for_register_confirm",
				Email:    "user_for_register_confirm@example.com",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			OTP:       "123456",
			ExpiredAt: time.Now().Add(time.Hour),
		}
		err := testEnvironment.Container.CacheManager.Set(
			context.Background(),
			"register.user_for_register_confirm@example.com",
			cacheEntry,
		)
		assert.NoError(t, err)

		// Send request
		reqBody := dto.RegisterConfirmInput{
			OTP:   "123456",
			Email: "user_for_register_confirm@example.com",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
