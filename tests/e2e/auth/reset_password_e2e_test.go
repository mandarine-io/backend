package auth_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"mandarine/internal/api/service/auth/dto"
	dto2 "mandarine/pkg/rest/dto"
	"mandarine/tests/e2e"
	"net/http"
	"testing"
	"time"
)

func Test_ResetPasswordHandler_RecoveryPassword(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/recovery-password"

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.RecoveryPasswordInput{
			nil,

			// Bad email
			{
				Email: "",
			},
			{
				Email: "user_for_recovery_password_example.com",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				// Send request
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
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

	t.Run("User not found", func(t *testing.T) {
		// Send request
		reqBody := dto.RecoveryPasswordInput{
			Email: "user_for_recovery_password_1@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		// Send request
		reqBody := dto.RecoveryPasswordInput{
			Email: "user_for_recovery_password@example.com",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		// Check cache
		code, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "get", "recovery-password.user_for_recovery_password@example.com"})
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

func Test_ResetPasswordHandler_VerifyRecoveryCode(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/recovery-password/verify"

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.VerifyRecoveryCodeInput{
			nil,

			// Bad email
			{
				Email: "",
				OTP:   "123456",
			},
			{
				Email: "user_for_recovery_password_example.com",
				OTP:   "123456",
			},

			// Bad OTP
			{
				OTP:   "",
				Email: "user_for_recovery_password@example.com",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				// Send request
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
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

	t.Run("Cache miss", func(t *testing.T) {
		// Send request
		reqBody := dto.VerifyRecoveryCodeInput{
			Email: "user_for_verify_recovery_code@example.com",
			OTP:   "123456",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("OTP mismatch", func(t *testing.T) {
		// Add cache
		entry := dto.RecoveryPasswordCache{
			Email:     "user_for_verify_recovery_code@example.com",
			OTP:       "654321",
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		serializedEntry, _ := json.Marshal(entry)
		_, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "set", "recovery_password.user_for_verify_recovery_code@example.com", string(serializedEntry)})
		assert.NoError(t, err)

		// Send request
		reqBody := dto.VerifyRecoveryCodeInput{
			Email: "user_for_verify_recovery_code@example.com",
			OTP:   "123456",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// Add cache
		entry := dto.RecoveryPasswordCache{
			Email:     "user_for_verify_recovery_code@example.com",
			OTP:       "123456",
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		serializedEntry, _ := json.Marshal(entry)
		_, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "set", "recovery_password.user_for_verify_recovery_code@example.com", string(serializedEntry)})
		assert.NoError(t, err)

		// Send request
		reqBody := dto.VerifyRecoveryCodeInput{
			Email: "user_for_verify_recovery_code@example.com",
			OTP:   "123456",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func Test_ResetPasswordHandler_ResetPassword(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/reset-password"

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.ResetPasswordInput{
			nil,

			// Bad email
			{
				Email:    "",
				OTP:      "123456",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},
			{
				Email:    "user_for_recovery_password_example.com",
				OTP:      "123456",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			},

			//// Bad OTP
			//{
			//	OTP:      "",
			//	Email:    "user_for_recovery_password@example.com",
			//	Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
			//},
			//
			//// Bad password
			//{
			//	OTP:      "123456",
			//	Email:    "user_for_recovery_password@example.com",
			//	Password: "",
			//},
			//{
			//	OTP:      "123456",
			//	Email:    "user_for_recovery_password@example.com",
			//	Password: "weak",
			//},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				// Send request
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
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

	t.Run("Cache miss", func(t *testing.T) {
		// Send request
		reqBody := dto.ResetPasswordInput{
			Email:    "user_for_reset_password@example.com",
			OTP:      "123456",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("OTP mismatch", func(t *testing.T) {
		// Add cache
		entry := dto.RecoveryPasswordCache{
			Email:     "user_for_reset_password@example.com",
			OTP:       "654321",
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		serializedEntry, _ := json.Marshal(entry)
		_, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "set", "recovery_password.user_for_reset_password@example.com", string(serializedEntry)})
		assert.NoError(t, err)

		// Send request
		reqBody := dto.ResetPasswordInput{
			Email:    "user_for_reset_password@example.com",
			OTP:      "123456",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		// Add cache
		entry := dto.RecoveryPasswordCache{
			Email:     "user_for_reset_password_1@example.com",
			OTP:       "123456",
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		serializedEntry, _ := json.Marshal(entry)
		_, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "set", "recovery_password.user_for_reset_password_1@example.com", string(serializedEntry)})
		assert.NoError(t, err)

		// Send request
		reqBody := dto.ResetPasswordInput{
			Email:    "user_for_reset_password_1@example.com",
			OTP:      "123456",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// Add cache
		entry := dto.RecoveryPasswordCache{
			Email:     "user_for_reset_password@example.com",
			OTP:       "123456",
			ExpiredAt: time.Now().Add(1 * time.Hour),
		}
		serializedEntry, _ := json.Marshal(entry)
		_, _, err := testEnvironment.RedisC.Exec(context.Background(), []string{"redis-cli", "set", "recovery_password.user_for_reset_password@example.com", string(serializedEntry)})
		assert.NoError(t, err)

		// Send request
		reqBody := dto.ResetPasswordInput{
			Email:    "user_for_reset_password@example.com",
			OTP:      "123456",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
		}
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
