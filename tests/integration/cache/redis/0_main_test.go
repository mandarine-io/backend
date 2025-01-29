package redis

import (
	"context"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	redis2 "github.com/mandarine-io/backend/internal/infrastructure/cache/redis"
	"github.com/mandarine-io/backend/tests/integration"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	ctx     = context.Background()
	rdb     redis.UniversalClient
	manager cache.Manager
)

type RedisCacheManagerSuite struct {
	suite.Suite
}

func TestRedisCacheManagerSuite(t *testing.T) {
	var err error
	rdb, err = redis2.NewClient(
		integration.Cfg.GetRedisConfig(),
	)
	require.NoError(t, err)

	manager, err = redis2.NewManager(rdb, redis2.WithTTL(500*time.Millisecond))
	require.NoError(t, err)

	suite.RunSuite(t, new(RedisCacheManagerSuite))
}

func (s *RedisCacheManagerSuite) Test(t provider.T) {
	s.RunSuite(t, new(DeleteSuite))
	s.RunSuite(t, new(GetSuite))
	s.RunSuite(t, new(InvalidateSuite))
	s.RunSuite(t, new(SetSuite))
	s.RunSuite(t, new(SetWithExpirationSuite))
}
