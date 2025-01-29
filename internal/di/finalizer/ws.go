package finalizer

import (
	"github.com/mandarine-io/backend/internal/di"
)

func Websocket(c *di.Container) di.Finalizer {
	return func() error {
		c.Logger.Debug().Msg("tear down ws pool")

		if c.Infrastructure.WSPool == nil {
			return nil
		}

		return c.Infrastructure.WSPool.Close()
	}
}
