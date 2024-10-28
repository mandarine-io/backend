package rest

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mandarine-io/Backend/internal/api/registry"
	validator2 "github.com/mandarine-io/Backend/pkg/rest/validator"
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

func NewServer(container *registry.Container) *http.Server {
	// Setup validators
	log.Debug().Msg("setup validators")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator2.PastDateValidator)
		_ = v.RegisterValidation("zxcvbn", validator2.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator2.UsernameValidator)
	}

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
