package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/service/domain"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/internal/transport/http/middleware"
	"github.com/mandarine-io/backend/internal/transport/http/util"
	_ "github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type handler struct {
	svc    domain.WebsocketService
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.WebsocketService, opts ...Option) apihandler.APIHandler {
	h := &handler{
		svc:    svc,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	log.Debug().Msg("register websocket routes")

	router.GET(
		"v0/ws",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
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
//	@Failure		400	{object}	v0.ErrorOutput
//	@Failure		401	{object}	v0.ErrorOutput
//	@Failure		503	{object}	v0.ErrorOutput
//	@Router			/v0/ws [get]
func (h *handler) Connect(ctx *gin.Context) {
	log.Debug().Msg("handle connect")

	authUser, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	// RegisterClient client in pool
	ctx.Writer.Header().Set("Content-Type", "application/json")
	_ = h.svc.RegisterClient(authUser.ID, ctx.Request, ctx.Writer)
}
