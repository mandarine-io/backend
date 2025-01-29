package scheduler

import (
	"context"
	scheduler2 "github.com/mandarine-io/backend/internal/scheduler"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type CronJobSuite struct {
	suite.Suite
}

func (s *CronJobSuite) Test_Success(t provider.T) {
	t.Title("Cron job - success")
	t.Severity(allure.NORMAL)
	t.Feature("Scheduler")
	t.Tags("Positive")
	t.Parallel()

	resCh := make(chan bool, 10)

	job := scheduler2.Job{
		Ctx:            context.Background(),
		Name:           "test",
		CronExpression: "* * * * * *",
		Action: func(ctx context.Context) error {
			resCh <- true
			return nil
		},
	}

	_, err := scheduler.AddJob(job)
	t.Require().NoError(err)

	scheduler.Start()
	defer func() {
		_ = scheduler.Shutdown()
	}()

	timeoutCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		t.Require().NoError(timeoutCtx.Err())
	case res := <-resCh:
		t.Require().True(res)
	}
}
