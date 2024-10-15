package logging

import (
	slogmulti "github.com/samber/slog-multi"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"

	"github.com/rs/zerolog"
	"github.com/timandy/routine"
	"log/slog"
	"os"
)

const (
	logEncodingText = "text"
	logEncodingJson = "json"
)

type Config struct {
	Console ConsoleLoggerConfig
	File    FileLoggerConfig
}

type ConsoleLoggerConfig struct {
	Level    string
	Encoding string
}

type FileLoggerConfig struct {
	Enable  bool
	Level   string
	DirPath string
	MaxSize int
	MaxAge  int
}

func SetupLogger(cfg *Config) {
	var handlers []slog.Handler
	handlers = append(handlers, setupConsoleHandler(cfg.Console))
	if cfg.File.Enable {
		handlers = append(handlers, setupFileHandler(cfg.File))
	}

	logger := slog.New(slogmulti.Fanout(handlers...))
	slog.SetDefault(logger)
}

func setupConsoleHandler(cfg ConsoleLoggerConfig) slog.Handler {
	// Parse logging level
	level, err := parseLevel(string(cfg.Level))
	if err != nil {
		level = slog.LevelInfo
	}

	// Setup console handler options
	var opts slogzerolog.Option
	switch cfg.Encoding {
	case logEncodingText:
		zerologLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Hook(
			&pidHook{}, &gidHook{},
		).With().Timestamp().Logger()
		opts = slogzerolog.Option{
			Level:  level,
			Logger: &zerologLogger,
		}
	case logEncodingJson:
		zerologLogger := zerolog.New(os.Stdout).Hook(
			&pidHook{}, &gidHook{},
		).With().Timestamp().Logger()
		opts = slogzerolog.Option{
			Level:     level,
			Logger:    &zerologLogger,
			AddSource: true,
		}
	}

	return opts.NewZerologHandler()
}

func setupFileHandler(cfg FileLoggerConfig) slog.Handler {
	// Parse logging level
	level, err := parseLevel(cfg.Level)
	if err != nil {
		level = slog.LevelInfo
	}

	// Setup file handler options
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.DirPath + "/app.log",
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: 3,
	}
	zerologLogger := zerolog.New(lumberjackLogger).Hook(
		&pidHook{}, &gidHook{},
	).With().Timestamp().Logger()
	opts := slogzerolog.Option{
		Level:     level,
		Logger:    &zerologLogger,
		AddSource: true,
	}

	return opts.NewZerologHandler()
}

type pidHook struct{}

func (h *pidHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Int("pid", os.Getpid())
}

type gidHook struct{}

func (h *gidHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Int64("gid", routine.Goid())
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
