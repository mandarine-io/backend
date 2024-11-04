package common_e2e_test

import (
	"context"
	appconfig "github.com/mandarine-io/Backend/internal/api/config"
	http2 "github.com/mandarine-io/Backend/internal/api/transport/http"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testEnvironment *e2e.TestEnvironment
	server          *httptest.Server
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
			RPS:            5,
			MaxRequestSize: 500,
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

	testEnvironment = e2e.NewTestContainer()
	defer testEnvironment.Close()

	testEnvironment.MustInitialize(cfg)
	router := http2.SetupRouter(testEnvironment.Container)
	server = httptest.NewServer(router)
	defer server.Close()

	os.Exit(m.Run())
}

func Test_CommonHandler_NoMethod(t *testing.T) {
	url := server.URL + "/v0/auth/login"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := server.Client().Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	var body dto2.ErrorResponse
	err = e2e.ReadResponseBody(resp, &body)
	assert.NoError(t, err)
}

func Test_CommonHandler_NoRoute(t *testing.T) {
	url := server.URL + "/not-supported-route"
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := server.Client().Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var body dto2.ErrorResponse
	err = e2e.ReadResponseBody(resp, &body)
	assert.NoError(t, err)
}

func Test_CommonHandler_MaxRequestSize(t *testing.T) {
	url := server.URL + "/v0/auth/login"

	reqBody := map[string]interface{}{
		"test": strings.Repeat("a", 1025),
	}
	reqBodyReader, _ := e2e.NewJSONReader(reqBody)

	req, _ := http.NewRequest("POST", url, reqBodyReader)
	req.Header.Set("Authorization", "application/json")

	resp, err := server.Client().Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func Test_CommonHandler_RateLimiter(t *testing.T) {
	url := server.URL + "/health"
	req, _ := http.NewRequest("GET", url, nil)

	resps := make([]http.Response, 100)
	executor, _ := errgroup.WithContext(context.Background())
	for i := 0; i < 100; i++ {
		executor.Go(func() error {
			resp, err := server.Client().Do(req)
			if err != nil {
				return err
			}

			resps[i] = *resp
			return nil
		})
	}

	err := executor.Wait()
	if err != nil {
		require.NoError(t, err)
	}

	count := 0
	for _, resp := range resps {
		if resp.StatusCode == http.StatusTooManyRequests {
			count++
		}
	}

	assert.True(t, count > 0)
}
