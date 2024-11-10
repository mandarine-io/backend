package auth_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/helper/security"
	model2 "github.com/mandarine-io/Backend/internal/persistence/model"
	http2 "github.com/mandarine-io/Backend/internal/transport/http"
	"github.com/mandarine-io/Backend/pkg/oauth"
	mock3 "github.com/mandarine-io/Backend/pkg/oauth/mock"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	validator2 "github.com/mandarine-io/Backend/pkg/transport/http/validator"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	testEnvironment *e2e.TestEnvironment
	server          *httptest.Server
	fixtures        *testfixtures.Loader
)

type mailhogMessagesResponse struct {
	Total    int           `json:"total"`
	Start    int           `json:"start"`
	Count    int           `json:"count"`
	Messages []interface{} `json:"messages"`
}

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
			RPS:            100,
			MaxRequestSize: 524288000,
		},
		Database: appconfig.DatabaseConfig{
			Type: "postgres",
			Postgres: &appconfig.PostgresDatabaseConfig{
				Username: "mandarine",
				Password: "password",
				DBName:   "mandarine_test",
			},
		},
		Cache: appconfig.CacheConfig{
			TTL:  120,
			Type: "redis",
			Redis: &appconfig.RedisCacheConfig{
				Username: "default",
				Password: "password",
				DBIndex:  0,
			},
		},
		PubSub: appconfig.PubSubConfig{
			Type: "redis",
			Redis: &appconfig.RedisPubSubConfig{
				Username: "default",
				Password: "password",
				DBIndex:  0,
			},
		},
		S3: appconfig.S3Config{
			Type: "minio",
			Minio: &appconfig.MinioS3Config{
				AccessKey: "admin",
				SecretKey: "Password_10",
				Bucket:    "mandarine-test",
			},
		},
		SMTP: appconfig.SmtpConfig{
			Host:     "127.0.0.1",
			Port:     25,
			Username: "admin",
			Password: "password",
			From:     "Mandarine <admin@localhost>",
			SSL:      false,
		},
		Locale: appconfig.LocaleConfig{
			Path:     pwd + "/../../../../locales",
			Language: "ru",
		},
		Template: appconfig.TemplateConfig{
			Path: pwd + "/../../../../templates",
		},
		Migrations: appconfig.MigrationConfig{
			Path: pwd + "/../../../../migrations",
		},
		Logger: appconfig.LoggerConfig{
			Level: "debug",
			Console: appconfig.ConsoleLoggerConfig{
				Enable:   true,
				Encoding: "text",
			},
			File: appconfig.FileLoggerConfig{
				Enable: false,
			},
		},
		OAuthClients: map[string]appconfig.OauthClientConfig{
			"google": {
				ClientID:     "",
				ClientSecret: "",
			},
			"yandex": {
				ClientID:     "",
				ClientSecret: "",
			},
			"mailru": {
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
		},
		Websocket: appconfig.WebsocketConfig{
			PoolSize: 1024,
		},
	}

	// Initialize test environment
	testEnvironment = e2e.NewTestContainer()
	defer testEnvironment.Close()
	testEnvironment.MustInitialize(cfg)

	// Setup routes
	router := http2.SetupRouter(testEnvironment.Container)
	router.GET("echo", func(c *gin.Context) {
		c.Header("Referer", "http://"+c.Request.Host+c.Request.URL.Path)
		c.Status(200)
	})

	// Setup validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator2.PastDateValidator)
		_ = v.RegisterValidation("zxcvbn", validator2.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator2.UsernameValidator)
		_ = v.RegisterValidation("point", validator2.PointValidator)
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
		_, refreshToken, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model2.UserEntity{
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
		_, refreshToken, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model2.UserEntity{
			ID:       uuid.MustParse("dded243b-a58f-47ba-9007-2fc41cf950c6"),
			Username: "user_for_refresh_blocked",
			Email:    "user_for_refresh_blocked@example.com",
			Role:     model2.RoleEntity{Name: model2.RoleUser},
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
		_, refreshToken, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model2.UserEntity{
			ID:       uuid.MustParse("d7163725-df27-45de-ae9c-0b860c9ffd17"),
			Username: "user_for_refresh",
			Email:    "user_for_refresh@example.com",
			Role:     model2.RoleEntity{Name: model2.RoleUser},
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

func Test_LogoutHandler_Logout(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/logout"

	t.Run("Unauthorized", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// New access token
		accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model2.UserEntity{
			ID:       uuid.MustParse("a83d9587-b01f-4146-8b1f-80f137f53534"),
			Username: "user_for_logout",
			Email:    "user_for_logout@example.com",
			Role:     model2.RoleEntity{Name: model2.RoleUser},
		})

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check repeated request
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err = server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
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
		err := testEnvironment.Container.Cache.Manager.Set(
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
		err := testEnvironment.Container.Cache.Manager.Set(
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
		err := testEnvironment.Container.Cache.Manager.Set(
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

func Test_SocialLogin_SocialLogin(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	prefixUrl := server.URL + "/v0/auth/social"
	oauthProviderMock := testEnvironment.Container.OauthProviders["mock"].(*mock3.ProviderMock)

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
	oauthProviderMock := testEnvironment.Container.OauthProviders["mock"].(*mock3.ProviderMock)

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

func generateRefreshToken(secret string, method jwt.SigningMethod, claims jwt.Claims) string {
	refreshToken := jwt.NewWithClaims(method, claims)
	refreshTokenSigned, _ := refreshToken.SignedString([]byte(secret))
	return refreshTokenSigned
}
