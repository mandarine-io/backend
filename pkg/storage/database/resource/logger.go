package resource

import (
	"context"
	"gorm.io/gorm/logger"
	"log/slog"
	"mandarine/pkg/logging"
	"time"
)

type dbLogger struct {
	level logger.LogLevel
}

func (l dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l dbLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	if l.level >= logger.Error {
		slog.ErrorContext(ctx, "Database error", logging.ErrorStringAttr(msg), slog.Any("opts", opts))
	}
}

func (l dbLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	if l.level >= logger.Warn {
		slog.WarnContext(ctx, msg, slog.Any("opts", opts))
	}
}

func (l dbLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	if l.level >= logger.Info {
		slog.InfoContext(ctx, msg, slog.Any("opts", opts))
	}
}

func (l dbLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	var args []any

	args = append(args, slog.Duration("elapsed", time.Since(begin)))

	sql, rows := f()
	if sql != "" {
		args = append(args, slog.String("sql", sql))
	}
	if rows > -1 {
		args = append(args, slog.Int64("rows", rows))
	}

	slog.DebugContext(ctx, "", args...)
}
