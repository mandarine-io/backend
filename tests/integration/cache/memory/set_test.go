package memory

import (
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type SetSuite struct {
	suite.Suite
}

func (s *SetSuite) AfterEach(t provider.T) {
	t.Title("Set - after each")
	t.Feature("Redis cache manager")

	err := manager.Delete(ctx, "set_key")
	t.Require().NoError(err)
}

func (s *SetSuite) Test_Success(t provider.T) {
	t.Title("Set - success")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	err := manager.Set(ctx, "set_key", "\"value\"")
	t.Require().NoError(err)

	var value string
	err = manager.Get(ctx, "set_key", &value)
	t.Require().NoError(err)
	t.Require().Equal("\"value\"", value)

	time.Sleep(2 * time.Second)

	err = manager.Get(ctx, "set_key", &value)
	t.Require().Error(err)
	t.Require().ErrorIs(err, cache.ErrCacheEntryNotFound)
}
