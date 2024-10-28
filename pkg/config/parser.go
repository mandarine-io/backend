package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func MustLoadConfig(filePath string, envFilePath string, cfg interface{}) {
	if envFilePath != "" {
		log.Debug().Msg("loading env vars from dotenv file")
		_ = godotenv.Load(envFilePath)
	}

	if filePath != "" {
		log.Debug().Msg("loading config from file and env vars")
		mustParseYamlAndEnv(filePath, cfg)
	}

	log.Debug().Msg("validating config")
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to validate config")
	}
}

func mustParseYamlAndEnv(filePath string, cfg interface{}) {
	if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to read config file and env vars")
	}
}
