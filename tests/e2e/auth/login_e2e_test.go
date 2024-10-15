package auth_e2e_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	appconfig "mandarine/internal/api/config"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/rest"
	"mandarine/internal/api/service/auth/dto"
	dto2 "mandarine/pkg/rest/dto"
	validator3 "mandarine/pkg/rest/validator"
	"mandarine/tests/e2e"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	testEnvironment *e2e.TestEnvironment
	server          *httptest.Server
	fixtures        *testfixtures.Loader
)

func TestMain(m *testing.M) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cfg := &appconfig.Config{
		Server: appconfig.ServerConfig{
			Name:           "mandarine_test",
			Mode:           "test",
			ExternalOrigin: "http://localhost:8081",
			Port:           8081,
			Version:        "0.0.0",
		},
		Postgres: appconfig.PostgresConfig{
			Username: "mandarine",
			Password: "password",
			DBName:   "mandarine_test",
		},
		Redis: appconfig.RedisConfig{
			Host:     "127.0.0.1",
			Port:     6379,
			Username: "default",
			Password: "password",
			DBIndex:  0,
		},
		Minio: appconfig.MinioConfig{
			AccessKey:  "admin",
			SecretKey:  "Password_10",
			BucketName: "mandarine-test",
		},
		SMTP: appconfig.SmtpConfig{
			Host:     "127.0.0.1",
			Port:     25,
			Username: "admin",
			Password: "password",
			From:     "Mandarine <admin@localhost>",
			SSL:      false,
		},
		Cache: appconfig.CacheConfig{
			TTL: 120,
		},
		Locale: appconfig.LocaleConfig{
			Path:     pwd + "/../../../locales",
			Language: "ru",
		},
		Template: appconfig.TemplateConfig{
			Path: pwd + "/../../../templates",
		},
		Migrations: appconfig.MigrationConfig{
			Path: pwd + "/../../../migrations",
		},
		Logger: appconfig.LoggerConfig{
			Console: appconfig.ConsoleLoggerConfig{
				Level:    "debug",
				Encoding: "text",
			},
			File: appconfig.FileLoggerConfig{
				Enable: false,
			},
		},
		OAuthClient: appconfig.OAuthClientConfig{
			Google: appconfig.GoogleOAuthClientConfig{
				ClientID:     "",
				ClientSecret: "",
			},
			Yandex: appconfig.YandexOAuthClientConfig{
				ClientID:     "",
				ClientSecret: "",
			},
			MailRu: appconfig.MailRuOAuthClientConfig{
				ClientID:     "",
				ClientSecret: "",
			},
		},
		Security: appconfig.SecurityConfig{
			JWT: appconfig.JWTConfig{
				Secret:          "",
				AccessTokenTTL:  3600,
				RefreshTokenTTL: 86400,
			},
			OTP: appconfig.OTPConfig{
				Length: 6,
				TTL:    300,
			},
			RateLimit: appconfig.RateLimitConfig{
				RPS: 100,
			},
		},
	}

	// Initialize test environment
	testEnvironment = e2e.NewTestContainer()
	defer testEnvironment.Close()
	testEnvironment.MustInitialize(cfg)

	// Setup routes
	router := rest.SetupRouter(testEnvironment.Container)
	router.GET("echo", func(c *gin.Context) {
		c.Header("Referer", "http://"+c.Request.Host+c.Request.URL.Path)
		c.Status(200)
	})

	// Setup validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator3.PastDateValidator)
		_ = v.RegisterValidation("zxcvbn", validator3.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator3.UsernameValidator)
	}

	// Create server
	server = httptest.NewServer(router)
	defer server.Close()

	fixtures = e2e.MustNewFixtures(testEnvironment.Container.DB, pwd+"/fixtures/users.yml")

	os.Exit(m.Run())
}

