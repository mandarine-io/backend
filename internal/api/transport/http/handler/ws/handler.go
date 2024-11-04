package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	wsSvc "github.com/mandarine-io/Backend/internal/api/service/ws"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler"
	"github.com/mandarine-io/Backend/pkg/transport/http/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	svc      *wsSvc.Service
	upgrader websocket.Upgrader
}

func NewHandler(svc *wsSvc.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, middlewares handler.RouteMiddlewares) {
	log.Debug().Msg("register websocket routes")

	router.GET(
		"ws",
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
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		503		{object}	dto.ErrorResponse
//	@Router			/ws [get]
func (h *Handler) Connect(ctx *gin.Context) {
	log.Debug().Msg("handle connect")

	authUser, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// Register client in pool
	_ = h.svc.Register(authUser.ID, ctx.Request, ctx.Writer)
}
