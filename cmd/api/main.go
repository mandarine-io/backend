package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mandarine-io/Backend/internal/api/cli"
	appconfig "github.com/mandarine-io/Backend/internal/api/config"
	"github.com/mandarine-io/Backend/internal/api/job"
	"github.com/mandarine-io/Backend/internal/api/registry"
	"github.com/mandarine-io/Backend/internal/api/rest"
	"github.com/mandarine-io/Backend/pkg/config"
	"github.com/mandarine-io/Backend/pkg/logging"
	"github.com/mandarine-io/Backend/pkg/scheduler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	syshttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	versionEnv = "MANDARINE_SERVER__VERSION"
)

var (
	banner = fmt.Sprintf(
		"  __  __                 _            _       \n"+
			" |  \\/  |               | |          (_)      \n"+
			" | \\  / | __ _ _ __   __| | __ _ _ __ _ _ __  \n"+
			" | |\\/| |/ _` | '_ \\ / _` |/ _` | '__| | '_ \\ \n"+
			" | |  | | (_| | | | | (_| | (_| | |  | | | | |\n"+
			" |_|  |_|\\__,_|_| |_|\\__,_|\\__,_|_|  |_|_| |_|\n"+
			"\n"+
			"Mandarine: %s\n", getEnvWithDefault(versionEnv, "0.0.0"),
	)
)

func main() {
	// Parse command line arguments
	fmt.Println(banner)
	options := cli.MustParseCommandLine()

	// Setup logger
	log.Logger = log.Level(zerolog.FatalLevel)
	var loggerCfg appconfig.OnlyLoggerConfig
	config.MustLoadConfig(options.ConfigFilePath, options.EnvFilePath, &loggerCfg)
	logging.SetupLogger(mapAppLoggerConfigToLoggerConfig(&loggerCfg.Logger))

	// Load config
	var cfg appconfig.Config
	config.MustLoadConfig(options.ConfigFilePath, options.EnvFilePath, &cfg)

	// Setup container
	container := registry.NewContainer()
	container.MustInitialize(&cfg)
	defer func() {
		_ = container.Close()
	}()

	// Setup scheduler
	jobs := job.SetupJobs(container)
	cronScheduler := scheduler.MustSetupJobScheduler()
	for _, j := range jobs {
		_, err := cronScheduler.AddJob(j)
		if err != nil {
			log.Warn().Stack().Err(err).Msgf("job %s setup error", j.Name)
		}
	}
	cronScheduler.Start()
	defer func() {
		err := cronScheduler.Shutdown()
		if err != nil {
			log.Warn().Stack().Err(err).Msg("failed to shutdown scheduler")
		}
	}()

	// Create server
	srv := rest.NewServer(container)

	// Run server
	log.Info().Msgf("the server listens on port %d", cfg.Server.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, syshttp.ErrServerClosed) {
			log.Fatal().Stack().Err(err).Msg("failed to start server")
		}
	}()

	// SIGINT and SIGTERM handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	log.Info().Msg("to stop application press Ctrl+C")
	<-quit
	fmt.Println()
	log.Info().Msg("waiting for the server to complete")

	// Shutdown server
	shutdownCtx, shutdownRelease := context.WithTimeout(context.TODO(), 1*time.Second)
	defer shutdownRelease()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}

	log.Info().Msg("the server is shutting down")
}

func getEnvWithDefault(envName, defaultValue string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	return defaultValue
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
