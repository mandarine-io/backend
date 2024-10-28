package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type AgentMock struct {
	mock.Mock
}

func (a *AgentMock) Publish(ctx context.Context, topic string, msg interface{}) error {
	args := a.Called(ctx, topic, msg)
	return args.Error(0)
}

func (a *AgentMock) Subscribe(ctx context.Context, topic string) (<-chan interface{}, error) {
	args := a.Called(ctx, topic)
	return args.Get(0).(<-chan interface{}), args.Error(1)
}

func (a *AgentMock) Close() error {
	args := a.Called()
	return args.Error(0)
}
