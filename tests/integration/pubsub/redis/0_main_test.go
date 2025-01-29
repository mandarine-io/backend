package redis

import (
	"context"
	redis2 "github.com/mandarine-io/backend/internal/infrastructure/cache/redis"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	redis3 "github.com/mandarine-io/backend/internal/infrastructure/pubsub/redis"
	"github.com/mandarine-io/backend/tests/integration"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	ctx   = context.Background()
	rdb   redis.UniversalClient
	agent pubsub.Agent
)

type RedisPubSubSuite struct {
	suite.Suite
}

func TestRedisPubSubSuite(t *testing.T) {
	var err error
	rdb, err = redis2.NewClient(
		integration.Cfg.GetRedisConfig(),
	)
	require.NoError(t, err)

	agent, err = redis3.NewAgent(rdb)
	require.NoError(t, err)

	suite.RunSuite(t, new(RedisPubSubSuite))
}

func (s *RedisPubSubSuite) Test(t provider.T) {
	s.RunSuite(t, new(PublishSubscribeSuite))
}
