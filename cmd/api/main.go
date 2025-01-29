package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	appconfig "github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/di/finalizer"
	"github.com/mandarine-io/backend/internal/di/initializer"
	"github.com/mandarine-io/backend/internal/logging"
	"github.com/mandarine-io/backend/internal/scheduler"
	"github.com/mandarine-io/backend/internal/scheduler/job"
	"github.com/mandarine-io/backend/internal/transport/http"
	"github.com/mandarine-io/backend/internal/util/env"
	"github.com/num30/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	syshttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	version = env.GetEnvWithDefault("APP_SERVER_VERSION", "0.0.0")
	banner  = fmt.Sprintf(
		"  __  __              _          _          \n"+
			" |  \\/  |__ _ _ _  __| |__ _ _ _(_)_ _  ___ \n"+
			" | |\\/| / _` | ' \\/ _` / _` | '_| | ' \\/ -_)\n"+
			" |_|  |_\\__,_|_||_\\__,_\\__,_|_| |_|_||_\\___|\n"+
			"\n"+
			"Mandarine: %s\n", version,
	)
)

func init() {
	fmt.Println(banner)

	// Setup default logger
	log.Logger = zerolog.
		New(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			},
		).
		With().
		Timestamp().
		Caller().
		Logger()
}

func main() {
	// Load env
	_ = godotenv.Load()

	// Load config
	var cfg appconfig.Config
	configName := appconfig.GetConfigName()

	err := config.NewConfReader(configName).WithPrefix("APP").Read(&cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to load config")
	}

	mustStartApp(cfg)
}

func mustStartApp(cfg appconfig.Config) {
	// Setup logger
	log.Debug().Msg("setup logging")
	logging.SetupLogger(toLoggerConfig(cfg.Logger))

	// Setup container
	container := di.NewContainer(cfg)
	container.RegisterInitializers(
		initializer.Locale(container),
		initializer.Template(container),
		initializer.Cache(container),
		initializer.GormDatabase(container),
		initializer.S3(container),
		initializer.SMTP(container),
		initializer.PubSub(container),
		initializer.Websocket(container),
		initializer.ThirdParty(container),
		initializer.Metrics(container),
		initializer.GormRepositories(container),
		initializer.Services(container),
		initializer.Handlers(container),
		initializer.Scheduler(container),
	)
	container.RegisterFinalizers(
		finalizer.Cache(container),
		finalizer.GormDatabase(container),
		finalizer.PubSub(container),
		finalizer.Websocket(container),
		finalizer.Scheduler(container),
	)

	log.Info().Msg("initialize DI container")
	err := container.Initialize()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to initialize app")
	}

	defer func() {
		log.Info().Msg("finalize DI container")
		err = container.Finalize()
		if err != nil {
			log.Error().Err(err).Msg("failed to finalize app")
		}
	}()

	// Setup scheduler
	jobs := []scheduler.Job{
		job.DeleteExpiredDeletedUsersJob(container.Repos.User),
	}
	for _, j := range jobs {
		_, err = container.Infrastructure.Scheduler.AddJob(j)
		if err != nil {
			log.Warn().Stack().Err(err).Msgf("job %s setup error", j.Name)
		}
	}

	// Create server
	srv := http.NewServer(container)

	// Run server
	log.Info().Msgf("server listens on port %d", cfg.Server.Port)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, syshttp.ErrServerClosed) {
			log.Fatal().Stack().Err(err).Msg("failed to start server")
		}
	}()

	// SIGINT and SIGTERM handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wait for signal
	log.Info().Msg("to stop application press Ctrl+C")

	<-quit
	fmt.Println()

	log.Info().Msg("waiting for the server to complete")

	// Shutdown server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer shutdownCancel()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to shutdown server")
	}

	log.Info().Msg("the server is shutting down")
}

func toLoggerConfig(cfg appconfig.LoggerConfig) logging.Config {
	return logging.Config{
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
