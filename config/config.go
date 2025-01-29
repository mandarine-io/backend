package config

import (
	"github.com/mandarine-io/backend/internal/util/env"
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
	configPath := env.GetEnvWithDefault("APP_CONFIG_FILE", "config/config.yaml")
	oldnew := make([]string, 2*len(viper.SupportedExts))
	for i, ext := range viper.SupportedExts {
		oldnew[2*i] = "." + ext
		oldnew[2*i+1] = ""
	}
	return strings.NewReplacer(oldnew...).Replace(configPath)
}

type Config struct {
	Server             ServerConfig
	Database           PostgresDatabaseConfig
	Cache              RedisCacheConfig
	S3                 MinIOS3Config
	PubSub             RedisPubSubConfig
	SMTP               SMTPConfig
	Websocket          WebsocketConfig
	Locale             LocaleConfig
	Template           TemplateConfig
	Migrations         MigrationConfig
	Logger             LoggerConfig
	OAuthProviders     []OauthProviderItemConfig
	GeocodingProviders []GeocodingProviderItemConfig
	Security           SecurityConfig
}

////////// Server //////////

type ServerConfig struct {
	Name        string `default:"server" validate:"required"`
	Port        int    `default:"8000" validate:"min=-1,max=65535"`
	ExternalURL string `default:"http://localhost:8080" validate:"omitempty,http_url"`
	Mode        Mode   `default:"local" validate:"required,oneof=local development production test"`
	Version     string `default:"0.0.0"`
}

////////// Database //////////

type PostgresDatabaseConfig struct {
	Address  string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
}

////////// Cache //////////

type RedisCacheConfig struct {
	Address  string `validate:"required"`
	Username string `default:"default" validate:"required"`
	Password string `validate:"required"`
	DBIndex  int    `default:"0" validate:"min=0"`
	TTL      int    `default:"86400" validate:"required,min=0"`
}

////////// S3 //////////

type MinIOS3Config struct {
	Address   string `validate:"required"`
	AccessKey string `validate:"required"`
	SecretKey string `validate:"required"`
	Bucket    string `validate:"required"`
}

/////////// SMTP //////////

type SMTPConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Username string
	Password string
	SSL      bool
	From     string `validate:"required"`
}

////////// PubSub //////////

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

type OauthProviderItemConfig struct {
	Name         string `validate:"required"`
	ClientID     string `validate:"required"`
	ClientSecret string `validate:"required"`
}

////////// Geocoding Clients //////////

type GeocodingProviderItemConfig struct {
	Name   string `validate:"required"`
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
	LogEncodingJSON = LogEncoding("json")
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
