package memory

import (
	"context"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/infrastructure/cache/memory"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	ctx     = context.Background()
	manager cache.Manager
)

type MemoryCacheManagerSuite struct {
	suite.Suite
}

func TestMemoryCacheManagerSuite(t *testing.T) {
	var err error
	manager, err = memory.NewManager(memory.WithTTL(500 * time.Millisecond))
	require.NoError(t, err)

	suite.RunSuite(t, new(MemoryCacheManagerSuite))
}

func (s *MemoryCacheManagerSuite) Test(t provider.T) {
	s.RunSuite(t, new(DeleteSuite))
	s.RunSuite(t, new(GetSuite))
	s.RunSuite(t, new(InvalidateSuite))
	s.RunSuite(t, new(SetSuite))
	s.RunSuite(t, new(SetWithExpirationSuite))
}
