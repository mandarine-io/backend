package ws

import (
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/infrastructure/websocket"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type svc struct {
	pool   *websocket.Pool
	logger zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(pool *websocket.Pool, opts ...Option) domain.WebsocketService {
	s := &svc{
		pool:   pool,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *svc) RegisterClient(userID uuid.UUID, r *http.Request, w http.ResponseWriter) error {
	s.logger.Info().Msg("register websocket client")

	err := s.pool.Register(userID.String(), r, w)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register websocket client")
		return err
	}

	return nil
}
