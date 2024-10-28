package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"mandarine/pkg/config"
)

const (
	LocalMode       = "local"
	DevelopmentMode = "development"
	ProductionMode  = "production"
	TestMode        = "test"

	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warning"
	LogLevelError = "error"

	LogEncodingText = "text"
	LogEncodingJson = "json"
)

type ServerConfig struct {
	Name           string `yaml:"name" env-description:"Server name" env:"MANDARINE_SERVER__NAME" env-default:"mandarine-server" validate:"required"`
	Port           int    `yaml:"port" env-description:"Server port" env:"MANDARINE_SERVER__PORT" env-default:"8000" validate:"min=-1,max=65535"`
	ExternalOrigin string `yaml:"external_origin" env-description:"Server external origin" env:"MANDARINE_SERVER__EXTERNAL_ORIGIN" env-default:"rest://localhost:8000" validate:"omitempty,http_url"`
	Mode           string `yaml:"mode" env-description:"Server mode" env:"MANDARINE_SERVER__MODE" env-default:"local" validate:"required,oneof=local development production test"`
	Version        string `yaml:"version" env-description:"Server version" env:"MANDARINE_SERVER__VERSION" env-default:"0.0.0"`
}

type PostgresConfig struct {
	Host         string `yaml:"host" env-description:"PostgreSQL host" env:"MANDARINE_POSTGRES__HOST" env-default:"localhost" validate:"required,hostname"`
	Port         int    `yaml:"port" env-description:"PostgreSQL port" env:"MANDARINE_POSTGRES__PORT" env-default:"5432" validate:"min=-1,max=65535"`
	Username     string `yaml:"username" env-description:"PostgreSQL username" env:"MANDARINE_POSTGRES__USERNAME" validate:"required"`
	Password     string `yaml:"password" env-description:"PostgreSQL password (not recommended)" env:"MANDARINE_POSTGRES__PASSWORD" validate:"required"`
	PasswordFile string `yaml:"password_file" env-description:"PostgreSQL password file" env:"MANDARINE_POSTGRES__PASSWORD_FILE" validate:"omitempty,filepath"`
	DBName       string `yaml:"db_name" env-description:"PostgreSQL database name" env:"MANDARINE_POSTGRES__DB_NAME" validate:"required"`
}

type RedisConfig struct {
	Host         string `yaml:"host" env-description:"Redis host" env:"MANDARINE_REDIS__HOST" env-default:"localhost" validate:"required,hostname"`
	Port         int    `yaml:"port" env-description:"Redis port" env:"MANDARINE_REDIS__PORT" env-default:"6379" validate:"min=-1,max=65535"`
	Username     string `yaml:"username" env-description:"Redis username" env:"MANDARINE_REDIS__USERNAME" env-default:"default" validate:"required"`
	Password     string `yaml:"password" env-description:"Redis password (not recommended)" env:"MANDARINE_REDIS__PASSWORD" validate:"required"`
	PasswordFile string `yaml:"password_file" env-description:"Redis password file" env:"MANDARINE_REDIS__PASSWORD_FILE" validate:"omitempty,filepath"`
	DBIndex      int    `yaml:"db_index" env-description:"Redis database index" env:"MANDARINE_REDIS__DB_INDEX" env-default:"0"`
}

type MinioConfig struct {
	Host          string `yaml:"host" env-description:"MinIO host" env:"MANDARINE_MINIO__HOST" env-default:"localhost" validate:"required,hostname"`
	Port          int    `yaml:"port" env-description:"MinIO port" env:"MANDARINE_MINIO__PORT" env-default:"9000" validate:"min=-1,max=65535"`
	AccessKey     string `yaml:"access_key" env-description:"MinIO access key" env:"MANDARINE_MINIO__ACCESS_KEY" validate:"required"`
	SecretKey     string `yaml:"secret_key" env-description:"MinIO secret key" env:"MANDARINE_MINIO__SECRET_KEY" validate:"required"`
	SecretKeyFile string `yaml:"secret_key_file" env-description:"MinIO secret key file" env:"MANDARINE_MINIO__SECRET_KEY_FILE" validate:"omitempty,filepath"`
	BucketName    string `yaml:"bucket_name" env-description:"MinIO bucket name" env:"MANDARINE_MINIO__BUCKET_NAME" validate:"required"`
}

type SmtpConfig struct {
	Host         string `yaml:"host" env-description:"SMTP host" env:"MANDARINE_SMTP__HOST" validate:"required,hostname"`
	Port         int    `yaml:"port" env-description:"SMTP port" env:"MANDARINE_SMTP__PORT" validate:"min=-1,max=65535"`
	Username     string `yaml:"username" env-description:"SMTP username" env:"MANDARINE_SMTP__USERNAME" validate:"required,email"`
	Password     string `yaml:"password" env-description:"SMTP password" env:"MANDARINE_SMTP__PASSWORD" validate:"required"`
	PasswordFile string `yaml:"password_file" env-description:"SMTP password file" env:"MANDARINE_SMTP__PASSWORD_FILE" validate:"omitempty,filepath"`
	SSL          bool   `yaml:"ssl" env-description:"SMTP SSL mode" env:"MANDARINE_SMTP__SSL"`
	From         string `yaml:"from" env-description:"SMTP from" env:"MANDARINE_SMTP__FROM" validate:"omitempty"`
}

