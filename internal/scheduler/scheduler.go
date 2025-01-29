package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Job struct {
	Ctx            context.Context
	Name           string
	CronExpression string
	Action         func(context.Context) error
}

type Scheduler struct {
	scheduler gocron.Scheduler
}

func NewScheduler() (*Scheduler, error) {
	return NewSchedulerWithLogger(zerolog.Nop())
}

func NewSchedulerWithLogger(logger zerolog.Logger) (*Scheduler, error) {
	gocronLogger := schedulerLogger{logger}
	scheduler, err := gocron.NewScheduler(
		gocron.WithLogger(gocronLogger),
		gocron.WithLimitConcurrentJobs(10, gocron.LimitModeWait),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Scheduler{scheduler}, nil
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) AddJob(job Job) (uuid.UUID, error) {
	j, err := s.scheduler.NewJob(
		gocron.CronJob(job.CronExpression, true),
		gocron.NewTask(job.Action, job.Ctx),
		gocron.WithName(job.Name),
	)
	if j == nil {
		return uuid.Nil, err
	}
	return j.ID(), err
}

func (s *Scheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}
