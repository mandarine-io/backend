package redis

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

	var (
		cursor uint64
		keys   []string
	)
	for {
		var (
			k   []string
			err error
		)
		k, cursor, err = rdb.Scan(ctx, cursor, "*", 0).Result()
		t.Require().NoError(err)

		keys = append(keys, k...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		err := rdb.Del(ctx, keys...).Err()
		t.Require().NoError(err)
	}
}

func (s *SetWithExpirationSuite) Test_Success(t provider.T) {
	t.Title("SetWithExpiration - success")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	err := manager.SetWithExpiration(ctx, "key", "\"value\"", time.Second)
	t.Require().NoError(err)

	var value string
	err = manager.Get(ctx, "key", &value)
	t.Require().NoError(err)
	t.Require().Equal("\"value\"", value)

	time.Sleep(4 * time.Second)

	err = manager.Get(ctx, "key", &value)
	t.Require().Error(err)
	t.Require().ErrorIs(err, cache.ErrCacheEntryNotFound)
}
