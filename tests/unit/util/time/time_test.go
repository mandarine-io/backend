package time

import (
	time2 "github.com/mandarine-io/backend/internal/util/time"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
	"time"
)

type TimeUtilSuite struct {
	suite.Suite
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TimeUtilSuite))
}

func (s *TimeUtilSuite) Test_FormatDuration(t provider.T) {
	t.Title("Format duration")
	t.Severity(allure.NORMAL)
	t.Tag("positive")
	t.Parallel()

	type testcase struct {
		Name              string
		Duration          time.Duration
		FormattedDuration string
	}

	datas := []testcase{
		{"00:00:00->00h00m", 0 * time.Second, "00h00m"},
		{"00:00:01->00h00m", 1 * time.Second, "00h00m"},
		{"00:00:59->00h00m", 59 * time.Second, "00h00m"},
		{"00:01:00->00h01m", 60 * time.Second, "00h01m"},
		{"160:00:00->160h00m", 160 * time.Hour, "160h00m"},
	}

	for _, data := range datas {
		t.Run(
			data.Name, func(t provider.T) {
				t.Severity(allure.NORMAL)
				t.Tag("positive")
				t.Parallel()

				actual := time2.FormatDuration(data.Duration)
				t.Require().Equal(data.FormattedDuration, actual)
			},
		)
	}
}
