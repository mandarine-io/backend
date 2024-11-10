package profile_e2e_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/helper/ref"
	"github.com/mandarine-io/Backend/internal/helper/security"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	http2 "github.com/mandarine-io/Backend/internal/transport/http"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	validator2 "github.com/mandarine-io/Backend/pkg/transport/http/validator"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"os"
	"strconv"
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
			Path:     pwd + "/../../../../../locales",
			Language: "ru",
		},
		Template: appconfig.TemplateConfig{
			Path: pwd + "/../../../../../templates",
		},
		Migrations: appconfig.MigrationConfig{
			Path: pwd + "/../../../../../migrations",
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

	fixtures = e2e.MustNewFixtures(
		testEnvironment.Container.DB,
		pwd+"/fixtures/users.yml",
		pwd+"/fixtures/master_profiles.yml",
		//pwd+"/fixtures/master_profile_vectors.yml",
	)

	os.Exit(m.Run())
}

func Test_MasterProfileHandler_CreateMasterProfile(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	// Create access token
	userEntity := &model.UserEntity{
		ID:        uuid.MustParse("9c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
		Username:  "user_for_create_master_profile",
		Email:     "user_for_create_master_profile@example.com",
		Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
		Role:      model.RoleEntity{Name: model.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	url := server.URL + "/v0/masters/profile"

	t.Run("Unauthorized", func(t *testing.T) {
		// Create request
		req, err := http.NewRequest("POST", url, nil)
		require.NoError(t, err)

		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Check response
		var resBody dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &resBody)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Invalid body", func(t *testing.T) {
		bodies := []*dto.CreateMasterProfileInput{
			nil,

			// Bad display name
			{
				DisplayName: "",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0,0",
				AvatarID:    nil,
			},

			// Bad job
			{
				DisplayName: "display name",
				Job:         "",
				Description: nil,
				Address:     nil,
				Point:       "0,0",
				AvatarID:    nil,
			},

			// Bad point
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "",
				AvatarID:    nil,
			},
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0",
				AvatarID:    nil,
			},
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0,0,0",
				AvatarID:    nil,
			},
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "point",
				AvatarID:    nil,
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				req, _ := http.NewRequest("POST", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}

	})

	t.Run("Duplicate master profile", func(t *testing.T) {
		// Create access token
		anotherUserEntity := &model.UserEntity{
			ID:        uuid.MustParse("8c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
			Username:  "user_for_create_master_profile_exists",
			Email:     "user_for_create_master_profile_exists@example.com",
			Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, anotherUserEntity)

		reqBodyReader, _ := e2e.NewJSONReader(&dto.CreateMasterProfileInput{
			DisplayName: "display name",
			Job:         "job",
			Description: nil,
			Address:     nil,
			Point:       "0,0",
			AvatarID:    nil,
		})

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		reqBodyReader, _ := e2e.NewJSONReader(&dto.CreateMasterProfileInput{
			DisplayName: "display name",
			Job:         "job",
			Description: nil,
			Address:     nil,
			Point:       "0,0",
			AvatarID:    nil,
		})

		req, _ := http.NewRequest("POST", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var body dto.OwnMasterProfileOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)

		assert.Equal(t, "display name", body.DisplayName)
		assert.Equal(t, "job", body.Job)
		assert.Nil(t, body.Description)
		assert.Nil(t, body.Address)
		assert.Equal(t, 0.0, body.Point.Longitude)
		assert.Equal(t, 0.0, body.Point.Latitude)
		assert.Nil(t, body.AvatarID)
		assert.True(t, body.IsEnabled)
	})
}

func Test_MasterProfileHandler_UpdateMasterProfile(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	// Create access token
	userEntity := &model.UserEntity{
		ID:        uuid.MustParse("7c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
		Username:  "user_for_update_master_profile",
		Email:     "user_for_update_master_profile@example.com",
		Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
		Role:      model.RoleEntity{Name: model.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	url := server.URL + "/v0/masters/profile"

	t.Run("Unauthorized", func(t *testing.T) {
		// Create request
		req, err := http.NewRequest("PATCH", url, nil)
		require.NoError(t, err)

		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Check response
		var resBody dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &resBody)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("PATCH", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Invalid body", func(t *testing.T) {
		bodies := []*dto.UpdateMasterProfileInput{
			nil,

			// Bad display name
			{
				DisplayName: "",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0,0",
				AvatarID:    nil,
				IsEnabled:   ref.SafeRef(true),
			},

			// Bad job
			{
				DisplayName: "display name",
				Job:         "",
				Description: nil,
				Address:     nil,
				Point:       "0,0",
				AvatarID:    nil,
				IsEnabled:   ref.SafeRef(true),
			},

			// Bad point
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "",
				AvatarID:    nil,
				IsEnabled:   ref.SafeRef(true),
			},
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0",
				AvatarID:    nil,
				IsEnabled:   ref.SafeRef(true),
			},
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0,0,0",
				AvatarID:    nil,
				IsEnabled:   ref.SafeRef(true),
			},
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "point",
				AvatarID:    nil,
				IsEnabled:   ref.SafeRef(true),
			},

			// Bas is enabled
			{
				DisplayName: "display name",
				Job:         "job",
				Description: nil,
				Address:     nil,
				Point:       "0,0",
				AvatarID:    nil,
				IsEnabled:   nil,
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				var reqBodyReader io.Reader = nil
				if body != nil {
					reqBodyReader, _ = e2e.NewJSONReader(body)
				}

				req, _ := http.NewRequest("PATCH", url, reqBodyReader)
				req.Header.Set("Authorization", "Bearer "+accessToken)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var body dto2.ErrorResponse
				err = e2e.ReadResponseBody(resp, &body)
				assert.NoError(t, err)
			})
		}

	})

	t.Run("Master profile not exists", func(t *testing.T) {
		// Create access token
		anotherUserEntity := &model.UserEntity{
			ID:        uuid.MustParse("6c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
			Username:  "user_for_update_master_profile_not_exists",
			Email:     "user_for_update_master_profile_not_exists@example.com",
			Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, anotherUserEntity)

		reqBodyReader, _ := e2e.NewJSONReader(&dto.UpdateMasterProfileInput{
			DisplayName: "display name",
			Job:         "job",
			Description: nil,
			Address:     nil,
			Point:       "0,0",
			AvatarID:    nil,
			IsEnabled:   ref.SafeRef(true),
		})

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		reqBodyReader, _ := e2e.NewJSONReader(&dto.UpdateMasterProfileInput{
			DisplayName: "display name",
			Job:         "job",
			Description: nil,
			Address:     nil,
			Point:       "1,1",
			AvatarID:    nil,
			IsEnabled:   ref.SafeRef(false),
		})

		req, _ := http.NewRequest("PATCH", url, reqBodyReader)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.OwnMasterProfileOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)

		assert.Equal(t, "display name", body.DisplayName)
		assert.Equal(t, "job", body.Job)
		assert.Nil(t, body.Description)
		assert.Nil(t, body.Address)
		assert.Equal(t, 1.0, body.Point.Longitude)
		assert.Equal(t, 1.0, body.Point.Latitude)
		assert.Nil(t, body.AvatarID)
		assert.False(t, body.IsEnabled)
	})
}

