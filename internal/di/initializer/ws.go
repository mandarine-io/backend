package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/infrastructure/websocket"
)

func Websocket(c *di.Container) di.Initializer {
	return func() error {
		c.Logger.Debug().Msg("setup ws pool")

		var err error
		c.Infrastructure.WSPool, err = websocket.NewPool(
			c.Config.Websocket.PoolSize,
			websocket.WithLogger(c.Logger.With().Str("component", "ws-pool").Logger()),
		)

		return err
	}
}
