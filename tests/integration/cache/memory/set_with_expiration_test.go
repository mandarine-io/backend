package memory

import (
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type SetWithExpirationSuite struct {
	suite.Suite
}

func (s *SetWithExpirationSuite) AfterEach(t provider.T) {
	t.Title("SetWithExpiration - after each")
	t.Feature("Redis cache manager")

	err := manager.Delete(ctx, "set_with_exp_key")
	t.Require().NoError(err)
}

func (s *SetWithExpirationSuite) Test_Success(t provider.T) {
	t.Title("SetWithExpiration - success")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	err := manager.SetWithExpiration(ctx, "set_with_exp_key", "\"value\"", time.Second)
	t.Require().NoError(err)

	var value string
	err = manager.Get(ctx, "set_with_exp_key", &value)
	t.Require().NoError(err)
	t.Require().Equal("\"value\"", value)

	time.Sleep(4 * time.Second)

	err = manager.Get(ctx, "set_with_exp_key", &value)
	t.Require().Error(err)
	t.Require().ErrorIs(err, cache.ErrCacheEntryNotFound)
}
