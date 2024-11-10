package ws

import (
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/pkg/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

type svc struct {
	Pool *websocket.Pool
}

func NewService(pool *websocket.Pool) service.WebsocketService {
	return &svc{Pool: pool}
}

func (s *svc) RegisterClient(userId uuid.UUID, r *http.Request, w http.ResponseWriter) error {
	log.Info().Msg("register websocket client")

	err := s.Pool.Register(userId.String(), r, w)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to register websocket client")
		return err
	}

	return nil
}