func Test_LoginHandler_Login(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/login"

	t.Run("Bad body", func(t *testing.T) {
		bodies := []*dto.LoginInput{
			nil,
			{
				Login:    "user",
				Password: "",
			},
			{
				Login:    "",
				Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
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
		reqBody := dto.LoginInput{
			Login:    "not_exist_user@example.com",
			Password: "7676393c-0d28-4f68-807d-12aa6b88c039",
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

	t.Run("Wrong password", func(t *testing.T) {
		// Send request
		reqBody := dto.LoginInput{
			Login:    "user_for_login@example.com",
			Password: "wrong_password",
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

	t.Run("Blocked user", func(t *testing.T) {
		// Send request
		reqBody := dto.LoginInput{
			Login:    "user_for_login_blocked@example.com",
			Password: "user_for_login_blocked",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Success with email", func(t *testing.T) {
		// Send request
		reqBody := dto.LoginInput{
			Login:    "user_for_login@example.com",
			Password: "user_for_login",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)

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
		setCookie := resp.Header.Get("Set-Cookie")
		assert.NotEmpty(t, setCookie)
	})

	t.Run("Success with username", func(t *testing.T) {
		// Send request
		reqBody := dto.LoginInput{
			Login:    "user_for_login",
			Password: "user_for_login",
		}

		reqBodyReader, _ := e2e.NewJSONReader(reqBody)

		req, _ := http.NewRequest("POST", url, reqBodyReader)

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
		setCookie := resp.Header.Get("Set-Cookie")
		assert.NotEmpty(t, setCookie)
		assert.Contains(t, setCookie, "RefreshToken=")
	})
}

func Test_LoginHandler_RefreshTokens(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/refresh"

	t.Run("Not cookie", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("Invalid refresh token", func(t *testing.T) {
		jwtSecret := testEnvironment.Container.Config.Security.JWT.Secret
		refreshTokens := []string{
			"invalid_refresh_token",
			generateRefreshToken(jwtSecret, jwt.SigningMethodHS256, nil),
			generateRefreshToken(jwtSecret, jwt.SigningMethodHS256, jwt.MapClaims{}),
			generateRefreshToken(jwtSecret, jwt.SigningMethodHS256, jwt.MapClaims{
				"iss": "mandarine",
			}),
			generateRefreshToken(jwtSecret, jwt.SigningMethodHS256, jwt.MapClaims{
				"iss": "mandarine",
				"sub": uuid.New().String(),
			}),
			generateRefreshToken(jwtSecret, jwt.SigningMethodHS256, jwt.MapClaims{
				"iss": "mandarine",
				"sub": uuid.New().String(),
				"iat": time.Now().Unix(),
			}),
			generateRefreshToken(jwtSecret, jwt.SigningMethodES256, jwt.MapClaims{
				"iss": "mandarine",
				"sub": uuid.New().String(),
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour).Unix(),
			}),
		}

		for i, refreshToken := range refreshTokens {
			t.Run(fmt.Sprintf("Invalid refresh token %d", i), func(t *testing.T) {
				// Send request
				req, _ := http.NewRequest("GET", url, nil)
				req.AddCookie(&http.Cookie{
					Name:  "RefreshToken",
					Value: refreshToken,
				})

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
		// Create refresh token
		_, refreshToken, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model.UserEntity{
			ID: uuid.New(),
		})

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.AddCookie(&http.Cookie{
			Name:  "RefreshToken",
			Value: refreshToken,
		})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Check response
		var body dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
	})

	t.Run("User is blocked", func(t *testing.T) {
		// Create refresh token
		_, refreshToken, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model.UserEntity{
			ID:       uuid.MustParse("dded243b-a58f-47ba-9007-2fc41cf950c6"),
			Username: "user_for_refresh_blocked",
			Email:    "user_for_refresh_blocked@example.com",
			Role:     model.RoleEntity{Name: model.RoleUser},
		})

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.AddCookie(&http.Cookie{
			Name:  "RefreshToken",
			Value: refreshToken,
		})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// Create refresh token
		_, refreshToken, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model.UserEntity{
			ID:       uuid.MustParse("d7163725-df27-45de-ae9c-0b860c9ffd17"),
			Username: "user_for_refresh",
			Email:    "user_for_refresh@example.com",
			Role:     model.RoleEntity{Name: model.RoleUser},
		})

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.AddCookie(&http.Cookie{
			Name:  "RefreshToken",
			Value: refreshToken,
		})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check response
		var body dto.JwtTokensOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.NotEmpty(t, body.AccessToken)
		assert.Empty(t, body.RefreshToken)

		// Check header
		setCookie := resp.Header.Get("Set-Cookie")
		assert.NotEmpty(t, setCookie)
	})
}

func generateRefreshToken(secret string, method jwt.SigningMethod, claims jwt.Claims) string {
	refreshToken := jwt.NewWithClaims(method, claims)
	refreshTokenSigned, _ := refreshToken.SignedString([]byte(secret))
	return refreshTokenSigned
}
