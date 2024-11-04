package memory

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/mandarine-io/Backend/pkg/pubsub"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

type agent struct {
	mu     sync.Mutex
	subs   map[string][]chan pubsub.Event
	quit   chan struct{}
	closed bool
}

func NewAgent() pubsub.Agent {
	return &agent{
		subs: make(map[string][]chan pubsub.Event),
		quit: make(chan struct{}),
	}
}

func (a *agent) Publish(_ context.Context, topic string, msg interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	log.Debug().Msgf("publish message to topic: %s", topic)

	// Checks
	if a.closed {
		return pubsub.ErrAgentClosed
	}
	if _, ok := a.subs[topic]; !ok {
		return pubsub.ErrTopicNotFound
	}

	// Convert to JSON
	var payload []byte
	err := json.Unmarshal(payload, msg)
	if err != nil {
		return err
	}

	// Send to subscribers
	event := pubsub.Event{
		Topic:   topic,
		Payload: string(payload),
	}
	for _, ch := range a.subs[topic] {
		ch <- event
	}

	return nil
}

func (a *agent) Subscribe(_ context.Context, topics ...string) (<-chan pubsub.Event, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	log.Debug().Msgf("subscribe to topic: %s", strings.Join(topics, ", "))

	if a.closed {
		return nil, pubsub.ErrAgentClosed
	}

	ch := make(chan pubsub.Event)
	for _, topic := range topics {
		a.subs[topic] = append(a.subs[topic], ch)
	}
	return ch, nil
}

func (a *agent) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return pubsub.ErrAgentClosed
	}

	a.closed = true
	close(a.quit)

	for _, ch := range a.subs {
		for _, sub := range ch {
			close(sub)
		}
	}

	return nil
}
