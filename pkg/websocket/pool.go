package websocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/rs/zerolog/log"
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

	errorHandler = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Error().Stack().Err(reason).Msg("failed to encode error response")

		w.WriteHeader(status)

		errorResp := dto.NewErrorResponse(reason.Error(), status, r.URL.Path)
		if err := json.NewEncoder(w).Encode(errorResp); err != nil {
			log.Error().Stack().Err(err).Msg("failed to encode error response")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
)

type Handler func(msg ClientMessage)

type Pool struct {
	upgrader    websocket.Upgrader
	conns       *sync.Map
	handlers    []Handler
	msgCh       chan ClientMessage
	broadcastCh chan BroadcastMessage
	size        int

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewPool(size int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &Pool{
		upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			Error:             errorHandler,
			EnableCompression: true,
		},
		handlers:    make([]Handler, 0),
		conns:       &sync.Map{},
		msgCh:       make(chan ClientMessage),
		broadcastCh: make(chan BroadcastMessage),
		size:        size,
		cancel:      cancel,
	}

	pool.wg.Add(3)
	go pool.sendPingMessages(ctx)
	go pool.sendClientMessages()
	go pool.sendBroadcastMessages()

	return pool
}

func (p *Pool) Register(id string, r *http.Request, w http.ResponseWriter) error {
	log.Debug().Msgf("register client %s", id)

	if lenSyncMap(p.conns) >= p.size {
		errorHandler(w, r, http.StatusServiceUnavailable, ErrPoolIsFull)
		return ErrPoolIsFull
	}

	conn, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	conn.SetPingHandler(func(string) error {
		log.Debug().Msgf("ping client %s", id)
		_ = conn.WriteMessage(websocket.PongMessage, []byte("pong"))
		return nil
	})
	conn.SetPongHandler(func(string) error {
		log.Debug().Msgf("pong client %s", id)
		_ = conn.SetReadDeadline(time.Now().Add(readWait))
		return nil
	})
	conn.SetCloseHandler(func(int, string) error {
		log.Debug().Msgf("close client %s", id)
		_ = p.Unregister(id)
		return nil
	})

	p.conns.Store(id, conn)

	go p.receiveClientMessages(id)

	return nil
}

func (p *Pool) Unregister(id string) error {
	log.Debug().Msgf("unregister client %s", id)

	conn, ok := p.conns.LoadAndDelete(id)
	if !ok {
		return ErrClientNotFound
	}

	return conn.(*websocket.Conn).Close()
}

func (p *Pool) Count() int {
	return lenSyncMap(p.conns)
}

func (p *Pool) RegisterHandler(h Handler) {
	p.handlers = append(p.handlers, h)
}

func (p *Pool) Send(clientId string, msg []byte) {
	log.Debug().Msg("send client message")
	p.msgCh <- NewClientMessage(clientId, msg)
}

func (p *Pool) Broadcast(msg []byte) {
	log.Debug().Msg("send broadcast message")
	p.broadcastCh <- NewBroadcastMessage(msg)
}

func (p *Pool) Close() error {
	// Close all connections
	var errs []error
	p.conns.Range(func(k, v interface{}) bool {
		clientId := k.(string)
		conn := v.(*websocket.Conn)

		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
		p.conns.Delete(clientId)
		return true
	})
	log.Debug().Msg("all websocket connections are closed")

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
	log.Debug().Msg("start sending ping messages")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.wg.Done()
		log.Debug().Msg("ping message sender is stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.conns.Range(func(k, v interface{}) bool {
				clientId := k.(string)
				conn := v.(*websocket.Conn)

				_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Error().Stack().Err(err).Msg("failed to send ping message")
					_ = p.Unregister(clientId)
				}

				return true
			})
		}
	}
}

func (p *Pool) sendClientMessages() {
	log.Debug().Msg("start sending client messages")
	defer func() {
		p.wg.Done()
		log.Debug().Msg("client message sender is stopped")
	}()

	for {
		clientMsg, ok := <-p.msgCh
		if !ok {
			return
		}

		connAny, ok := p.conns.Load(clientMsg.ClientId)
		if !ok {
			continue
		}
		conn := connAny.(*websocket.Conn)

		_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
		err := conn.WriteMessage(websocket.TextMessage, clientMsg.Payload)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to send client message")
			_ = p.Unregister(clientMsg.ClientId)
		}
	}
}

func (p *Pool) sendBroadcastMessages() {
	log.Debug().Msg("start sending broadcast messages")
	defer func() {
		p.wg.Done()
		log.Debug().Msg("broadcast message sender is stopped")
	}()

	for {
		broadcastMsg, ok := <-p.broadcastCh
		if !ok {
			break
		}

		p.conns.Range(func(k, v interface{}) bool {
			clientId := k.(string)
			conn := v.(*websocket.Conn)

			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteMessage(websocket.TextMessage, broadcastMsg.Payload)
			if err != nil {
				log.Error().Stack().Err(err).Msg("failed to send broadcast message")
				_ = p.Unregister(clientId)
			}

			return true
		})
	}
}

func (p *Pool) receiveClientMessages(clientId string) {
	for {
		connAny, ok := p.conns.Load(clientId)
		if !ok {
			break
		}
		conn := connAny.(*websocket.Conn)

		_ = conn.SetReadDeadline(time.Now().Add(readWait))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// During normal close connection `conn.ReadMessage` returns error with code `websocket.CloseNormalClosure`
			// To don`t print log in this case we check that error is `websocket.CloseNormalClosure`
			var wsErr *websocket.CloseError
			if errors.As(err, &wsErr) && wsErr.Code == websocket.CloseNormalClosure {
				break
			}

			log.Error().Stack().Err(err).Msg("failed to receive client message")
			_ = p.Unregister(clientId)
			break
		}

		clientMsg := NewClientMessage(clientId, msg)
		for _, h := range p.handlers {
			h(clientMsg)
		}
	}
}

func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}
