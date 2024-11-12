package geocoding_e2e_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service/geocoding"
	"github.com/mandarine-io/Backend/internal/helper/ref"
	"github.com/mandarine-io/Backend/internal/helper/security"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	http2 "github.com/mandarine-io/Backend/internal/transport/http"
	geocoding3 "github.com/mandarine-io/Backend/internal/transport/http/handler/v0/geocoding"
	geocoding2 "github.com/mandarine-io/Backend/pkg/geocoding"
	mock3 "github.com/mandarine-io/Backend/pkg/geocoding/mock"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	validator2 "github.com/mandarine-io/Backend/pkg/transport/http/validator"
	"github.com/mandarine-io/Backend/tests/e2e"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	ctx               = context.TODO()
	testEnvironment   *e2e.TestEnvironment
	server            *httptest.Server
	geocodingProvider = new(mock3.ProviderMock)
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
		GeocodingClients: map[string]appconfig.GeocodingClientConfig{},
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

	// Add mock geocoding provider
	geocodingSvc := geocoding.NewService(
		[]geocoding2.Provider{geocodingProvider},
		testEnvironment.Container.Cache.Manager,
	)
	testEnvironment.Container.SVCs.Geocoding = geocodingSvc

	for i, handler := range testEnvironment.Container.Handlers {
		name := ref.GetType(handler)
		if strings.HasSuffix(name, "geocoding.handler") {
			testEnvironment.Container.Handlers[i] = geocoding3.NewHandler(geocodingSvc)
		}
	}

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

	os.Exit(m.Run())
}

func Test_GeocodingHandler_Geocode(t *testing.T) {
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

	url := server.URL + "/v0/geocode/forward"
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

	t.Run("Invalid input", func(t *testing.T) {
		params := []dto.GeocodingInput{
			// Bad address
			{Address: "", Limit: 1},

			// Bad limit
			{Address: "address", Limit: 0},
		}

		for i, param := range params {
			t.Run(fmt.Sprintf("Invalid input %d", i), func(t *testing.T) {
				// Send request
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.URL.RawQuery = mapGeocodingInputToQueryString(param)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		}
	})

	t.Run("Geocoding services unavailable", func(t *testing.T) {
		_, _, _ = testEnvironment.RedisC.Exec(ctx, []string{"FLUSHALL"})

		geocodingProvider.On("GeocodeWithContext", mock.Anything, "address", mock.Anything).
			Return(nil, geocoding2.ErrGeocodeProvidersUnavailable).Once()

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.URL.RawQuery = mapGeocodingInputToQueryString(dto.GeocodingInput{Address: "address", Limit: 1})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("Unexpected geocoding error", func(t *testing.T) {
		_, _, _ = testEnvironment.RedisC.Exec(ctx, []string{"FLUSHALL"})

		geocodingProvider.On("GeocodeWithContext", mock.Anything, "address", mock.Anything).
			Return(nil, errors.New("unexpected error")).Once()

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.URL.RawQuery = mapGeocodingInputToQueryString(dto.GeocodingInput{Address: "address", Limit: 1})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		_, _, _ = testEnvironment.RedisC.Exec(ctx, []string{"FLUSHALL"})

		geocodingProvider.On("GeocodeWithContext", mock.Anything, "address", mock.Anything).
			Return([]*geocoding2.Location{{Lat: 1.0, Lng: 1.0}}, nil).Once()

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.URL.RawQuery = mapGeocodingInputToQueryString(dto.GeocodingInput{Address: "address", Limit: 1})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.GeocodingOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, 1, body.Count)
		assert.Equal(t, 1.0, body.Data[0].Latitude)
		assert.Equal(t, 1.0, body.Data[0].Longitude)
	})
}

func Test_GeocodingHandler_ReverseGeocode(t *testing.T) {
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

	url := server.URL + "/v0/geocode/reverse"
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

	t.Run("Invalid input", func(t *testing.T) {
		params := []dto.ReverseGeocodingInput{
			// Bad point
			{Point: "", Limit: 1},

			// Bad limit
			{Point: "1,1", Limit: 0},
		}

		for i, param := range params {
			t.Run(fmt.Sprintf("Invalid input %d", i), func(t *testing.T) {
				// Send request
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.URL.RawQuery = mapReverseGeocodingInputToQueryString(param)

				resp, err := server.Client().Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		}
	})

	t.Run("Geocoding services unavailable", func(t *testing.T) {
		_, _, _ = testEnvironment.RedisC.Exec(ctx, []string{"FLUSHALL"})

		geocodingProvider.On("ReverseGeocodeWithContext", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, geocoding2.ErrGeocodeProvidersUnavailable).Once()

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.URL.RawQuery = mapReverseGeocodingInputToQueryString(dto.ReverseGeocodingInput{Point: "1,1", Limit: 1})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("Unexpected geocoding error", func(t *testing.T) {
		_, _, _ = testEnvironment.RedisC.Exec(ctx, []string{"FLUSHALL"})

		geocodingProvider.On("ReverseGeocodeWithContext", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("unexpected error")).Once()

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.URL.RawQuery = mapReverseGeocodingInputToQueryString(dto.ReverseGeocodingInput{Point: "1,1", Limit: 1})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		_, _, _ = testEnvironment.RedisC.Exec(ctx, []string{"FLUSHALL"})

		geocodingProvider.On("ReverseGeocodeWithContext", mock.Anything, mock.Anything, mock.Anything).
			Return([]*geocoding2.Address{{FormattedAddress: "address"}}, nil).Once()

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.URL.RawQuery = mapReverseGeocodingInputToQueryString(dto.ReverseGeocodingInput{Point: "1,1", Limit: 1})

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body dto.ReverseGeocodingOutput
		err = e2e.ReadResponseBody(resp, &body)
		assert.NoError(t, err)
		assert.Equal(t, 1, body.Count)
		assert.Equal(t, "address", body.Data[0].FormattedAddress)
	})
}

func mapGeocodingInputToQueryString(input dto.GeocodingInput) string {
	q := neturl.Values{}

	q.Add("address", input.Address)
	if input.Limit == 0 {
		q.Add("limit", strconv.Itoa(input.Limit))
	}

	return q.Encode()
}

func mapReverseGeocodingInputToQueryString(input dto.ReverseGeocodingInput) string {
	q := neturl.Values{}

	q.Add("point", input.Point)
	if input.Limit == 0 {
		q.Add("limit", strconv.Itoa(input.Limit))
	}

	return q.Encode()
}
