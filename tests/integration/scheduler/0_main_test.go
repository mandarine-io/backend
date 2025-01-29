package scheduler

import (
	scheduler2 "github.com/mandarine-io/backend/internal/scheduler"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	scheduler *scheduler2.Scheduler
)

type SchedulerSuite struct {
	suite.Suite
}

func TestSchedulerSuite(t *testing.T) {
	var err error
	scheduler, err = scheduler2.NewScheduler()
	require.NoError(t, err)

	suite.RunSuite(t, new(SchedulerSuite))
}

func (s *SchedulerSuite) Test(t provider.T) {
	s.RunSuite(t, new(CronJobSuite))
}
