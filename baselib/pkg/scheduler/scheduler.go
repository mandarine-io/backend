package scheduler

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func MustSetupJobScheduler() *Scheduler {
	scheduler, err := gocron.NewScheduler(
		gocron.WithLogger(schedulerLogger{}),
		gocron.WithLimitConcurrentJobs(10, gocron.LimitModeWait),
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to setup job scheduler")
	}

	return &Scheduler{scheduler}
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) AddJob(job Job) (uuid.UUID, error) {
	j, err := s.scheduler.NewJob(
		gocron.CronJob(job.CronExpression, false),
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
