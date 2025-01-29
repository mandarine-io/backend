package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/pubsub"
	"github.com/redis/go-redis/v9"
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
	rdb    redis.UniversalClient
	logger zerolog.Logger

	mu   sync.Mutex
	subs map[string][]*redis.PubSub
}

func NewAgent(rdb redis.UniversalClient, opts ...Option) (pubsub.Agent, error) {
	a := &agent{
		rdb:    rdb,
		logger: zerolog.Nop(),
		subs:   make(map[string][]*redis.PubSub),
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return a, nil
}

func (a *agent) Publish(ctx context.Context, topic string, msg any) error {
	a.logger.Debug().Msgf("publish message to topic: %s", topic)
	return a.rdb.Publish(ctx, topic, msg).Err()
}

func (a *agent) Subscribe(ctx context.Context, topic string) (<-chan pubsub.Event, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logger.Debug().Msgf("subscribe to topics: %s", topic)
	p := a.rdb.Subscribe(ctx, topic)
	a.subs[topic] = append(a.subs[topic], p)

	eventChan := make(chan pubsub.Event, 1024)

	go func() {
		defer close(eventChan)
		for msg := range p.Channel() {
			event := pubsub.Event{
				Topic:   msg.Channel,
				Payload: msg.Payload,
			}
			eventChan <- event
		}
	}()

	return eventChan, nil

}

func (a *agent) Unsubscribe(ctx context.Context, topic string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.logger.Debug().Msgf("unsubscribe from topic: %s", topic)
	if _, ok := a.subs[topic]; !ok {
		return pubsub.ErrTopicNotFound
	}

	errs := make([]error, 0)
	for _, p := range a.subs[topic] {
		err := p.Unsubscribe(ctx, topic)
		if err != nil {
			errs = append(errs, err)
		}

		err = p.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	delete(a.subs, topic)

	return errors.Join(errs...)
}

func (a *agent) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	errs := make([]error, 0)
	for topic, sub := range a.subs {
		for _, p := range sub {
			err := p.Unsubscribe(context.Background(), topic)
			if err != nil {
				errs = append(errs, err)
			}

			err = p.Close()
			if err != nil {
				errs = append(errs, err)
			}
		}
		delete(a.subs, topic)
	}

	return errors.Join(errs...)
}
