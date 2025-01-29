package memory

import (
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type DeleteSuite struct {
	suite.Suite
}

func (s *DeleteSuite) BeforeEach(t provider.T) {
	t.Title("Delete - before each")
	t.Feature("Redis cache manager")

	err := manager.SetWithExpiration(ctx, "delete_key", "value", -1)
	t.Require().NoError(err)
}

func (s *DeleteSuite) AfterEach(t provider.T) {
	t.Title("Delete - after each")
	t.Feature("Redis cache manager")

	err := manager.Delete(ctx, "delete_key")
	t.Require().NoError(err)
}

func (s *DeleteSuite) Test_Success(t provider.T) {
	t.Title("Delete - success")
	t.Severity(allure.NORMAL)
	t.Feature("Redis cache manager")
	t.Tags("Positive")

	err := manager.Delete(ctx, "delete_key")
	t.Require().NoError(err)
}

func (s *DeleteSuite) Test_NoSuchKey(t provider.T) {
	t.Title("Delete - no such delete_key")
	t.Severity(allure.CRITICAL)
	t.Feature("Redis cache manager")
	t.Tags("Negative")

	err := manager.Delete(ctx, "no-such-delete_key")
	t.Require().NoError(err)
}
