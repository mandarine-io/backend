package pubsub

import (
	"context"
	"errors"
)

var (
	ErrAgentClosed   = errors.New("agent closed")
	ErrTopicNotFound = errors.New("topic not found")
)

type Event struct {
	Topic   string
	Payload any
}

type Agent interface {
	Publish(ctx context.Context, topic string, msg any) error
	Subscribe(ctx context.Context, topic string) (<-chan Event, error)
	Unsubscribe(ctx context.Context, topic string) error
	Close() error
}
