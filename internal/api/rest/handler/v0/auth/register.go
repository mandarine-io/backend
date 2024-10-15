package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	"net/http"
)

type RegisterHandler struct {
	registerService *auth.RegisterService
}

func NewRegisterHandler(registerService *auth.RegisterService) RegisterHandler {
	return RegisterHandler{
		registerService: registerService,
	}
}

// Register godoc
//
//	@Id				Register
//	@Summary		Register
//	@Description	Request for creating new user. At the end will be sent confirmation email with code
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.RegisterInput	true	"Register request body"
//	@Success		202
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/register [post]
func (h RegisterHandler) Register(ctx *gin.Context) {
	req := dto.RegisterInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.registerService.Register(ctx, req); err != nil {
		switch {
		case errors.Is(err, auth.ErrDuplicateUser):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// RegisterConfirm godoc
//
//	@Id				RegisterConfirm
//	@Summary		Register confirmation
//	@Description	Request for confirming registration. At the end will be created new user
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.RegisterConfirmInput	true	"Register confirm body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/register/confirm [post]
func (h RegisterHandler) RegisterConfirm(ctx *gin.Context) {
	req := dto.RegisterConfirmInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.registerService.RegisterConfirm(ctx, req); err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidOrExpiredOtp):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, auth.ErrDuplicateUser):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}
