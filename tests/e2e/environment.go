package e2e

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	appconfig "mandarine/internal/api/config"
	"mandarine/internal/api/registry"
	"mandarine/internal/api/service/auth"
	"mandarine/pkg/logging"
	mock3 "mandarine/pkg/oauth/mock"
	"os"
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
	tc.Container = registry.MustNewContainer(cfg)

	// Add mock oauth2 provider
	oauthProvider := new(mock3.OAuthProviderMock)
	tc.Container.OauthProviders["mock"] = oauthProvider

	// Add mock social login service
	tc.Container.Services.SocialLogins["mock"] = auth.NewSocialLoginService(tc.Container.Repositories.User, oauthProvider, "mock", cfg)
}

func (tc *TestEnvironment) Close() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.initCount--
	if tc.initCount > 0 {
		return
	}

	// Terminate running containers
	err := tc.PostgresC.Terminate(ctx)
	if err != nil {
		slog.Warn("Postgres container terminate error", logging.ErrorAttr(err))
	}
	err = tc.RedisC.Terminate(ctx)
	if err != nil {
		slog.Warn("Redis container terminate error", logging.ErrorAttr(err))
	}
	err = tc.MinioC.Terminate(ctx)
	if err != nil {
		slog.Warn("Minio container terminate error", logging.ErrorAttr(err))
	}
	err = tc.SmtpC.Terminate(ctx)
	if err != nil {
		slog.Warn("Smtp container terminate error", logging.ErrorAttr(err))
	}
	err = tc.Container.Close()
	if err != nil {
		slog.Warn("Container closing error", logging.ErrorAttr(err))
	}
}

func mustSetupPostgresContainer(cfg *appconfig.Config) testcontainers.Container {
	// https://github.com/go-testfixtures/testfixtures/blob/c756c9973ec0c741014dce19106369780dc88d37/testfixtures.go#L54
	if !strings.HasSuffix(cfg.Postgres.DBName, "_test") {
		cfg.Postgres.DBName += "_test"
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgis/postgis:17-3.4-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     cfg.Postgres.Username,
				"POSTGRES_PASSWORD": cfg.Postgres.Password,
				"POSTGRES_DB":       cfg.Postgres.DBName,
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		slog.Error("Postgres container setup error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.Postgres.Host, err = postgresC.Host(ctx)
	if err != nil {
		slog.Error("Postgres container host error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		slog.Error("Postgres container port error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.Postgres.Port = port.Int()

	slog.Info(fmt.Sprintf("Postgres container running at %s:%d", cfg.Postgres.Host, cfg.Postgres.Port))
	return postgresC
}

func mustSetupRedisContainer(cfg *appconfig.Config) testcontainers.Container {
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7.4.1-alpine3.20",
			ExposedPorts: []string{"6379/tcp"},
			Env: map[string]string{
				"REDIS_PASSWORD": cfg.Redis.Password,
			},
			WaitingFor: wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	if err != nil {
		slog.Error("Redis container setup error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.Redis.Host, err = redisC.Host(ctx)
	if err != nil {
		slog.Error("Redis container host error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		slog.Error("Redis container port error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.Redis.Port = port.Int()
	cfg.Redis.Username = "default"
	cfg.Redis.DBIndex = 0

	slog.Info(fmt.Sprintf("Redis container running at %s:%d", cfg.Redis.Host, cfg.Redis.Port))
	return redisC
}

func mustSetupMinioContainer(cfg *appconfig.Config) testcontainers.Container {
	minioC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "minio/minio:RELEASE.2024-10-02T17-50-41Z",
			ExposedPorts: []string{"9000/tcp"},
			Cmd:          []string{"server", "/data"},
			Env: map[string]string{
				"MINIO_ROOT_USER":     cfg.Minio.AccessKey,
				"MINIO_ROOT_PASSWORD": cfg.Minio.SecretKey,
			},
			WaitingFor: wait.ForListeningPort("9000/tcp"),
		},
		Started: true,
	})
	if err != nil {
		slog.Error("Minio container setup error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.Minio.Host, err = minioC.Host(ctx)
	if err != nil {
		slog.Error("Minio container host error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	port, err := minioC.MappedPort(ctx, "9000")
	if err != nil {
		slog.Error("Minio container port error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.Minio.Port = port.Int()

	slog.Info(fmt.Sprintf("Minio container running at %s:%d", cfg.Minio.Host, cfg.Minio.Port))
	return minioC
}

func mustSetupSmtpContainer(cfg *appconfig.Config) testcontainers.Container {
	smtpC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mailhog/mailhog",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor:   wait.ForListeningPort("1025/tcp"),
		},
		Started: true,
	})
	if err != nil {
		slog.Error("Mailhog container setup error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.SMTP.Host, err = smtpC.Host(ctx)
	if err != nil {
		slog.Error("Mailhog container host error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	port, err := smtpC.MappedPort(ctx, "1025")
	if err != nil {
		slog.Error("Mailhog container port error", logging.ErrorAttr(err))
		os.Exit(1)
	}
	cfg.SMTP.Port = port.Int()
	cfg.SMTP.Username = ""
	cfg.SMTP.Password = ""
	cfg.SMTP.SSL = false

	slog.Info(fmt.Sprintf("Mailhog container running at %s:%d", cfg.SMTP.Host, cfg.SMTP.Port))
	return smtpC
}

func mapAppLoggerConfigToLoggerConfig(cfg *appconfig.LoggerConfig) *logging.Config {
	return &logging.Config{
		Console: logging.ConsoleLoggerConfig{
			Level:    cfg.Console.Level,
			Encoding: cfg.Console.Encoding,
		},
		File: logging.FileLoggerConfig{
			Enable:  cfg.File.Enable,
			Level:   cfg.File.Level,
			DirPath: cfg.File.DirPath,
			MaxSize: cfg.File.MaxSize,
			MaxAge:  cfg.File.MaxAge,
		},
	}
}

type testcontainersLogger struct{}

func (t testcontainersLogger) Printf(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...))
}