type GoogleOAuthClientConfig struct {
	ClientID         string `yaml:"client_id" env-description:"Google OAuth client ID" env:"MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_ID" validate:"required"`
	ClientSecret     string `yaml:"client_secret" env-description:"Google OAuth client secret" env:"MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_SECRET" validate:"required"`
	ClientSecretFile string `yaml:"client_secret_file" env-description:"Google OAuth client secret" env:"MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_SECRET_FILE" validate:"omitempty,filepath"`
}

type YandexOAuthClientConfig struct {
	ClientID         string `yaml:"client_id" env-description:"Yandex OAuth client ID" env:"MANDARINE_YANDEX_OAUTH_CLIENT__CLIENT_ID" validate:"required"`
	ClientSecret     string `yaml:"client_secret" env-description:"Yandex OAuth client secret" env:"MANDARINE_YANDEX_OAUTH_CLIENT__CLIENT_SECRET" validate:"required"`
	ClientSecretFile string `yaml:"client_secret_file" env-description:"Yandex OAuth client secret" env:"MANDARINE_YANDEX_OAUTH_CLIENT__CLIENT_SECRET_FILE" validate:"omitempty,filepath"`
}

type MailRuOAuthClientConfig struct {
	ClientID         string `yaml:"client_id" env-description:"Mail RU OAuth client ID" env:"MANDARINE_MAIL_RU_OAUTH_CLIENT__CLIENT_ID" validate:"required"`
	ClientSecret     string `yaml:"client_secret" env-description:"Mail RU OAuth client secret" env:"MANDARINE_MAIL_RU_OAUTH_CLIENT__CLIENT_SECRET" validate:"required"`
	ClientSecretFile string `yaml:"client_secret_file" env-description:"Mail RU OAuth client secret" env:"MANDARINE_MAIL_RU_OAUTH_CLIENT_CLIENT__SECRET_FILE" validate:"omitempty,filepath"`
}

type OAuthClientConfig struct {
	Google GoogleOAuthClientConfig `yaml:"google"`
	Yandex YandexOAuthClientConfig `yaml:"yandex"`
	MailRu MailRuOAuthClientConfig `yaml:"mail_ru"`
}

type CacheConfig struct {
	TTL int `yaml:"ttl" env-description:"Cache TTL (seconds)" env:"MANDARINE_CACHE__TTL" env-default:"120" validate:"required,min=0"`
}

type JWTConfig struct {
	Secret          string `yaml:"secret" env-description:"JWT secret" env:"MANDARINE_JWT__SECRET" validate:"required"`
	SecretFile      string `yaml:"secret_file" env-description:"JWT secret file" env:"MANDARINE_JWT__SECRET_FILE" validate:"omitempty,filepath"`
	AccessTokenTTL  int    `yaml:"access_token_ttl" env-description:"JWT access token TTL (seconds)" env:"MANDARINE_JWT__ACCESS_TOKEN_TTL" env-default:"3600"  validate:"required,min=0"`
	RefreshTokenTTL int    `yaml:"refresh_token_ttl" env-description:"JWT refresh token TTL (seconds)" env:"MANDARINE_JWT__REFRESH_TOKEN_TTL" env-default:"86400"  validate:"required,min=0"`
}

type OTPConfig struct {
	Length int `yaml:"length" env-description:"OTP length" env:"MANDARINE_OTP__LENGTH" env-default:"6" validate:"required,min=4"`
	TTL    int `yaml:"ttl" env-description:"OTP TTL (seconds)" env:"MANDARINE_OTP__TTL" env-default:"600" validate:"required,min=0"`
}

type RateLimitConfig struct {
	RPS int `yaml:"rps" env-description:"Rate limit RPS" env:"MANDARINE_RATE_LIMIT__RPS" env-default:"100" validate:"required,min=1"`
}

