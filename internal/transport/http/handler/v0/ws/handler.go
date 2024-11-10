package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/transport/http/handler"
	"github.com/mandarine-io/Backend/pkg/transport/http/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	svc service.WebsocketService
}

func NewHandler(svc service.WebsocketService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, middlewares handler.RouteMiddlewares) {
	log.Debug().Msg("register websocket routes")

	router.GET(
		"v0/ws",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.Connect,
	)
}

// Connect godoc
//
//	@Id				WsConnect
//	@Summary		Connect to websocket server
//	@Description	Request for connect to websocket server. If pool is not full, a new websocket connection is created.
//	@Tags			Websocket API
//	@Security		BearerAuth
//	@Success		101
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		503	{object}	dto.ErrorResponse
//	@Router			/v0/ws [get]
func (h *Handler) Connect(ctx *gin.Context) {
	log.Debug().Msg("handle connect")

	authUser, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// RegisterClient client in pool
	_ = h.svc.RegisterClient(authUser.ID, ctx.Request, ctx.Writer)
}
