package http

import (
	"fmt"
	"github.com/mandarine-io/backend/internal/di"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"log/slog"
	"net/http"
	"time"
)

const (
	serverReadTimeout = 30 * time.Second
	serverIdleTimeout = 30 * time.Second
)

func NewServer(container *di.Container) *http.Server {
	// Create server
	log.Debug().Msg("create server")
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", container.Config.Server.Port),
		Handler:           SetupRouter(container),
		ReadTimeout:       serverReadTimeout,
		ReadHeaderTimeout: serverReadTimeout,
		IdleTimeout:       serverIdleTimeout,
		ErrorLog: slog.NewLogLogger(
			slog.New(
				slogzerolog.Option{
					Level:  slog.LevelDebug,
					Logger: &log.Logger,
				}.NewZerologHandler(),
			).Handler(),
			slog.LevelError,
		),
	}
}
