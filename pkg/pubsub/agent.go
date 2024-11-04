package pubsub

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrAgentClosed   = errors.New("agent closed")
	ErrTopicNotFound = errors.New("topic not found")
)

type Event struct {
	Topic   string
	Payload string
}

type Agent interface {
	Publish(ctx context.Context, topic string, msg interface{}) error
	Subscribe(ctx context.Context, topics ...string) (<-chan Event, error)
	Close() error
}
