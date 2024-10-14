package scheduler

import "log/slog"

type schedulerLogger struct{}

func (s schedulerLogger) Debug(msg string, args ...any) {
	slog.Debug(msg, slog.Any("args", args))
}

func (s schedulerLogger) Info(msg string, args ...any) {
	slog.Info(msg, slog.Any("args", args))
}

func (s schedulerLogger) Warn(msg string, args ...any) {
	slog.Warn(msg, slog.Any("args", args))
}

func (s schedulerLogger) Error(msg string, args ...any) {
	slog.Error(msg, slog.Any("args", args))
}
