package websocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
	"time"
)

const (
	pingPeriod = 30 * time.Second
	writeWait  = 1 * time.Minute
	readWait   = 1 * time.Minute
)

var (
	ErrPoolIsFull     = fmt.Errorf("pool is full")
	ErrClientNotFound = fmt.Errorf("client not found")
)

type Option func(*Pool) error

func WithLogger(logger zerolog.Logger) Option {
	return func(pool *Pool) error {
		pool.logger = logger
		return nil
	}
}

type Handler func(pool *Pool, msg ClientMessage)

type Pool struct {
	upgrader    websocket.Upgrader
	logger      zerolog.Logger
	conns       map[string]*websocket.Conn
	handlers    []Handler
	msgCh       chan ClientMessage
	broadcastCh chan BroadcastMessage
	size        int

	mu     sync.RWMutex
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewPool(size int, opts ...Option) (*Pool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &Pool{
		upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			Error:             handleError,
			EnableCompression: true,
		},
		logger:      zerolog.Nop(),
		handlers:    make([]Handler, 0),
		conns:       make(map[string]*websocket.Conn),
		msgCh:       make(chan ClientMessage, 1024),
		broadcastCh: make(chan BroadcastMessage, 1024),
		size:        size,
		cancel:      cancel,
	}

	for _, opt := range opts {
		if err := opt(pool); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	pool.wg.Add(3)
	go pool.sendPingMessages(ctx)
	go pool.sendClientMessages()
	go pool.sendBroadcastMessages()

	return pool, nil
}

func handleError(w http.ResponseWriter, r *http.Request, status int, reason error) {
	w.WriteHeader(status)

	errorResp := v0.NewErrorOutput(reason.Error(), status, r.URL.Path)
	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (p *Pool) Register(id string, r *http.Request, w http.ResponseWriter) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger.Debug().Msgf("register client %s", id)

	if len(p.conns) >= p.size {
		handleError(w, r, http.StatusServiceUnavailable, ErrPoolIsFull)
		return ErrPoolIsFull
	}

	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	conn.SetPingHandler(
		func(string) error {
			_ = conn.WriteMessage(websocket.PongMessage, []byte("pong"))
			return nil
		},
	)
	conn.SetPongHandler(
		func(string) error {
			_ = conn.SetReadDeadline(time.Now().Add(readWait))
			return nil
		},
	)
	conn.SetCloseHandler(
		func(int, string) error {
			p.logger.Debug().Msgf("close client %s", id)
			_ = p.Unregister(id)
			return nil
		},
	)

	p.conns[id] = conn

	go p.receiveClientMessages(id)

	return nil
}

func (p *Pool) Unregister(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger.Debug().Msgf("unregister client %s", id)

	conn, ok := p.conns[id]
	if !ok {
		return ErrClientNotFound
	}

	delete(p.conns, id)

	err := conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

func (p *Pool) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.conns)
}

func (p *Pool) RegisterHandler(h Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, h)
}

func (p *Pool) Send(clientID string, msg []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.logger.Debug().Msg("send client message")
	p.msgCh <- NewClientMessage(clientID, msg)
}

func (p *Pool) Broadcast(msg []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.logger.Debug().Msg("send broadcast message")
	p.broadcastCh <- NewBroadcastMessage(msg)
}

func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Close all connections
	var errs []error
	for clientID, conn := range p.conns {
		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}

		delete(p.conns, clientID)
	}
	p.logger.Debug().Msg("all websocket connections are closed")

	// Close channels
	close(p.msgCh)
	close(p.broadcastCh)

	// Cancel context
	p.cancel()

	p.wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (p *Pool) sendPingMessages(ctx context.Context) {
	p.logger.Debug().Msg("start sending ping messages")

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.wg.Done()
		p.logger.Debug().Msg("ping message sender is stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for clientID, conn := range p.conns {
				_ = conn.SetWriteDeadline(time.Now().Add(writeWait))

				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					p.logger.Error().Stack().Err(err).Msg("failed to send ping message")
					_ = p.Unregister(clientID)
				}
			}
		}
	}
}

func (p *Pool) sendClientMessages() {
	p.logger.Debug().Msg("start sending client messages")
	defer func() {
		p.wg.Done()
		p.logger.Debug().Msg("client message sender is stopped")
	}()

	for clientMsg := range p.msgCh {
		conn, ok := p.conns[clientMsg.ClientID]
		if !ok {
			continue
		}

		_ = conn.SetWriteDeadline(time.Now().Add(writeWait))

		err := conn.WriteMessage(websocket.TextMessage, clientMsg.Payload)
		if err != nil {
			p.logger.Error().Stack().Err(err).Msg("failed to send client message")
			_ = p.Unregister(clientMsg.ClientID)
		}
	}
}

func (p *Pool) sendBroadcastMessages() {
	p.logger.Debug().Msg("start sending broadcast messages")
	defer func() {
		p.wg.Done()
		p.logger.Debug().Msg("broadcast message sender is stopped")
	}()

	for broadcastMsg := range p.broadcastCh {
		for clientID, conn := range p.conns {
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))

			err := conn.WriteMessage(websocket.TextMessage, broadcastMsg.Payload)
			if err != nil {
				p.logger.Error().Stack().Err(err).Msg("failed to send broadcast message")
				_ = p.Unregister(clientID)
			}
		}
	}
}

func (p *Pool) receiveClientMessages(clientID string) {
	for {
		conn, ok := p.conns[clientID]
		if !ok {
			break
		}

		_ = conn.SetReadDeadline(time.Now().Add(readWait))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// During normal close connection `conn.ReadMessage` returns error with code `websocket.CloseNormalClosure`
			// To don`t print p.logger in this case we check that error is `websocket.CloseNormalClosure`
			var wsErr *websocket.CloseError
			if errors.As(err, &wsErr) && wsErr.Code == websocket.CloseNormalClosure {
				break
			}

			p.logger.Error().Stack().Err(err).Msg("failed to receive client message")
			_ = p.Unregister(clientID)
			break
		}

		clientMsg := NewClientMessage(clientID, msg)
		for _, h := range p.handlers {
			h(p, clientMsg)
		}
	}
}
