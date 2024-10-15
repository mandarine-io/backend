package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"mandarine/pkg/logging"
	"os"
)

type Job struct {
	Name       string
	Definition gocron.JobDefinition
	Task       gocron.Task
}

func MustSetupJobScheduler(jobs []Job) gocron.Scheduler {
	// Create scheduler
	scheduler, err := gocron.NewScheduler(gocron.WithLogger(schedulerLogger{}))
	if err != nil {
		slog.Error("Job scheduler setup error", logging.ErrorAttr(err))
		os.Exit(1)
	}

	// Add jobs
	for _, j := range jobs {
		_, err1 := scheduler.NewJob(j.Definition, j.Task)
		if err1 == nil {
			slog.Info("gocron: job added: " + j.Name)
		} else {
			slog.Error("gocron: job error: "+j.Name, logging.ErrorAttr(err1))
		}
	}

	return scheduler
}
