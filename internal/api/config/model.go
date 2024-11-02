package config

import (
	"github.com/ilyakaznacheev/cleanenv"
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
	MaxRequestSize int    `yaml:"max_request_size" env-description:"Max request size in bytes" env:"MANDARINE_SERVER__MAX_REQUEST_SIZE" env-default:"500" validate:"min=0"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-description:"PostgreSQL host" env:"MANDARINE_POSTGRES__HOST" env-default:"localhost" validate:"required,hostname"`
	Port     int    `yaml:"port" env-description:"PostgreSQL port" env:"MANDARINE_POSTGRES__PORT" env-default:"5432" validate:"min=-1,max=65535"`
	Username string `yaml:"username" env-description:"PostgreSQL username" env:"MANDARINE_POSTGRES__USERNAME" validate:"required"`
	Password string `yaml:"password" env-description:"PostgreSQL password (not recommended)" env:"MANDARINE_POSTGRES__PASSWORD" validate:"required"`
	DBName   string `yaml:"db_name" env-description:"PostgreSQL database name" env:"MANDARINE_POSTGRES__DB_NAME" validate:"required"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env-description:"Redis host" env:"MANDARINE_REDIS__HOST" env-default:"localhost" validate:"required,hostname"`
	Port     int    `yaml:"port" env-description:"Redis port" env:"MANDARINE_REDIS__PORT" env-default:"6379" validate:"min=-1,max=65535"`
	Username string `yaml:"username" env-description:"Redis username" env:"MANDARINE_REDIS__USERNAME" env-default:"default" validate:"required"`
	Password string `yaml:"password" env-description:"Redis password (not recommended)" env:"MANDARINE_REDIS__PASSWORD" validate:"required"`
	DBIndex  int    `yaml:"db_index" env-description:"Redis database index" env:"MANDARINE_REDIS__DB_INDEX" env-default:"0"`
}

type MinioConfig struct {
	Host       string `yaml:"host" env-description:"MinIO host" env:"MANDARINE_MINIO__HOST" env-default:"localhost" validate:"required,hostname"`
	Port       int    `yaml:"port" env-description:"MinIO port" env:"MANDARINE_MINIO__PORT" env-default:"9000" validate:"min=-1,max=65535"`
	AccessKey  string `yaml:"access_key" env-description:"MinIO access key" env:"MANDARINE_MINIO__ACCESS_KEY" validate:"required"`
	SecretKey  string `yaml:"secret_key" env-description:"MinIO secret key" env:"MANDARINE_MINIO__SECRET_KEY" validate:"required"`
	BucketName string `yaml:"bucket_name" env-description:"MinIO bucket name" env:"MANDARINE_MINIO__BUCKET_NAME" validate:"required"`
}

type SmtpConfig struct {
	Host     string `yaml:"host" env-description:"SMTP host" env:"MANDARINE_SMTP__HOST" validate:"required,hostname"`
	Port     int    `yaml:"port" env-description:"SMTP port" env:"MANDARINE_SMTP__PORT" validate:"min=-1,max=65535"`
	Username string `yaml:"username" env-description:"SMTP username" env:"MANDARINE_SMTP__USERNAME"`
	Password string `yaml:"password" env-description:"SMTP password" env:"MANDARINE_SMTP__PASSWORD"`
	SSL      bool   `yaml:"ssl" env-description:"SMTP SSL mode" env:"MANDARINE_SMTP__SSL"`
	From     string `yaml:"from" env-description:"SMTP from" env:"MANDARINE_SMTP__FROM" validate:"omitempty"`
}

type WebsocketConfig struct {
	PoolSize int `yaml:"pool_size" env-description:"Websocket pool size" env:"MANDARINE_WEBSOCKET__POOL_SIZE" env-default:"1024" validate:"min=0"`
}

type GoogleOAuthClientConfig struct {
	ClientID     string `yaml:"client_id" env-description:"Google OAuth client id" env:"MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_ID" validate:"required"`
	ClientSecret string `yaml:"client_secret" env-description:"Google OAuth client secret" env:"MANDARINE_GOOGLE_OAUTH_CLIENT__CLIENT_SECRET" validate:"required"`
}

type YandexOAuthClientConfig struct {
	ClientID     string `yaml:"client_id" env-description:"Yandex OAuth client id" env:"MANDARINE_YANDEX_OAUTH_CLIENT__CLIENT_ID" validate:"required"`
	ClientSecret string `yaml:"client_secret" env-description:"Yandex OAuth client secret" env:"MANDARINE_YANDEX_OAUTH_CLIENT__CLIENT_SECRET" validate:"required"`
}

type MailRuOAuthClientConfig struct {
	ClientID     string `yaml:"client_id" env-description:"Mail RU OAuth client id" env:"MANDARINE_MAIL_RU_OAUTH_CLIENT__CLIENT_ID" validate:"required"`
	ClientSecret string `yaml:"client_secret" env-description:"Mail RU OAuth client secret" env:"MANDARINE_MAIL_RU_OAUTH_CLIENT__CLIENT_SECRET" validate:"required"`
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
	Websocket   WebsocketConfig   `yaml:"websocket"`
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

func GetDescription() string {
	help, _ := cleanenv.GetDescription(&Config{}, nil)
	return help
}
