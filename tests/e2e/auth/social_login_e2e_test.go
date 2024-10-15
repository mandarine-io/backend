package auth_e2e_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"io"
	"mandarine/internal/api/service/auth/dto"
	"mandarine/pkg/oauth"
	mock3 "mandarine/pkg/oauth/mock"
	dto2 "mandarine/pkg/rest/dto"
	"mandarine/tests/e2e"
	"net/http"
	"strings"
	"testing"
)

func Test_SocialLogin_SocialLogin(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	prefixUrl := server.URL + "/v0/auth/social"
	oauthProviderMock := testEnvironment.Container.OauthProviders["mock"].(*mock3.OAuthProviderMock)

	t.Run("Unsupported provider", func(t *testing.T) {
		// Send request
		url := prefixUrl + "/unsupported-provider"
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Not redirect url", func(t *testing.T) {
		// Send request
		url := prefixUrl + "/mock"
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		url := prefixUrl + "/mock?redirectUrl=https://mandarine.dev"

		// Mock
		oauthProviderMock.On("GetConsentPageUrl", "https://mandarine.dev").Return(server.URL+"/echo", "oauthState").Once()

		// Send request
		resp, err := server.Client().Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("Referer"), server.URL+"/echo")
	})
}

func Test_SocialLogin_SocialCallback(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	prefixUrl := server.URL + "/v0/auth/social"
	oauthProviderMock := testEnvironment.Container.OauthProviders["mock"].(*mock3.OAuthProviderMock)

	t.Run("Unsupported provider", func(t *testing.T) {
		// Send request
		url := prefixUrl + "/unsupported-provider/callback"

		req, _ := http.NewRequest("POST", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.SocialLoginCallbackInput{
			nil,
			{
				Code:  "",
				State: "oauthState",
			},
			{
				Code:  "code",
				State: "",
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				url := prefixUrl + "/mock/callback"

				// Send request
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				req, _ := http.NewRequest("POST", url, reqBodyReader)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("Not found state cookie", func(t *testing.T) {
		url := prefixUrl + "/mock/callback"
		reqBody := dto.SocialLoginCallbackInput{
			Code:  "code",
			State: "oauthState",
		}

		// Send request
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Mismatch state", func(t *testing.T) {
		url := prefixUrl + "/mock/callback"
		reqBody := dto.SocialLoginCallbackInput{
			Code:  "code",
			State: "oauthState",
		}

		// Send request
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "anotherOauthState"})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Failed to exchange code", func(t *testing.T) {
		url := prefixUrl + "/mock/callback"
		reqBody := dto.SocialLoginCallbackInput{
			Code:  "code",
			State: "oauthState",
		}

		// Mock
		oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, "code", mock.Anything).Return(nil, errors.New("failed to exchange code")).Once()

		// Send request
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "oauthState"})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Failed to get user info", func(t *testing.T) {
		url := prefixUrl + "/mock/callback"
		reqBody := dto.SocialLoginCallbackInput{
			Code:  "code",
			State: "oauthState",
		}

		// Mock
		token := &oauth2.Token{AccessToken: "accessToken"}
		oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, "code", mock.Anything).Return(token, nil).Once()
		oauthProviderMock.On("GetUserInfo", mock.Anything, token).Return(oauth.UserInfo{}, errors.New("failed to get user info")).Once()

		// Send request
		reqBodyReader, _ := e2e.NewJSONReader(reqBody)
		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "oauthState"})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("User with such email already exists", func(t *testing.T) {
		url := prefixUrl + "/mock/callback"
		reqBody := dto.SocialLoginCallbackInput{
			Code:  "code",
			State: "oauthState",
		}

		t.Run("User is blocked", func(t *testing.T) {
			// Mock
			token := &oauth2.Token{AccessToken: "accessToken"}
			userInfo := oauth.UserInfo{
				Email:           "user_for_social_blocked@example.com",
				Username:        "user_for_social_blocked",
				IsEmailVerified: true,
			}
			oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, "code", mock.Anything).Return(token, nil).Once()
			oauthProviderMock.On("GetUserInfo", mock.Anything, token).Return(userInfo, nil).Once()

			// Send request
			reqBodyReader, _ := e2e.NewJSONReader(reqBody)
			req, _ := http.NewRequest("POST", url, reqBodyReader)
			req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "oauthState"})

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)

			var body dto2.ErrorResponse
			err = e2e.ReadResponseBody(resp, &body)
			assert.NoError(t, err)
		})

		t.Run("Success", func(t *testing.T) {
			// Mock
			token := &oauth2.Token{AccessToken: "accessToken"}
			userInfo := oauth.UserInfo{
				Email:           "user_for_social@example.com",
				Username:        "user_for_social",
				IsEmailVerified: true,
			}
			oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, "code", mock.Anything).Return(token, nil).Once()
			oauthProviderMock.On("GetUserInfo", mock.Anything, token).Return(userInfo, nil).Once()

			// Send request
			reqBodyReader, _ := e2e.NewJSONReader(reqBody)
			req, _ := http.NewRequest("POST", url, reqBodyReader)
			req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "oauthState"})

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Check body
			var body dto.JwtTokensOutput
			err = e2e.ReadResponseBody(resp, &body)
			assert.NoError(t, err)
			assert.NotEmpty(t, body.AccessToken)
			assert.Empty(t, body.RefreshToken)

			// Check header
			setCookies := resp.Header["Set-Cookie"]
			for _, setCookie := range setCookies {
				if strings.Contains(setCookie, "RefreshToken") {
					assert.NotEmpty(t, setCookie)
					assert.Contains(t, setCookie, "RefreshToken=")
					return
				}
			}
			assert.Fail(t, "Refresh token cookie not found")
		})
	})

	t.Run("User with such email not exists", func(t *testing.T) {
		url := prefixUrl + "/mock/callback"
		reqBody := dto.SocialLoginCallbackInput{
			Code:  "code",
			State: "oauthState",
		}

		t.Run("Without searching unique username", func(t *testing.T) {
			// Mock
			token := &oauth2.Token{AccessToken: "accessToken"}
			userInfo := oauth.UserInfo{
				Email:           "user_for_social_unique@example.com",
				Username:        "user_for_social_unique",
				IsEmailVerified: true,
			}
			oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, "code", mock.Anything).Return(token, nil).Once()
			oauthProviderMock.On("GetUserInfo", mock.Anything, token).Return(userInfo, nil).Once()

			// Send request
			reqBodyReader, _ := e2e.NewJSONReader(reqBody)
			req, _ := http.NewRequest("POST", url, reqBodyReader)
			req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "oauthState"})

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Check body
			var body dto.JwtTokensOutput
			err = e2e.ReadResponseBody(resp, &body)
			assert.NoError(t, err)
			assert.NotEmpty(t, body.AccessToken)
			assert.Empty(t, body.RefreshToken)

			// Check header
			setCookies := resp.Header["Set-Cookie"]
			for _, setCookie := range setCookies {
				if strings.Contains(setCookie, "RefreshToken") {
					assert.NotEmpty(t, setCookie)
					assert.Contains(t, setCookie, "RefreshToken=")
					return
				}
			}
			assert.Fail(t, "Refresh token cookie not found")
		})

		t.Run("With searching unique username", func(t *testing.T) {
			// Mock
			token := &oauth2.Token{AccessToken: "accessToken"}
			userInfo := oauth.UserInfo{
				Email:           "user_for_social_unique_1@example.com",
				Username:        "user_for_social",
				IsEmailVerified: true,
			}
			oauthProviderMock.On("ExchangeCodeToToken", mock.Anything, "code", mock.Anything).Return(token, nil).Once()
			oauthProviderMock.On("GetUserInfo", mock.Anything, token).Return(userInfo, nil).Once()

			// Send request
			reqBodyReader, _ := e2e.NewJSONReader(reqBody)
			req, _ := http.NewRequest("POST", url, reqBodyReader)
			req.AddCookie(&http.Cookie{Name: "OAuthState", Value: "oauthState"})

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Check body
			var body dto.JwtTokensOutput
			err = e2e.ReadResponseBody(resp, &body)
			assert.NoError(t, err)
			assert.NotEmpty(t, body.AccessToken)
			assert.Empty(t, body.RefreshToken)

			// Check header
			setCookies := resp.Header["Set-Cookie"]
			for _, setCookie := range setCookies {
				if strings.Contains(setCookie, "RefreshToken") {
					assert.NotEmpty(t, setCookie)
					assert.Contains(t, setCookie, "RefreshToken=")
					return
				}
			}
			assert.Fail(t, "Refresh token cookie not found")
		})
	})
}
