package ws_e2e_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	appconfig "github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/internal/api/helper/security"
	"github.com/mandarine-io/Backend/internal/api/persistence/model"
	"github.com/mandarine-io/Backend/internal/api/rest"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			Level: "debug",
			Console: appconfig.ConsoleLoggerConfig{
				Enable:   true,
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
		Websocket: appconfig.WebsocketConfig{
			PoolSize: 1,
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

func Test_Connect(t *testing.T) {
	url := fmt.Sprintf("%s/ws", strings.Replace(server.URL, "http", "ws", 1))

	userEntity := &model.UserEntity{
		ID:        uuid.MustParse("fa4d4574-7e1a-4a1d-b262-760423c3743d"),
		Username:  "user",
		Email:     "user@example.com",
		Password:  "$2a$12$3VuuVLeH1Psvc8umGMjyHusALuKz8OsmqiGE79JQXdeaKj5ih1u8e",
		Role:      model.RoleEntity{Name: model.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	t.Run("Success connection", func(t *testing.T) {
		ws, _, err := websocket.DefaultDialer.Dial(url, http.Header{"Authorization": {"Bearer " + accessToken}})
		require.NoError(t, err)

		err = ws.Close()
		require.NoError(t, err)
	})

	t.Run("Pool is full", func(t *testing.T) {
		anotherUserEntity := &model.UserEntity{
			ID:        uuid.MustParse("7d26639d-7082-45df-b1d0-77f663f5003a"),
			Username:  "user1",
			Email:     "user1@example.com",
			Password:  "$2a$12$4XWfvkfvvLxLlLyPQ9CA7eNhkUIFSj7sF3768lAMJi9G2kl4XjGve",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, anotherUserEntity)

		// First connection
		ws, resp, err := websocket.DefaultDialer.Dial(url, http.Header{"Authorization": {"Bearer " + accessToken}})
		require.NoError(t, err)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

		// Second connection
		_, resp, err = websocket.DefaultDialer.Dial(url, http.Header{"Authorization": {"Bearer " + anotherAccessToken}})
		require.Error(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

		err = ws.Close()
		require.NoError(t, err)
	})
}
