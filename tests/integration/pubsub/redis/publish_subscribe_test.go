package redis

import (
	"context"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"time"
)

type PublishSubscribeSuite struct {
	suite.Suite
}

func (suite *PublishSubscribeSuite) Test_Success(t provider.T) {
	t.Title("Publish and subscribe - success")
	t.Severity(allure.NORMAL)
	t.Feature("Redis pubsub")
	t.Tags("Positive")

	deadlineCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()

	subscriber, err := agent.Subscribe(deadlineCtx, "test")
	t.Require().NoError(err)

	err = agent.Publish(deadlineCtx, "test", "message")
	t.Require().NoError(err)

	select {
	case <-deadlineCtx.Done():
		t.Require().Error(deadlineCtx.Err())
	case event := <-subscriber:
		t.Require().Equal("test", event.Topic)
		t.Require().Equal("message", event.Payload)
	}

	err = agent.Unsubscribe(ctx, "test")
	t.Require().NoError(err)
}

func (suite *PublishSubscribeSuite) Test_NotFoundTopic(t provider.T) {
	t.Title("Publish and subscribe - not found topic")
	t.Severity(allure.CRITICAL)
	t.Feature("Redis pubsub")
	t.Tags("Negative")

	err := agent.Unsubscribe(ctx, "test")
	t.Require().Error(err)
	t.Require().ErrorIs(err, pubsub.ErrTopicNotFound)
}
