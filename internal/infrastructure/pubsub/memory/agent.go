package memory

import (
	"context"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	"github.com/rs/zerolog"
	"sync"
)

type Option func(*agent) error

func WithLogger(logger zerolog.Logger) Option {
	return func(a *agent) error {
		a.logger = logger
		return nil
	}
}

type agent struct {
	mu     sync.Mutex
	subs   map[string][]chan pubsub.Event
	logger zerolog.Logger
}

func NewAgent(opts ...Option) (pubsub.Agent, error) {
	a := &agent{
		subs:   make(map[string][]chan pubsub.Event),
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return a, nil
}

func (a *agent) Publish(_ context.Context, topic string, msg any) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logger.Debug().Msgf("publish message to topic: %s", topic)

	if _, ok := a.subs[topic]; !ok {
		return pubsub.ErrTopicNotFound
	}

	// Send to subscribers
	event := pubsub.Event{
		Topic:   topic,
		Payload: msg,
	}
	for _, ch := range a.subs[topic] {
		ch <- event
	}

	return nil
}

func (a *agent) Subscribe(_ context.Context, topic string) (<-chan pubsub.Event, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logger.Debug().Msgf("subscribe to topic: %s", topic)

	ch := make(chan pubsub.Event, 1024)
	a.subs[topic] = append(a.subs[topic], ch)

	return ch, nil
}

func (a *agent) Unsubscribe(_ context.Context, topic string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logger.Debug().Msgf("unsubscribe to topic: %s", topic)

	if _, ok := a.subs[topic]; !ok {
		return pubsub.ErrTopicNotFound
	}

	for _, ch := range a.subs[topic] {
		close(ch)
	}
	delete(a.subs, topic)

	return nil
}

func (a *agent) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, ch := range a.subs {
		for _, sub := range ch {
			close(sub)
		}
	}

	return nil
}
