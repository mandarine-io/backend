package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mandarine/internal/api/cli"
	appconfig "mandarine/internal/api/config"
	"mandarine/internal/api/job"
	"mandarine/internal/api/registry"
	"mandarine/internal/api/rest"
	"mandarine/pkg/config"
	"mandarine/pkg/logging"
	"mandarine/pkg/scheduler"
	syshttp "net/http"
	"os"
	"os/signal"
	"syscall"
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
	var loggerCfg appconfig.OnlyLoggerConfig
	config.MustLoadConfig(options.ConfigFilePath, options.EnvFilePath, &loggerCfg)
	logging.SetupLogger(mapAppLoggerConfigToLoggerConfig(&loggerCfg.Logger))

	// Load config
	var cfg appconfig.Config
	config.MustLoadConfig(options.ConfigFilePath, options.EnvFilePath, &cfg)

	// Setup container
	container := registry.MustNewContainer(&cfg)
	defer func() {
		_ = container.Close()
	}()

	// Setup scheduler
	jobs := job.SetupJobs(container)
	cronScheduler := scheduler.MustSetupJobScheduler(jobs)
	defer func() {
		err := cronScheduler.Shutdown()
		if err != nil {
			slog.Error("Job scheduler shutdown error", logging.ErrorAttr(err))
		}
	}()
	cronScheduler.Start()

	// SIGINT and SIGTERM handler
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create server
	srv := rest.NewServer(container)

	// Run server
	slog.Info(fmt.Sprintf("The server listens on port %d", cfg.Server.Port))
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, syshttp.ErrServerClosed) {
			slog.Error("Server error", logging.ErrorAttr(err))
			stop()
		}
	}()

	// Wait for signal
	slog.Info("To stop server press Ctrl+C")
	<-ctx.Done()
	stop()
	slog.Info("Waiting for the server to complete")

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", logging.ErrorAttr(err))
	}

	slog.Info("Server is shutdown")
}

func getEnvWithDefault(envName, defaultValue string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	return defaultValue
}

func mapAppLoggerConfigToLoggerConfig(cfg *appconfig.LoggerConfig) *logging.Config {
	return &logging.Config{
		Console: logging.ConsoleLoggerConfig{
			Level:    cfg.Console.Level,
			Encoding: cfg.Console.Encoding,
		},
		File: logging.FileLoggerConfig{
			Enable:  cfg.File.Enable,
			Level:   cfg.File.Level,
			DirPath: cfg.File.DirPath,
			MaxSize: cfg.File.MaxSize,
			MaxAge:  cfg.File.MaxAge,
		},
	}
}
