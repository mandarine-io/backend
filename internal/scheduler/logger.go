package scheduler

import (
	"github.com/rs/zerolog"
)

type schedulerLogger struct {
	logger zerolog.Logger
}

func (s schedulerLogger) Debug(msg string, args ...any) {
	log(s.logger.Debug(), msg, args...)
}

func (s schedulerLogger) Info(msg string, args ...any) {
	log(s.logger.Info(), msg, args...)
}

func (s schedulerLogger) Warn(msg string, args ...any) {
	log(s.logger.Warn(), msg, args...)
}

func (s schedulerLogger) Error(msg string, args ...any) {
	log(s.logger.Error(), msg, args...)
}

func log(event *zerolog.Event, msg string, args ...any) {
	if args == nil || len(args)%2 != 0 {
		event.Msg(msg)
		return
	}

	for i := 0; i < len(args)/2; i++ {
		key, ok := args[i*2].(string)
		if !ok {
			continue
		}

		event = event.Any(key, args[i*2+1])
	}

	event.Msg(msg)
}
