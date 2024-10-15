package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func MustLoadConfig(filePath string, envFilePath string, cfg IConfig) {
	if envFilePath != "" {
		_ = godotenv.Load(envFilePath)
	}
	if filePath != "" {
		mustParseYamlAndEnv(filePath, cfg)
	}
	mustParseSecrets(cfg)

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		slog.Error("Config validation error:", "error", err)
		os.Exit(1)
	}
}

func mustParseSecrets(cfg IConfig) {
	secretConfigInfos := cfg.GetSecretInfos()
	for _, item := range secretConfigInfos {
		mustParseSpecificSecret(item)
	}
}

func mustParseSpecificSecret(item SecretConfigInfo) {
	if item.SecretFileName == "" {
		slog.Warn("It is recommended to use files with secrets: " + item.SecretFileEnvName)
		return
	}

	bytes, err := os.ReadFile(item.SecretFileName)
	if err != nil {
		slog.Warn("Secret file reading error", "error", err)
		return
	}

	*item.SecretValuePtr = string(bytes)
}

func mustParseYamlAndEnv(filePath string, cfg IConfig) {
	if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
		slog.Error("Configuration and environment variables reading error", "error", err)
		os.Exit(1)
	}
}
