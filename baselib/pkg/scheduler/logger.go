package scheduler

import (
	"github.com/rs/zerolog/log"
)

type schedulerLogger struct{}

func (s schedulerLogger) Debug(msg string, args ...any) {
	log.Debug().Interface("args", args).Msg(msg)
}

func (s schedulerLogger) Info(msg string, args ...any) {
	log.Info().Interface("args", args).Msg(msg)
}

func (s schedulerLogger) Warn(msg string, args ...any) {
	log.Warn().Interface("args", args).Msg(msg)
}

func (s schedulerLogger) Error(msg string, args ...any) {
	log.Error().Interface("args", args).Msg(msg)
}
