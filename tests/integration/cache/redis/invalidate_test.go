package redis

import (
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type InvalidateSuite struct {
	suite.Suite
}

func (s *InvalidateSuite) BeforeEach(t provider.T) {
	t.Title("Invalidate - before each")
	t.Feature("Redis cache manager")

	rdb.Set(ctx, "key", "\"value\"", -1)
}

func (s *InvalidateSuite) AfterEach(t provider.T) {
	t.Title("Invalidate - after each")
	t.Feature("Redis cache manager")

	rdb.Del(ctx, "key")
}

func (s *InvalidateSuite) Test_SuccessWithIncludeRegex(t provider.T) {
	t.Title("Invalidate - success with include regex")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	var value string
	err := manager.Get(ctx, "key", &value)
	t.Require().NoError(err)
	t.Require().Equal("value", value)

	err = manager.Invalidate(ctx, "ke*")
	t.Require().NoError(err)

	err = manager.Get(ctx, "key", &value)
	t.Require().Error(err)
	t.Require().ErrorIs(err, cache.ErrCacheEntryNotFound)
}

func (s *InvalidateSuite) Test_SuccessWithExcludeRegex(t provider.T) {
	t.Title("Invalidate - success with exclude regex")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	var value string
	err := manager.Get(ctx, "key", &value)
	t.Require().NoError(err)
	t.Require().Equal("value", value)

	err = manager.Invalidate(ctx, "ke*1")
	t.Require().NoError(err)

	err = manager.Get(ctx, "key", &value)
	t.Require().NoError(err)
	t.Require().Equal("value", value)
}
