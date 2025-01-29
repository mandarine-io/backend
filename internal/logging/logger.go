package logging

import (
	"github.com/timandy/routine"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logEncodingText = "text"
	logEncodingJSON = "json"
)

type Config struct {
	Level   string
	Console ConsoleLoggerConfig
	File    FileLoggerConfig
}

type ConsoleLoggerConfig struct {
	Enable   bool
	Encoding string
}

type FileLoggerConfig struct {
	Enable  bool
	DirPath string
	MaxSize int
	MaxAge  int
}

func SetupLogger(cfg Config) {
	var writers []io.Writer

	if consoleWriter := setupConsoleWriter(cfg.Console); consoleWriter != nil {
		writers = append(writers, consoleWriter)
	}

	if fileWriter := setupFileWriter(cfg.File); fileWriter != nil {
		writers = append(writers, fileWriter)
	}

	mw := io.MultiWriter(writers...)

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return file + ":" + strconv.Itoa(line)
	}
	zerolog.SetGlobalLevel(parseLevel(cfg.Level))
	log.Logger = zerolog.
		New(mw).
		Hook(pidHook{}, gidHook{}).
		With().
		Timestamp().
		Caller().
		Logger()
}

func setupConsoleWriter(cfg ConsoleLoggerConfig) io.Writer {
	if !cfg.Enable {
		return nil
	}

	if cfg.Encoding == logEncodingJSON {
		return os.Stdout
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	return consoleWriter
}

func setupFileWriter(cfg FileLoggerConfig) io.Writer {
	if cfg.Enable {
		return newRollingFile(cfg)
	}

	return nil
}

func newRollingFile(cfg FileLoggerConfig) io.Writer {
	if err := os.MkdirAll(cfg.DirPath, 0750); err != nil {
		log.Error().Stack().Err(err).Msgf("can't create log directory %s", cfg.DirPath)
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(cfg.DirPath, "app.log"),
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: 3,
	}
}

func parseLevel(s string) zerolog.Level {
	switch s {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

type pidHook struct{}

func (h pidHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Int("pid", os.Getpid())
}

type gidHook struct{}

func (h gidHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Int64("gid", routine.Goid())
}