type SecurityConfig struct {
	JWT       JWTConfig       `yaml:"jwt"`
	OTP       OTPConfig       `yaml:"otp"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type LocaleConfig struct {
	Path     string `yaml:"path" env-description:"Locales path" env:"MANDARINE_LOCALE__PATH" env-default:"locales" validate:"required"`
	Language string `yaml:"language" env-description:"Locale default language" env:"MANDARINE_LOCALE__LANGUAGE" env-default:"ru" validate:"required"`
}

type TemplateConfig struct {
	Path string `yaml:"path" env-description:"Templates path" env:"MANDARINE_TEMPLATE__PATH" env-default:"templates" validate:"required"`
}

type MigrationConfig struct {
	Path string `yaml:"path" env-description:"Migrations path" env:"MANDARINE_MIGRATION__PATH" env-default:"migrations" validate:"required"`
}

type LoggerConfig struct {
	Level   string              `yaml:"level" env-description:"Logger level" env:"MANDARINE_LOGGER__LEVEL" env-default:"info" validate:"required"`
	Console ConsoleLoggerConfig `yaml:"console"`
	File    FileLoggerConfig    `yaml:"file"`
}

type ConsoleLoggerConfig struct {
	Enable   bool   `yaml:"enable" env-description:"Console logger is enabled" env:"MANDARINE_LOGGER__CONSOLE_ENABLE" env-default:"false"`
	Encoding string `yaml:"encoding" env-description:"Console logger encoding" env:"MANDARINE_LOGGER__CONSOLE_ENCODING" env-default:"text" validate:"required_with=Enable"`
}

type FileLoggerConfig struct {
	Enable  bool   `yaml:"enable" env-description:"File logger is enabled" env:"MANDARINE_LOGGER__FILE_ENABLE" env-default:"false"`
	DirPath string `yaml:"dir_path" env-description:"Log directory path" env:"MANDARINE_LOGGER__FILE_DIR_PATH" env-default:"logs" validate:"required_with=Enable"`
	MaxSize int    `yaml:"max_size" env-description:"Log file max size (MB)" env:"MANDARINE_LOGGER__FILE_MAX_SIZE" env-default:"100" validate:"required_with=Enable,min=0"`
	MaxAge  int    `yaml:"max_age" env-description:"Log file max age (days)" env:"MANDARINE_LOGGER__FILE_MAX_AGE" env-default:"30" validate:"required_with=Enable,min=0"`
}

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Postgres    PostgresConfig    `yaml:"postgres"`
	Redis       RedisConfig       `yaml:"redis"`
	Minio       MinioConfig       `yaml:"minio"`
	SMTP        SmtpConfig        `yaml:"smtp"`
	Cache       CacheConfig       `yaml:"cache"`
	Locale      LocaleConfig      `yaml:"locale"`
	Template    TemplateConfig    `yaml:"template"`
	Migrations  MigrationConfig   `yaml:"migrations"`
	Logger      LoggerConfig      `yaml:"logger"`
	OAuthClient OAuthClientConfig `yaml:"oauth_client"`
	Security    SecurityConfig    `yaml:"security"`
}

type OnlyLoggerConfig struct {
	Logger LoggerConfig `yaml:"logger"`
}

func (c *OnlyLoggerConfig) GetSecretInfos() []config.SecretConfigInfo {
	return make([]config.SecretConfigInfo, 0)
}

func (c *Config) GetSecretInfos() []config.SecretConfigInfo {
	return []config.SecretConfigInfo{
		{
			SecretFileEnvName: "MANDARINE_POSTGRES__PASSWORD_FILE",
			SecretFileName:    c.Postgres.PasswordFile,
			SecretValuePtr:    &c.Postgres.Password,
		},
		{
			SecretFileEnvName: "MANDARINE_REDIS__PASSWORD_FILE",
			SecretFileName:    c.Redis.PasswordFile,
			SecretValuePtr:    &c.Redis.Password,
		},
		{
			SecretFileEnvName: "MANDARINE_MINIO__SECRET_KEY_FILE",
			SecretFileName:    c.Minio.SecretKeyFile,
			SecretValuePtr:    &c.Minio.SecretKey,
		},
		{
			SecretFileEnvName: "MANDARINE_SMTP__PASSWORD_FILE",
			SecretFileName:    c.SMTP.PasswordFile,
			SecretValuePtr:    &c.SMTP.Password,
		},
		{
			SecretFileEnvName: "MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_SECRET_FILE",
			SecretFileName:    c.OAuthClient.Google.ClientSecretFile,
			SecretValuePtr:    &c.OAuthClient.Google.ClientSecret,
		},
		{
			SecretFileEnvName: "MANDARINE_YANDEX_OAUTH_CLIENT__CLIENT_SECRET_FILE",
			SecretFileName:    c.OAuthClient.Yandex.ClientSecretFile,
			SecretValuePtr:    &c.OAuthClient.Yandex.ClientSecret,
		},
		{
			SecretFileEnvName: "MANDARINE_MAIL_RU_OAUTH_CLIENT__CLIENT_SECRET_FILE",
			SecretFileName:    c.OAuthClient.MailRu.ClientSecretFile,
			SecretValuePtr:    &c.OAuthClient.MailRu.ClientSecret,
		},
		{
			SecretFileEnvName: "MANDARINE_JWT__SECRET_FILE",
			SecretFileName:    c.Security.JWT.SecretFile,
			SecretValuePtr:    &c.Security.JWT.Secret,
		},
	}
}

func GetDescription() string {
	help, _ := cleanenv.GetDescription(&Config{}, nil)
	return help
}
