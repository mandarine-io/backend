package plugin

import (
	"context"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
}

func (l Logger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Error(ctx context.Context, msg string, opts ...any) {
	log.Ctx(ctx).Error().Stack().Msgf(msg, opts...)
}

func (l Logger) Warn(ctx context.Context, msg string, opts ...any) {
	log.Ctx(ctx).Warn().Msgf(msg, opts...)
}

func (l Logger) Info(ctx context.Context, msg string, opts ...any) {
	log.Ctx(ctx).Info().Msgf(msg, opts...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	sql, _ := f()
	log.Ctx(ctx).Debug().Dur("elapsed", time.Since(begin)).Msg(sql)

	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("sql error")
	}
}
