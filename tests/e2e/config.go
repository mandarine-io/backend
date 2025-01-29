package e2e

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mandarine-io/backend/internal/infrastructure/cache/redis"
	"github.com/mandarine-io/backend/internal/infrastructure/database/gorm/postgres"
	"github.com/mandarine-io/backend/internal/infrastructure/s3/minio"
	"github.com/mandarine-io/backend/internal/infrastructure/smtp"
	"github.com/rs/zerolog/log"
)

var Cfg = mustLoadConfig()

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Minio    MinioConfig
	Mailhog  MailhogConfig
	Server   ServerConfig
}

type PostgresConfig struct {
	Port     int    `env:"APP_TEST_POSTGRES_PORT" env-default:"35432"`
	User     string `env:"APP_TEST_POSTGRES_USER" env-default:"admin"`
	Password string `env:"APP_TEST_POSTGRES_PASSWORD" env-default:"password"`
	DB       string `env:"APP_TEST_POSTGRES_DB" env-default:"mandarine"`
}

type RedisConfig struct {
	Port     int    `env:"APP_TEST_REDIS_PORT" env-default:"36379"`
	User     string `env:"APP_TEST_REDIS_USER" env-default:"default"`
	Password string `env:"APP_TEST_REDIS_PASSWORD" env-default:"password"`
	DBIndex  int    `env:"APP_TEST_REDIS_DBINDEX" env-default:"0"`
}

type MinioConfig struct {
	Port      int    `env:"APP_TEST_MINIO_PORT" env-default:"39000"`
	AccessKey string `env:"APP_TEST_MINIO_ACCESSKEY" env-default:"admin"`
	SecretKey string `env:"APP_TEST_MINIO_SECRETKEY" env-default:"Password_10"`
	Bucket    string `env:"APP_TEST_MINIO_BUCKET" env-default:"mandarine"`
}

type MailhogConfig struct {
	SMTPPort int `env:"APP_TEST_MAILHOG_SMTPPORT" env-default:"31025"`
	APIPort  int `env:"APP_TEST_MAILHOG_APIPORT" env-default:"38025"`
}

type ServerConfig struct {
	Port      int    `env:"APP_TEST_SERVER_PORT" env-default:"38080"`
	JWTSecret string `env:"APP_TEST_SERVER_JWTSECRET" env-default:"9bd8b3e960d752f050950dcec783aaae1e0437baa2f29310d556116448b9471c"`
}

func mustLoadConfig() Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatal().Err(err).Msg("failed to read config from env")
	}

	return config
}

func (c Config) GetPostgresConfig() postgres.Config {
	return postgres.Config{
		Address:  fmt.Sprintf("localhost:%d", c.Postgres.Port),
		Username: c.Postgres.User,
		Password: c.Postgres.Password,
		DBName:   c.Postgres.DB,
	}
}

func (c Config) GetRedisConfig() redis.Config {
	return redis.Config{
		Address:  fmt.Sprintf("localhost:%d", c.Redis.Port),
		Username: c.Redis.User,
		Password: c.Redis.Password,
		DBIndex:  c.Redis.DBIndex,
	}
}

func (c Config) GetMinioConfig() minio.Config {
	return minio.Config{
		Address:    fmt.Sprintf("localhost:%d", c.Minio.Port),
		AccessKey:  c.Minio.AccessKey,
		SecretKey:  c.Minio.SecretKey,
		BucketName: c.Minio.Bucket,
	}
}

func (c Config) GetSMTPConfig() smtp.Config {
	return smtp.Config{
		Host: "localhost",
		Port: c.Mailhog.SMTPPort,
		SSL:  false,
	}
}

func (c Config) GetMailhogAPIURL() string {
	return fmt.Sprintf("http://localhost:%d", c.Mailhog.APIPort)
}

func (c Config) GetServerURL() string {
	return fmt.Sprintf("http://localhost:%d", c.Server.Port)
}

func (c Config) GetServerJWTSecret() string {
	return c.Server.JWTSecret
}
