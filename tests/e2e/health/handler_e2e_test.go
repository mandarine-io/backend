package health_test

import (
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	http2 "github.com/mandarine-io/Backend/internal/transport/http"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/stretchr/testify/assert"
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

func Test_HealthCheck(t *testing.T) {
	url := server.URL + "/health"

	resp, err := server.Client().Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//// Check response
	var body []dto.HealthOutput
	err = e2e.ReadResponseBody(resp, &body)
	assert.NoError(t, err)
	for _, v := range body {
		assert.True(t, v.Pass)
	}
}
