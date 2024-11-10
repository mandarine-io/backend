package ws_e2e_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/helper/security"
	model2 "github.com/mandarine-io/Backend/internal/persistence/model"
	http2 "github.com/mandarine-io/Backend/internal/transport/http"
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
			PoolSize: 1,
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

func Test_Connect(t *testing.T) {
	url := fmt.Sprintf("%s/v0/ws", strings.Replace(server.URL, "http", "ws", 1))

	userEntity := &model2.UserEntity{
		ID:        uuid.MustParse("fa4d4574-7e1a-4a1d-b262-760423c3743d"),
		Username:  "user",
		Email:     "user@example.com",
		Password:  "$2a$12$3VuuVLeH1Psvc8umGMjyHusALuKz8OsmqiGE79JQXdeaKj5ih1u8e",
		Role:      model2.RoleEntity{Name: model2.RoleUser},
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
		anotherUserEntity := &model2.UserEntity{
			ID:        uuid.MustParse("7d26639d-7082-45df-b1d0-77f663f5003a"),
			Username:  "user1",
			Email:     "user1@example.com",
			Password:  "$2a$12$4XWfvkfvvLxLlLyPQ9CA7eNhkUIFSj7sF3768lAMJi9G2kl4XjGve",
			Role:      model2.RoleEntity{Name: model2.RoleUser},
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
