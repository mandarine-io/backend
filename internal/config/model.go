package config

import (
	"github.com/mandarine-io/Backend/internal/helper/env"
	"github.com/spf13/viper"
	"strings"
)

type Mode string

const (
	LocalMode       = Mode("local")
	DevelopmentMode = Mode("development")
	ProductionMode  = Mode("production")
	TestMode        = Mode("test")
)

func GetConfigName() string {
	configPath := env.GetEnvWithDefault("MANDARINE_CONFIG_FILE", "config/config.yaml")
	oldnew := make([]string, 2*len(viper.SupportedExts))
	for i, ext := range viper.SupportedExts {
		oldnew[2*i] = "." + ext
		oldnew[2*i+1] = ""
	}
	return strings.NewReplacer(oldnew...).Replace(configPath)
}

type Config struct {
	Server           ServerConfig
	Database         DatabaseConfig
	Cache            CacheConfig
	S3               S3Config
	PubSub           PubSubConfig
	SMTP             SmtpConfig
	Websocket        WebsocketConfig
	Locale           LocaleConfig
	Template         TemplateConfig
	Migrations       MigrationConfig
	Logger           LoggerConfig
	OAuthClients     map[string]OauthClientConfig
	GeocodingClients map[string]GeocodingClientConfig
	Security         SecurityConfig
}

////////// Server //////////

type ServerConfig struct {
	Name           string `default:"server" validate:"required"`
	Port           int    `default:"8000" validate:"min=-1,max=65535"`
	ExternalOrigin string `default:"http://localhost:8080" validate:"omitempty,http_url"`
	Mode           Mode   `default:"local" validate:"required,oneof=local development production test"`
	Version        string `default:"0.0.0"`
	MaxRequestSize int    `default:"524288000" validate:"min=0"`
	RPS            int    `default:"100" validate:"required,min=1"`
}

////////// Database //////////

type DatabaseType string

const (
	PostgresDatabaseType = DatabaseType("postgres")
)

type DatabaseConfig struct {
	Type     DatabaseType            `default:"postgres" validate:"required,oneof=postgres"`
	Postgres *PostgresDatabaseConfig `validate:"required_if=Type postgres"`
}

type PostgresDatabaseConfig struct {
	Address  string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
}

////////// Cache //////////

type CacheType string

const (
	MemoryCacheType       = CacheType("memory")
	RedisCacheType        = CacheType("redis")
	RedisClusterCacheType = CacheType("redis_cluster")
)

type CacheConfig struct {
	TTL   int               `default:"120" validate:"required,min=0"`
	Type  CacheType         `default:"memory" validate:"required,oneof=memory redis"`
	Redis *RedisCacheConfig `validate:"required_if=Type redis"`
}

type RedisCacheConfig struct {
	Address  string `validate:"required"`
	Username string `default:"default" validate:"required"`
	Password string `validate:"required"`
	DBIndex  int    `default:"0" validate:"min=0"`
}

////////// S3 //////////

type S3Type string

const (
	MinioS3Type = S3Type("minio")
)

type S3Config struct {
	Type  S3Type         `default:"minio" validate:"required,oneof=minio"`
	Minio *MinioS3Config `validate:"required_if=Type minio"`
}

type MinioS3Config struct {
	Address   string `validate:"required"`
	AccessKey string `validate:"required"`
	SecretKey string `validate:"required"`
	Bucket    string `validate:"required"`
}

/////////// SMTP //////////

type SmtpConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Username string
	Password string
	SSL      bool
	From     string `validate:"required"`
}

////////// PubSub //////////

type PubSubType string

const (
	MemoryPubSubType = PubSubType("memory")
	RedisPubSubType  = PubSubType("redis")
)

type PubSubConfig struct {
	Type  PubSubType         `default:"memory" validate:"required,oneof=memory redis"`
	Redis *RedisPubSubConfig `validate:"required_if=Type redis"`
}

type RedisPubSubConfig struct {
	Address  string `validate:"required"`
	Username string
	Password string
	DBIndex  int `default:"0" validate:"min=0"`
}

////////// Websocket //////////

type WebsocketConfig struct {
	PoolSize int `default:"1024" validate:"min=0"`
}

////////// Oauth 2.0 Clients //////////

type OauthClientConfig struct {
	ClientID     string `validate:"required"`
	ClientSecret string `validate:"required"`
}

////////// Geocoding Clients //////////

type GeocodingClientConfig struct {
	APIKey string `validate:"required"`
}

////////// Security //////////

type SecurityConfig struct {
	JWT JWTConfig
	OTP OTPConfig
}

type JWTConfig struct {
	Secret          string `validate:"required"`
	AccessTokenTTL  int    `default:"3600" validate:"required,min=0"`
	RefreshTokenTTL int    `default:"86400" validate:"required,min=0"`
}

type OTPConfig struct {
	Length int `default:"6" validate:"required,min=4"`
	TTL    int `default:"600" validate:"required,min=0"`
}

////////// Locale //////////

type LocaleConfig struct {
	Path     string `default:"locales" validate:"required"`
	Language string `default:"ru" validate:"required"`
}

///////// Template //////////

type TemplateConfig struct {
	Path string `default:"templates" validate:"required"`
}

///////// Migrations //////////

type MigrationConfig struct {
	Path string `default:"migrations" validate:"required"`
}

///////// Logger //////////

type LogEncoding string

const (
	LogLevelDebug = LogLevel("debug")
	LogLevelInfo  = LogLevel("info")
	LogLevelWarn  = LogLevel("warning")
	LogLevelError = LogLevel("error")
)

type LogLevel string

const (
	LogEncodingText = LogEncoding("text")
	LogEncodingJson = LogEncoding("json")
)

type LoggerConfig struct {
	Level   string `default:"info" validate:"required"`
	Console ConsoleLoggerConfig
	File    FileLoggerConfig
}

type ConsoleLoggerConfig struct {
	Enable   bool   `default:"true"`
	Encoding string `default:"text" validate:"required_with=Enable"`
}

type FileLoggerConfig struct {
	Enable  bool   `default:"false"`
	DirPath string `default:"logs" validate:"required_with=Enable"`
	MaxSize int    `default:"100" validate:"required_with=Enable,min=0"`
	MaxAge  int    `default:"30" validate:"required_with=Enable,min=0"`
}