func Test_MasterProfileHandler_GetMasterProfile(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	// Create access token
	userEntity := &model.UserEntity{
		ID:        uuid.MustParse("5c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
		Username:  "user_for_get_master_profile",
		Email:     "user_for_get_master_profile@example.com",
		Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
		Role:      model.RoleEntity{Name: model.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	url := server.URL + "/v0/masters/profile"

	t.Run("Unauthorized", func(t *testing.T) {
		// Create request
		req, err := http.NewRequest("GET", url+"/username", nil)
		require.NoError(t, err)

		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Check response
		var resBody dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &resBody)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url+"/username", nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url+"/username", nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Get own master profile", func(t *testing.T) {
		t.Run("Master profile not exists", func(t *testing.T) {
			// Create access token
			anotherUserEntity := &model.UserEntity{
				ID:        uuid.MustParse("2c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
				Username:  "user_for_get_master_profile_not_exists",
				Email:     "user_for_get_master_profile_not_exists@example.com",
				Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
				Role:      model.RoleEntity{Name: model.RoleUser},
				IsEnabled: true,
				DeletedAt: nil,
			}
			anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, anotherUserEntity)

			req, _ := http.NewRequest("GET", url+"/user_for_get_master_profile_not_exists", nil)
			req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Success", func(t *testing.T) {
			req, _ := http.NewRequest("GET", url+"/user_for_get_master_profile", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var body dto.OwnMasterProfileOutput
			err = e2e.ReadResponseBody(resp, &body)
			assert.NoError(t, err)

			assert.Equal(t, "name3", body.DisplayName)
			assert.Equal(t, "job3", body.Job)
			assert.NotNil(t, body.Description)
			assert.Equal(t, "description3", *body.Description)
			assert.NotNil(t, body.Address)
			assert.Equal(t, "address3", *body.Address)
			assert.Equal(t, 3.0, body.Point.Longitude)
			assert.Equal(t, 3.0, body.Point.Latitude)
			assert.Nil(t, body.AvatarID)
			assert.True(t, body.IsEnabled)
		})
	})

	t.Run("Get another master profile", func(t *testing.T) {
		t.Run("Master profile disabled", func(t *testing.T) {
			req, _ := http.NewRequest("GET", url+"/user_for_get_master_profile_disabled", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("Success", func(t *testing.T) {
			req, _ := http.NewRequest("GET", url+"/user_for_get_master_profile_enabled", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var body dto.MasterProfileOutput
			err = e2e.ReadResponseBody(resp, &body)
			assert.NoError(t, err)

			assert.Equal(t, "name4", body.DisplayName)
			assert.Equal(t, "job4", body.Job)
			assert.NotNil(t, body.Description)
			assert.Equal(t, "description4", *body.Description)
			assert.NotNil(t, body.Address)
			assert.Equal(t, "address4", *body.Address)
			assert.Equal(t, 4.0, body.Point.Longitude)
			assert.Equal(t, 4.0, body.Point.Latitude)
			assert.Nil(t, body.AvatarID)
		})
	})
}

func Test_MasterProfileHandler_FindMasterProfiles(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	// Create access token
	userEntity := &model.UserEntity{
		ID:        uuid.MustParse("5c4778b8-8af3-41fc-bf06-6b5bfeddbbad"),
		Username:  "user_for_find_master_profiles",
		Email:     "user_for_find_master_profiles@example.com",
		Password:  "$2a$12$rKquaEVs3ltdaMJj1yCaReka5T1TMm61AUYfiK3VsQJCOvaJLiOk2",
		Role:      model.RoleEntity{Name: model.RoleUser},
		IsEnabled: true,
		DeletedAt: nil,
	}
	accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, userEntity)

	url := server.URL + "/v0/masters/profile"

	t.Run("Unauthorized", func(t *testing.T) {
		// Create request
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := server.Client().Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Check response
		var resBody dto2.ErrorResponse
		err = e2e.ReadResponseBody(resp, &resBody)
		assert.NoError(t, err)
	})

	t.Run("Banned user", func(t *testing.T) {
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: false,
			DeletedAt: nil,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Deleted user", func(t *testing.T) {
		deletedTime := time.Now().UTC()
		anotherUserEntity := model.UserEntity{
			ID:        uuid.New(),
			Username:  "user",
			Email:     "user@example.com",
			Role:      model.RoleEntity{Name: model.RoleUser},
			IsEnabled: true,
			DeletedAt: &deletedTime,
		}
		anotherAccessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &anotherUserEntity)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+anotherAccessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Invalid data", func(t *testing.T) {
		bodies := []dto.FindMasterProfilesInput{
			// Bad point
			{
				FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
					Point: ref.SafeRef("invalid"),
				},
			},
			{
				FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
					Point: ref.SafeRef("invalid,0"),
				},
			},
			{
				FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
					Point: ref.SafeRef("0,0,0"),
				},
			},

			// Bad radius
			{
				FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
					Radius: ref.SafeRef("invalid"),
				},
			},
			{
				FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
					Radius: ref.SafeRef("0,0"),
				},
			},

			// Bad sort field
			{
				SortInput: &dto.SortInput{
					Field: "",
					Order: "desc",
				},
			},

			// Bad sort order
			{
				SortInput: &dto.SortInput{
					Field: "displayName",
					Order: "",
				},
			},

			// Bad pagination page
			{
				PaginationInput: &dto.PaginationInput{
					Page:     -1,
					PageSize: 10,
				},
			},

			// Bad pagination page size
			{
				PaginationInput: &dto.PaginationInput{
					Page:     0,
					PageSize: 0,
				},
			},
			{
				PaginationInput: &dto.PaginationInput{
					Page:     0,
					PageSize: 1000,
				},
			},
		}

		for i, body := range bodies {
			t.Run(fmt.Sprintf("Bad body %d", i), func(t *testing.T) {
				// Send request
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.URL.RawQuery = mapFindMastersInputToQueryString(body)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		}
	})

	t.Run("SortInput", func(t *testing.T) {
		t.Run("Unavailable field", func(t *testing.T) {
			// Send request
			input := dto.FindMasterProfilesInput{
				SortInput: &dto.SortInput{
					Field: "unavailable",
					Order: "desc",
				},
			}

			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			req.URL.RawQuery = mapFindMastersInputToQueryString(input)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Success", func(t *testing.T) {
			inputs := []dto.FindMasterProfilesInput{
				{
					SortInput: &dto.SortInput{
						Field: "display_name",
						Order: "asc",
					},
				},
				{
					SortInput: &dto.SortInput{
						Field: "display_name",
						Order: "desc",
					},
				},
				{
					SortInput: &dto.SortInput{
						Field: "job",
						Order: "asc",
					},
				},
				{
					SortInput: &dto.SortInput{
						Field: "job",
						Order: "desc",
					},
				},
				{
					SortInput: &dto.SortInput{
						Field: "address",
						Order: "asc",
					},
				},
				{
					SortInput: &dto.SortInput{
						Field: "address",
						Order: "desc",
					},
				},
				{
					FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
						Point:  ref.SafeRef("0,0"),
						Radius: ref.SafeRef("1000000000"),
					},
					SortInput: &dto.SortInput{
						Field: "point",
						Order: "asc",
					},
				},
				{
					FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
						Point:  ref.SafeRef("0,0"),
						Radius: ref.SafeRef("1000000000"),
					},
					SortInput: &dto.SortInput{
						Field: "point",
						Order: "desc",
					},
				},
			}

			for _, input := range inputs {
				t.Run(fmt.Sprintf("Field=%s;Order=%s", input.SortInput.Field, input.SortInput.Order), func(t *testing.T) {
					// Send request
					req, _ := http.NewRequest("GET", url, nil)
					req.Header.Set("Authorization", "Bearer "+accessToken)
					req.URL.RawQuery = mapFindMastersInputToQueryString(input)

					resp, err := server.Client().Do(req)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)

					// Check response
					var resBody dto.MasterProfilesOutput
					err = e2e.ReadResponseBody(resp, &resBody)
					assert.NoError(t, err)

					assert.Greater(t, len(resBody.Data), 2)
					assert.Greater(t, resBody.Count, 2)

					if input.SortInput.Order == "asc" {
						assert.Less(t, resBody.Data[0].Point.Latitude, resBody.Data[1].Point.Latitude)
					} else {
						assert.Greater(t, resBody.Data[0].Point.Latitude, resBody.Data[1].Point.Latitude)
					}
				})
			}
		})
	})

	t.Run("PaginationInput", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// Send request
			input := dto.FindMasterProfilesInput{
				PaginationInput: &dto.PaginationInput{
					Page:     0,
					PageSize: 2,
				},
			}

			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.URL.RawQuery = mapFindMastersInputToQueryString(input)

			resp, err := server.Client().Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Check response
			var resBody dto.MasterProfilesOutput
			err = e2e.ReadResponseBody(resp, &resBody)
			assert.NoError(t, err)

			assert.Less(t, len(resBody.Data), 3)
		})
	})

	t.Run("Filter", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			inputs := []dto.FindMasterProfilesInput{
				// Point
				{
					FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
						Point:  ref.SafeRef("1,1"),
						Radius: ref.SafeRef("1"),
					},
				},

				// Display name
				{
					FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
						DisplayName: ref.SafeRef("name1"),
					},
				},

				// Job
				{
					FindMasterProfilesFilterInput: &dto.FindMasterProfilesFilterInput{
						Job: ref.SafeRef("job1"),
					},
				},
			}

			for _, input := range inputs {
				inputRaw, _ := json.Marshal(input)

				t.Run(string(inputRaw), func(t *testing.T) {
					req, _ := http.NewRequest("GET", url, nil)
					req.Header.Set("Authorization", "Bearer "+accessToken)
					req.URL.RawQuery = mapFindMastersInputToQueryString(input)

					resp, err := server.Client().Do(req)
					assert.NoError(t, err)
					assert.Equal(t, http.StatusOK, resp.StatusCode)

					// Check response
					var resBody dto.MasterProfilesOutput
					err = e2e.ReadResponseBody(resp, &resBody)
					assert.NoError(t, err)

					assert.Equal(t, len(resBody.Data), 1)
					assert.Equal(t, resBody.Count, 1)
				})
			}
		})
	})
}

func mapFindMastersInputToQueryString(input dto.FindMasterProfilesInput) string {
	q := neturl.Values{}

	if input.FindMasterProfilesFilterInput != nil {
		if input.FindMasterProfilesFilterInput.Point != nil {
			q.Add("point", *input.FindMasterProfilesFilterInput.Point)
		}
		if input.FindMasterProfilesFilterInput.Radius != nil {
			q.Add("radius", *input.FindMasterProfilesFilterInput.Radius)
		}
		if input.FindMasterProfilesFilterInput.Job != nil {
			q.Add("job", *input.FindMasterProfilesFilterInput.Job)
		}
		if input.FindMasterProfilesFilterInput.DisplayName != nil {
			q.Add("displayName", *input.FindMasterProfilesFilterInput.DisplayName)
		}
	}

	if input.SortInput != nil {
		q.Add("field", input.SortInput.Field)
		q.Add("order", input.SortInput.Order)
	}

	if input.PaginationInput != nil {
		q.Add("page", strconv.Itoa(input.PaginationInput.Page))
		q.Add("pageSize", strconv.Itoa(input.PaginationInput.PageSize))
	}

	return q.Encode()
}
