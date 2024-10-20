package health_test

import (
	"github.com/stretchr/testify/assert"
	appconfig "mandarine/internal/api/config"
	"mandarine/internal/api/rest"
	"mandarine/internal/api/rest/handler/health"
	"mandarine/tests/e2e"
	"net/http"
	"net/http/httptest"
	"os"
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

	testEnvironment = e2e.NewTestContainer()
	defer testEnvironment.Close()

	testEnvironment.MustInitialize(cfg)
	router := rest.SetupRouter(testEnvironment.Container)
	server = httptest.NewServer(router)
	defer server.Close()

	os.Exit(m.Run())
}

func Test_HealthCheck(t *testing.T) {
	url := server.URL + "/health"

	resp, err := server.Client().Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//// Check response
	var body []health.Response
	err = e2e.ReadResponseBody(resp, &body)
	assert.NoError(t, err)
	for _, v := range body {
		assert.True(t, v.Pass)
	}
}
