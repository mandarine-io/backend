package rest

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"mandarine/internal/api/registry"
	validator2 "mandarine/pkg/rest/validator"
	"net/http"
	"time"
)

const (
	serverReadTimeout = 30 * time.Second
	serverIdleTimeout = 30 * time.Second
)

func NewServer(container *registry.Container) *http.Server {
	// Setup validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator2.PastDateValidator)
		_ = v.RegisterValidation("zxcvbn", validator2.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator2.UsernameValidator)
	}

	// Create server
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", container.Config.Server.Port),
		Handler:           SetupRouter(container),
		ReadTimeout:       serverReadTimeout,
		ReadHeaderTimeout: serverReadTimeout,
		IdleTimeout:       serverIdleTimeout,
		ErrorLog:          slog.NewLogLogger(container.Logger.Handler(), slog.LevelError),
	}
}
