package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"mandarine/internal/api/config"
	"mandarine/internal/api/registry"
	validator3 "mandarine/pkg/rest/validator"
	"net/http"
	"time"
)

const (
	serverReadTimeout     = 30 * time.Second
	serverIdleTimeout     = 30 * time.Second
	serverShutdownTimeout = 5 * time.Second
)

type Server struct {
	server *http.Server
	cfg    *config.Config
}

func NewServer(container *registry.Container) *Server {
	// Setup routes
	router := SetupRouter(container)

	// Setup validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator3.PastDateValidator)
		_ = v.RegisterValidation("zxcvbn", validator3.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator3.UsernameValidator)
	}

	// Create server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", container.Config.Server.Port),
		Handler:           router,
		ErrorLog:          slog.NewLogLogger(container.Logger.Handler(), slog.LevelError),
		ReadTimeout:       serverReadTimeout,
		ReadHeaderTimeout: serverReadTimeout,
		IdleTimeout:       serverIdleTimeout,
	}

	return &Server{
		server: server,
		cfg:    container.Config,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	cancelCtx, cancel := context.WithTimeout(ctx, serverShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(cancelCtx)
}
