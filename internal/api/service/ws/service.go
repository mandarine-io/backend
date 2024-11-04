package ws

import (
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/pkg/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Service struct {
	Pool *websocket.Pool
}

func NewService(pool *websocket.Pool) *Service {
	return &Service{Pool: pool}
}

func (s *Service) Register(userId uuid.UUID, r *http.Request, w http.ResponseWriter) error {
	log.Info().Msg("register websocket client")
	factoryErr := func(err error) error {
		log.Error().Stack().Err(err).Msg("failed to register websocket client")
		return err
	}

	err := s.Pool.Register(userId.String(), r, w)
	if err != nil {
		return factoryErr(err)
	}

	return nil
}
