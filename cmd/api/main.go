package main

import (
	"context"
	"fmt"
	appconfig "github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/helper/env"
	"github.com/mandarine-io/Backend/internal/job"
	"github.com/mandarine-io/Backend/internal/registry"
	"github.com/mandarine-io/Backend/internal/transport/http"
	"github.com/mandarine-io/Backend/pkg/logging"
	"github.com/mandarine-io/Backend/pkg/scheduler"
	"github.com/num30/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	syshttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
			"Mandarine: %s\n", env.GetEnvWithDefault("SERVER_VERSION", "0.0.0"),
	)
)

func init() {
	// Print banner
	fmt.Println(banner)
}

func main() {
	// Setup default logger
	log.Logger = zerolog.
		New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).
		With().
		Timestamp().
		Caller().
		Logger()

	// Load config
	var cfg appconfig.Config
	configName := appconfig.GetConfigName()
	err := config.NewConfReader(configName).WithPrefix("MANDARINE").Read(&cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}

	// Setup logger
	logging.SetupLogger(mapAppLoggerConfigToLoggerConfig(&cfg.Logger))

	// Setup container
	container := registry.NewContainer()
	container.MustInitialize(&cfg)
	defer func() {
		_ = container.Close()
	}()

	// Setup scheduler
	jobs := []scheduler.Job{
		job.DeleteExpiredTokensJob(container.Repos.BannedToken),
		job.DeleteExpiredDeletedUsersJob(container.Repos.User),
	}
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
	srv := http.NewServer(container)

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

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to shutdown server")
	}

	log.Info().Msg("the server is shutting down")
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
