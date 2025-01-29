package finalizer

import (
	"github.com/mandarine-io/backend/internal/di"
)

func PubSub(c *di.Container) di.Finalizer {
	return func() error {
		c.Logger.Debug().Msg("tear down pub/sub agent")

		if c.Infrastructure.PubSubRDB == nil {
			return nil
		}

		err := c.Infrastructure.PubSubRDB.Close()
		if err != nil {
			return err
		}

		return c.Infrastructure.PubSubAgent.Close()
	}
}
