package redis

import (
	"context"
	"github.com/mandarine-io/Backend/pkg/pubsub"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"strings"
)

type agent struct {
	rdb *redis.Client
}

func NewAgent(rdb *redis.Client) pubsub.Agent {
	return &agent{rdb: rdb}
}

func (a *agent) Publish(ctx context.Context, topic string, msg interface{}) error {
	log.Debug().Msgf("publish message to topic: %s", topic)
	return a.rdb.Publish(ctx, topic, msg).Err()
}

func (a *agent) Subscribe(ctx context.Context, topics ...string) (<-chan pubsub.Event, error) {
	log.Debug().Msgf("subscribe to topics: %s", strings.Join(topics, ", "))
	p := a.rdb.Subscribe(ctx, topics...)

	eventChan := make(chan pubsub.Event)

	// TODO: remake without goroutine
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

func (a *agent) Close() error {
	return nil
}
