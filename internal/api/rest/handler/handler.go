package handler

import (
	"mandarine/internal/api/config"
	v0 "mandarine/internal/api/rest/handler/v0"
	"mandarine/internal/api/service"
)

type Handlers struct {
	V0 *v0.Handlers
}

func NewHandlers(services *service.Services, cfg *config.Config) *Handlers {
	return &Handlers{
		V0: v0.NewHandlers(services, cfg),
	}
}
