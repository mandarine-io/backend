package memory

import (
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type PublishSubscribeSuite struct {
	suite.Suite
}

func (suite *PublishSubscribeSuite) Test_Success(t provider.T) {
	t.Title("Publish and subscribe - success")
	t.Severity(allure.NORMAL)
	t.Feature("Memory pubsub")
	t.Tags("Positive")

	subscriber, err := agent.Subscribe(ctx, "test")
	t.Require().NoError(err)

	err = agent.Publish(ctx, "test", "message")
	t.Require().NoError(err)

	event := <-subscriber
	t.Require().Equal("test", event.Topic)
	t.Require().Equal("message", event.Payload)

	err = agent.Unsubscribe(ctx, "test")
	t.Require().NoError(err)
}

func (suite *PublishSubscribeSuite) Test_NotFoundTopic(t provider.T) {
	t.Title("Publish and subscribe - not found topic")
	t.Severity(allure.CRITICAL)
	t.Feature("Memory pubsub")
	t.Tags("Negative")

	err := agent.Unsubscribe(ctx, "test")
	t.Require().Error(err)
	t.Require().ErrorIs(err, pubsub.ErrTopicNotFound)
}
