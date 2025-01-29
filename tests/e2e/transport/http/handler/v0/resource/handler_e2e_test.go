package resource_e2e_test

//import (
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	e2e2 "github.com/mandarine-io/backend/.github/tests/e2e"
//	"github.com/mandarine-io/backend/config"
//	appconfig "github.com/mandarine-io/backend/internal/config"
//	model3 "github.com/mandarine-io/backend/internal/infrastructure/s3"
//	model2 "github.com/mandarine-io/backend/internal/persistence/entity"
//	http2 "github.com/mandarine-io/backend/internal/transport/http"
//	"github.com/mandarine-io/backend/internal/util/security"
//	"github.com/mandarine-io/backend/pkg/model"
//	model2 "github.com/mandarine-io/backend/pkg/model"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"mime/multipart"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"path/filepath"
//	"strings"
//	"testing"
//	"time"
//)
//
//var (
//	testEnvironment *e2e2.TestEnvironment
//	server          *httptest.Server
//
//	ctx = context.Background()
//)
//
//func TestMain(m *testing.M) {
//	pwd, err := os.Getwd()
//	if err != nil {
//		panic(err)
//	}
//
//	cfg := &config.Config{
//		Server: config.ServerConfig{
//			Name:           "mandarine_test",
//			Mode:           "test",
//			ExternalURL:    "http://localhost:8081",
//			Port:           8081,
//			Version:        "0.0.0",
//			RPS:            100,
//			MaxRequestSize: 524288000,
//		},
//		Database: config.PostgresDatabaseConfig{
//			Type: "postgres",
//			Postgres: &config.PostgresDatabaseConfig{
//				Username: "mandarine",
//				Password: "password",
//				DBName:   "mandarine_test",
//			},
//		},
//		Cache: config.RedisCacheConfig{
//			TTL:  120,
//			Type: "redis",
//			Redis: &config.RedisCacheConfig{
//				Username: "default",
//				Password: "password",
//				DBIndex:  0,
//			},
//		},
//		PubSub: config.RedisPubSubConfig{
//			Type: "redis",
//			Redis: &config.RedisPubSubConfig{
//				Username: "default",
//				Password: "password",
//				DBIndex:  0,
//			},
//		},
//		S3: config.MinIOS3Config{
//			Type: "minio",
//			Minio: &appconfig.MinioS3Config{
//				AccessKey: "admin",
//				SecretKey: "Password_10",
//				Bucket:    "mandarine-test",
//			},
//		},
//		SMTP: config.SMTPConfig{
//			Host:     "127.0.0.1",
//			Port:     25,
//			Username: "admin",
//			Password: "password",
//			From:     "Mandarine <admin@localhost>",
//			SSL:      false,
//		},
//		Locale: config.LocaleConfig{
//			Path:     pwd + "/../../../../locales",
//			Language: "ru",
//		},
//		Template: config.TemplateConfig{
//			Path: pwd + "/../../../../templates",
//		},
//		Migrations: config.MigrationConfig{
//			Path: pwd + "/../../../../migrations",
//		},
//		Logger: config.LoggerConfig{
//			Level: "debug",
//			Console: config.ConsoleLoggerConfig{
//				Enable:   true,
//				Encoding: "text",
//			},
//			File: config.FileLoggerConfig{
//				Enable: false,
//			},
//		},
//		OAuthProviders: map[string]config.OauthProviderItemConfig{
//			"google": {
//				ClientID:     "",
//				ClientSecret: "",
//			},
//			"yandex": {
//				ClientID:     "",
//				ClientSecret: "",
//			},
//			"mailru": {
//				ClientID:     "",
//				ClientSecret: "",
//			},
//		},
//		Security: config.SecurityConfig{
//			JWT: config.JWTConfig{
//				Secret:          "",
//				AccessTokenTTL:  3600,
//				RefreshTokenTTL: 86400,
//			},
//			OTP: config.OTPConfig{
//				Length: 6,
//				TTL:    300,
//			},
//		},
//		Websocket: config.WebsocketConfig{
//			PoolSize: 1024,
//		},
//	}
//
//	testEnvironment = e2e2.NewTestContainer()
//	defer testEnvironment.Close()
//
//	testEnvironment.MustInitialize(cfg)
//	router := http2.SetupRouter(testEnvironment.Container)
//	server = httptest.NewServer(router)
//	defer server.Close()
//
//	os.Exit(m.Run())
//}
//
//func Test_ResourceHandler_UploadResource(t *testing.T) {
//	// Create access token
//	userEntity := &model2.UserEntity{
//		ID:        uuid.MustParse("a02fc7e1-c19a-4c1a-b66e-29fed1ed452f"),
//		Username:  "user1",
//		Email:     "user1@example.com",
//		Password:  "$2a$12$4XWfvkfvvLxLlLyPQ9CA7eNhkUIFSj7sF3768lAMJi9G2kl4XjGve",
//		Role:      model2.RoleEntity{Name: model2.RoleUser},
//		IsEnabled: true,
//		DeletedAt: nil,
//	}
//	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)
//
//	// Create temp file
//	file, err := os.CreateTemp("", "test")
//	require.Nil(t, err)
//	defer os.Remove(file.Name())
//
//	_, _ = file.WriteString(strings.Repeat("a", 1024))
//	_, _ = file.Seek(0, 0)
//
//	url := server.URL + "/v0/resources/one"
//
//	t.Run(
//		"Unauthorized", func(t *testing.T) {
//			// Create request
//			req, err := http.NewRequest("POST", url, nil)
//			require.NoError(t, err)
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Banned user", func(t *testing.T) {
//			anotherUserEntity := model2.UserEntity{
//				ID:        uuid.New(),
//				Username:  "user",
//				Email:     "user@example.com",
//				Role:      model2.RoleEntity{Name: model2.RoleUser},
//				IsEnabled: false,
//				DeletedAt: nil,
//			}
//			anotherAccessToken, _, _ := security.GenerateTokens(
//				testEnvironment.Container.Config.Security.JWT,
//				&anotherUserEntity,
//			)
//			req, _ := http.NewRequest("POST", url, nil)
//			req.Header.Set("Authorization", "Bearer "+anotherAccessToken)
//
//			resp, err := server.Client().Do(req)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
//		},
//	)
//
//	t.Run(
//		"Deleted user", func(t *testing.T) {
//			deletedTime := time.Now().UTC()
//			anotherUserEntity := model2.UserEntity{
//				ID:        uuid.New(),
//				Username:  "user",
//				Email:     "user@example.com",
//				Role:      model2.RoleEntity{Name: model2.RoleUser},
//				IsEnabled: true,
//				DeletedAt: &deletedTime,
//			}
//			anotherAccessToken, _, _ := security.GenerateTokens(
//				testEnvironment.Container.Config.Security.JWT,
//				&anotherUserEntity,
//			)
//			req, _ := http.NewRequest("POST", url, nil)
//			req.Header.Set("Authorization", "Bearer "+anotherAccessToken)
//
//			resp, err := server.Client().Do(req)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
//		},
//	)
//
//	t.Run(
//		"Body is empty", func(t *testing.T) {
//			// Create request
//			req, err := http.NewRequest("POST", url, nil)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"File not uploaded", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			_, err := mw.CreateFormFile("redis", "file.txt")
//			require.NoError(t, err)
//
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", mw.FormDataContentType())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusCreated, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Invalid field name", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			fw, err := mw.CreateFormFile("invalid_resource", file.Name())
//			require.NoError(t, err)
//
//			_, err = io.Copy(fw, file)
//			require.NoError(t, err)
//
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", gin.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Success", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			fw, err := mw.CreateFormFile("redis", file.Name())
//			require.NoError(t, err)
//
//			_, err = io.Copy(fw, file)
//			require.NoError(t, err)
//
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", gin.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusCreated, resp.StatusCode)
//
//			// Check response
//			var resBody model.UploadResourceOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//			assert.True(t, strings.HasSuffix(resBody.ObjectID, filepath.Base(file.Name())))
//		},
//	)
//}
//
//func Test_ResourceHandler_UploadResources(t *testing.T) {
//	// Create access token
//	userEntity := &model2.UserEntity{
//		ID:        uuid.MustParse("a02fc7e1-c19a-4c1a-b66e-29fed1ed452f"),
//		Username:  "user1",
//		Email:     "user1@example.com",
//		Password:  "$2a$12$4XWfvkfvvLxLlLyPQ9CA7eNhkUIFSj7sF3768lAMJi9G2kl4XjGve",
//		Role:      model2.RoleEntity{Name: model2.RoleUser},
//		IsEnabled: true,
//		DeletedAt: nil,
//	}
//	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)
//
//	// Create temp files
//	var err error
//	files := make([]*os.File, 3)
//	defer func() {
//		for _, file := range files {
//			if file != nil {
//				os.Remove(file.Name())
//			}
//		}
//	}()
//	for i := 0; i < 3; i++ {
//		files[i], err = os.CreateTemp("", "test")
//		require.Nil(t, err)
//
//		_, _ = files[i].WriteString(strings.Repeat("a", 1024))
//		_, _ = files[i].Seek(0, 0)
//	}
//
//	url := server.URL + "/v0/resources/many"
//
//	t.Run(
//		"Unauthorized", func(t *testing.T) {
//			// Create request
//			req, err := http.NewRequest("POST", url, nil)
//			require.NoError(t, err)
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Banned user", func(t *testing.T) {
//			anotherUserEntity := model2.UserEntity{
//				ID:        uuid.New(),
//				Username:  "user",
//				Email:     "user@example.com",
//				Role:      model2.RoleEntity{Name: model2.RoleUser},
//				IsEnabled: false,
//				DeletedAt: nil,
//			}
//			anotherAccessToken, _, _ := security.GenerateTokens(
//				testEnvironment.Container.Config.Security.JWT,
//				&anotherUserEntity,
//			)
//			req, _ := http.NewRequest("POST", url, nil)
//			req.Header.Set("Authorization", "Bearer "+anotherAccessToken)
//
//			resp, err := server.Client().Do(req)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
//		},
//	)
//
//	t.Run(
//		"Deleted user", func(t *testing.T) {
//			deletedTime := time.Now().UTC()
//			anotherUserEntity := model2.UserEntity{
//				ID:        uuid.New(),
//				Username:  "user",
//				Email:     "user@example.com",
//				Role:      model2.RoleEntity{Name: model2.RoleUser},
//				IsEnabled: true,
//				DeletedAt: &deletedTime,
//			}
//			anotherAccessToken, _, _ := security.GenerateTokens(
//				testEnvironment.Container.Config.Security.JWT,
//				&anotherUserEntity,
//			)
//			req, _ := http.NewRequest("POST", url, nil)
//			req.Header.Set("Authorization", "Bearer "+anotherAccessToken)
//
//			resp, err := server.Client().Do(req)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
//		},
//	)
//
//	t.Run(
//		"Body is empty", func(t *testing.T) {
//			// Create request
//			req, err := http.NewRequest("POST", url, nil)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Files not uploaded", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			_, err := mw.CreateFormFile("resources", "file.txt")
//			require.NoError(t, err)
//
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", mw.FormDataContentType())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusCreated, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Invalid field name", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			for _, file := range files {
//				fw, err := mw.CreateFormFile("invalid_resources", file.Name())
//				require.NoError(t, err)
//
//				_, err = io.Copy(fw, file)
//				require.NoError(t, err)
//			}
//
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", gin.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Duplicate files", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			for i := 0; i < 2; i++ {
//				for _, file := range files {
//					fw, err := mw.CreateFormFile("resources", file.Name())
//					require.NoError(t, err)
//
//					_, err = io.Copy(fw, file)
//					require.NoError(t, err)
//				}
//			}
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", gin.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusCreated, resp.StatusCode)
//
//			// Check response
//			var resBody model.UploadResourcesOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//			assert.Equal(t, len(files), resBody.Count)
//			for _, file := range files {
//				fileName := filepath.Base(file.Name())
//				fileResp, ok := resBody.Data[fileName]
//				assert.True(t, ok)
//
//				objectId := fileResp.ObjectID
//				assert.True(t, strings.HasSuffix(objectId, fileName))
//			}
//		},
//	)
//
//	t.Run(
//		"Success", func(t *testing.T) {
//			// Create body
//			var body bytes.Buffer
//
//			mw := multipart.NewWriter(&body)
//			for _, file := range files {
//				fw, err := mw.CreateFormFile("resources", file.Name())
//				require.NoError(t, err)
//
//				_, err = io.Copy(fw, file)
//				require.NoError(t, err)
//			}
//			err = mw.Close()
//			require.NoError(t, err)
//
//			// Create request
//			req, err := http.NewRequest("POST", url, &body)
//			require.NoError(t, err)
//
//			req.Header.Set("Authorization", "Bearer "+accessToken)
//			req.Header.Set("Content-Type", gin.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())
//
//			resp, err := server.Client().Do(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusCreated, resp.StatusCode)
//
//			// Check response
//			var resBody model.UploadResourcesOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//			assert.Equal(t, len(files), resBody.Count)
//			for _, file := range files {
//				fileName := filepath.Base(file.Name())
//				fileResp, ok := resBody.Data[fileName]
//				assert.True(t, ok)
//
//				objectId := fileResp.ObjectID
//				assert.True(t, strings.HasSuffix(objectId, fileName))
//			}
//		},
//	)
//}
//
//func Test_ResourceHandler_DownloadResource(t *testing.T) {
//	// Create temp file
//	file, err := os.CreateTemp("", "test")
//	require.Nil(t, err)
//	defer os.Remove(file.Name())
//
//	_, _ = file.WriteString(strings.Repeat("a", 1024))
//	_, _ = file.Seek(0, 0)
//
//	url := server.URL + "/v0/resources"
//
//	t.Run(
//		"Not found", func(t *testing.T) {
//			// Create request
//			req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, "not-found"), nil)
//			require.NoError(t, err)
//
//			resp, err := server.Client().Do(req)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
//
//			// Check response
//			var resBody model2.ErrorOutput
//			err = e2e2.ReadResponseBody(resp, &resBody)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusNotFound, resBody.Status)
//		},
//	)
//
//	t.Run(
//		"Success", func(t *testing.T) {
//			// Upload file
//			objectId := filepath.Base(file.Name())
//			testEnvironment.Container.S3.Client.CreateOne(
//				ctx, &model3.FileData{
//					ID:     objectId,
//					Size:   1024,
//					Reader: file,
//				},
//			)
//
//			// Create request
//			req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", url, objectId), nil)
//			require.NoError(t, err)
//
//			resp, err := server.Client().Do(req)
//			assert.NoError(t, err)
//
//			assert.Equal(t, http.StatusOK, resp.StatusCode)
//
//			// Check response
//			var resBody bytes.Buffer
//			_, err = io.Copy(&resBody, resp.Body)
//			assert.NoError(t, err)
//			assert.Equal(t, 1024, resBody.Len())
//			assert.Equal(t, strings.Repeat("a", 1024), resBody.String())
//		},
//	)
//}
