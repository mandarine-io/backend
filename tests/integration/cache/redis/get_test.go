package redis

import (
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type GetSuite struct {
	suite.Suite
}

func (s *GetSuite) BeforeEach(t provider.T) {
	t.Title("Get - before each")
	t.Feature("Redis cache manager")

	rdb.Set(ctx, "key", "\"value\"", -1)
}

func (s *GetSuite) AfterEach(t provider.T) {
	t.Title("Get - after each")
	t.Feature("Redis cache manager")

	rdb.Del(ctx, "key")
}

func (s *GetSuite) Test_Success(t provider.T) {
	t.Title("Get - success")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	var value string
	err := manager.Get(ctx, "key", &value)

	t.Require().NoError(err)
	t.Require().Equal("value", value)
}

func (s *GetSuite) Test_NoSuchKey(t provider.T) {
	t.Title("Get - no such key")
	t.Severity(allure.CRITICAL)
	t.Feature("Redis cache manager")
	t.Tags("Negative")

	var value string
	err := manager.Get(ctx, "no-such-key", &value)
	t.Require().Error(err)
	t.Require().ErrorIs(err, cache.ErrCacheEntryNotFound)
}
