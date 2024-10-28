package resource

import (
	"context"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"time"
)

type dbLogger struct {
}

func (l dbLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l dbLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Error().Stack().Msgf(msg, opts...)
}

func (l dbLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Warn().Msgf(msg, opts...)
}

func (l dbLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	log.Ctx(ctx).Info().Msgf(msg, opts...)
}

func (l dbLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	sql, _ := f()
	log.Ctx(ctx).Debug().Dur("elapsed", time.Since(begin)).Msg(sql)

	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("sql error")
	}
}
