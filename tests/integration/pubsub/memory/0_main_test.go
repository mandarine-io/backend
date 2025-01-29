package memory

import (
	"context"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub/memory"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	ctx   = context.Background()
	agent pubsub.Agent
)

type MemoryPubSubSuite struct {
	suite.Suite
}

func TestRedisPubSubSuite(t *testing.T) {
	var err error
	agent, err = memory.NewAgent()
	require.NoError(t, err)

	suite.RunSuite(t, new(MemoryPubSubSuite))
}

func (s *MemoryPubSubSuite) Test(t provider.T) {
	s.RunSuite(t, new(PublishSubscribeSuite))
}
