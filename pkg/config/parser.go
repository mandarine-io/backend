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
		log.Debug().Msg("loading env vars from dotenv file")
		_ = godotenv.Load(envFilePath)
	}
	if filePath != "" {
		log.Debug().Msg("loading config from file and env vars")
		mustParseYamlAndEnv(filePath, cfg)
	}
	mustParseSecrets(cfg)

	log.Debug().Msg("validating config")
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to validate config")
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
		log.Fatal().Stack().Err(err).Msg("failed to read config file and env vars")
	}
}
