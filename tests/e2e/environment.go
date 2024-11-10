package e2e

import (
	"context"
	"fmt"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/registry"
	"github.com/mandarine-io/Backend/pkg/logging"
	mock3 "github.com/mandarine-io/Backend/pkg/oauth/mock"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
	"sync"
)

var (
	ctx                  = context.Background()
	te  *TestEnvironment = nil
)

type TestEnvironment struct {
	PostgresC testcontainers.Container
	RedisC    testcontainers.Container
	MinioC    testcontainers.Container
	SmtpC     testcontainers.Container
	Container *registry.Container

	mu        sync.Mutex
	initCount int64
}

func NewTestContainer() *TestEnvironment {
	if te == nil {
		te = &TestEnvironment{}
	}
	return te
}

func (tc *TestEnvironment) MustInitialize(cfg *appconfig.Config) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.initCount++
	if tc.initCount > 1 {
		return
	}

	logging.SetupLogger(mapAppLoggerConfigToLoggerConfig(&cfg.Logger))

	// Running containers
	testcontainers.Logger = testcontainersLogger{}
	tc.PostgresC = mustSetupPostgresContainer(cfg)
	tc.RedisC = mustSetupRedisContainer(cfg)
	tc.MinioC = mustSetupMinioContainer(cfg)
	tc.SmtpC = mustSetupSmtpContainer(cfg)

	// Setup container
	tc.Container = registry.NewContainer()
	tc.Container.MustInitialize(cfg)

	// Add mock oauth2 provider
	oauthProvider := new(mock3.ProviderMock)
	tc.Container.OauthProviders["mock"] = oauthProvider
}

func (tc *TestEnvironment) Close() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.initCount--
	if tc.initCount > 0 {
		return
	}

	// Terminate running containers
	_ = tc.PostgresC.Terminate(ctx)
	_ = tc.RedisC.Terminate(ctx)
	_ = tc.MinioC.Terminate(ctx)
	_ = tc.SmtpC.Terminate(ctx)
	_ = tc.Container.Close()
}

func mustSetupPostgresContainer(cfg *appconfig.Config) testcontainers.Container {
	// https://github.com/go-testfixtures/testfixtures/blob/c756c9973ec0c741014dce19106369780dc88d37/testfixtures.go#L54
	if !strings.HasSuffix(cfg.Database.Postgres.DBName, "_test") {
		cfg.Database.Postgres.DBName += "_test"
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgis/postgis:17-3.4-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     cfg.Database.Postgres.Username,
				"POSTGRES_PASSWORD": cfg.Database.Postgres.Password,
				"POSTGRES_DB":       cfg.Database.Postgres.DBName,
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to start postgres container")
	}
	host, err := postgresC.Host(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get postgres container host")
	}
	portRaw, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get postgres container port")
	}
	port := portRaw.Int()

	cfg.Database.Postgres.Address = fmt.Sprintf("%s:%d", host, port)
	return postgresC
}

func mustSetupRedisContainer(cfg *appconfig.Config) testcontainers.Container {
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7.4.1-alpine3.20",
			ExposedPorts: []string{"6379/tcp"},
			Env: map[string]string{
				"REDIS_PASSWORD": cfg.Cache.Redis.Password,
			},
			WaitingFor: wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to start redis container")
	}
	host, err := redisC.Host(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get redis container host")
	}
	portRaw, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get redis container port")
	}
	port := portRaw.Int()

	cfg.Cache.Redis.Address = fmt.Sprintf("%s:%d", host, port)
	cfg.Cache.Redis.Username = "default"
	cfg.Cache.Redis.DBIndex = 0

	cfg.PubSub.Redis.Address = fmt.Sprintf("%s:%d", host, port)
	cfg.PubSub.Redis.Username = "default"
	cfg.PubSub.Redis.DBIndex = 0

	return redisC
}

func mustSetupMinioContainer(cfg *appconfig.Config) testcontainers.Container {
	minioC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "minio/minio:RELEASE.2024-10-02T17-50-41Z",
			ExposedPorts: []string{"9000/tcp"},
			Cmd:          []string{"server", "/data"},
			Env: map[string]string{
				"MINIO_ROOT_USER":     cfg.S3.Minio.AccessKey,
				"MINIO_ROOT_PASSWORD": cfg.S3.Minio.SecretKey,
			},
			WaitingFor: wait.ForListeningPort("9000/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to start minio container")
	}
	host, err := minioC.Host(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get minio container host")
	}
	portRaw, err := minioC.MappedPort(ctx, "9000")
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to get minio container port")
	}
	port := portRaw.Int()

	cfg.S3.Minio.Address = fmt.Sprintf("%s:%d", host, port)

	return minioC
}

func mustSetupSmtpContainer(cfg *appconfig.Config) testcontainers.Container {
	smtpC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mailhog/mailhog:v1.0.1",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor:   wait.ForListeningPort("1025/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to start mailhog container")
	}
	cfg.SMTP.Host, err = smtpC.Host(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("Mailhog container host getting error: %s", err.Error())
	}
	port, err := smtpC.MappedPort(ctx, "1025")
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("Mailhog container port getting error: %s", err.Error())
	}
	cfg.SMTP.Port = port.Int()
	cfg.SMTP.Username = ""
	cfg.SMTP.Password = ""
	cfg.SMTP.SSL = false

	return smtpC
}

func mapAppLoggerConfigToLoggerConfig(cfg *appconfig.LoggerConfig) *logging.Config {
	return &logging.Config{
		Level: cfg.Level,
		Console: logging.ConsoleLoggerConfig{
			Enable:   cfg.Console.Enable,
			Encoding: cfg.Console.Encoding,
		},
		File: logging.FileLoggerConfig{
			Enable:  cfg.File.Enable,
			DirPath: cfg.File.DirPath,
			MaxSize: cfg.File.MaxSize,
			MaxAge:  cfg.File.MaxAge,
		},
	}
}

type testcontainersLogger struct{}

func (t testcontainersLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}
