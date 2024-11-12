package account

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	apihandler "github.com/mandarine-io/Backend/internal/transport/http/handler"
	middleware2 "github.com/mandarine-io/Backend/pkg/transport/http/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type handler struct {
	svc service.AccountService
}

func NewHandler(svc service.AccountService) apihandler.ApiHandler {
	return &handler{svc: svc}
}

func (h *handler) RegisterRoutes(router *gin.Engine, middlewares apihandler.RouteMiddlewares) {
	log.Debug().Msg("register service routes")

	router.GET(
		"v0/account",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.GetAccount,
	)
	router.PATCH(
		"v0/account/username",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.UpdateUsername)
	router.PATCH(
		"v0/account/email",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.UpdateEmail)
	router.POST(
		"v0/account/email/verify",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.VerifyEmail)
	router.POST(
		"v0/account/password",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.SetPassword)
	router.PATCH(
		"v0/account/password",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.UpdatePassword)
	router.DELETE(
		"v0/account",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.DeleteAccount)
	router.GET(
		"v0/account/restore",
		middlewares.Auth,
		middlewares.BannedUser,
		h.RestoreAccount,
	)
}

// GetAccount godoc
//
//	@Id				GetAccount
//	@Summary		Get service
//	@Description	Request for receiving own service. User must be logged in. In response will be returned own service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.AccountOutput	"Account info"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found user"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account [get]
func (h *handler) GetAccount(ctx *gin.Context) {
	log.Debug().Msg("handle get service")
	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.svc.GetAccount(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// UpdateUsername godoc
//
//	@Id				UpdateUsername
//	@Summary		Update username
//	@Description	Request for updating username. User must be logged in. In response will be returned updated service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateUsernameInput	true	"Update username request body"
//	@Success		200		{object}	dto.AccountOutput	"Account info"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404		{object}	dto.ErrorResponse	"Not found user"
//	@Failure		409		{object}	dto.ErrorResponse	"Duplicate username"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account/username [patch]
func (h *handler) UpdateUsername(ctx *gin.Context) {
	log.Debug().Msg("handle update username")

	req := dto.UpdateUsernameInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.svc.UpdateUsername(ctx, principal.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrDuplicateUsername):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// UpdateEmail godoc
//
//	@Id				UpdateEmail
//	@Summary		Update email
//	@Description	Request for updating email. User must be logged in. In process will be sent verification email. In response will be returned updated service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateEmailInput	true	"Update email request body"
//	@Success		200		{object}	dto.AccountOutput	"Account info (email is verified)"
//	@Success		202		{object}	dto.AccountOutput	"Account info (email is not verified)"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404		{object}	dto.ErrorResponse	"Not found user"
//	@Failure		409		{object}	dto.ErrorResponse	"Duplicate email"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account/email [patch]
func (h *handler) UpdateEmail(ctx *gin.Context) {
	log.Debug().Msg("handle update email")

	req := dto.UpdateEmailInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	log.Debug().Msg("get localizer")
	localizer := ctx.Value(middleware2.LocalizerKey).(*i18n.Localizer)

	res, err := h.svc.UpdateEmail(ctx, principal.ID, req, localizer)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrDuplicateEmail):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		case errors.Is(err, service.ErrSendEmail):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	if res.IsEmailVerified {
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusAccepted, res)
	}
}

// VerifyEmail godoc
//
//	@Id				VerifyEmail
//	@Summary		Verify email
//	@Description	Request for verify email. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.VerifyEmailInput	true	"Verify email request body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found user"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account/email/verify [post]
func (h *handler) VerifyEmail(ctx *gin.Context) {
	log.Debug().Msg("handle verify email")

	req := dto.VerifyEmailInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.VerifyEmail(ctx, principal.ID, req); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrExpiredOtp):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// SetPassword godoc
//
//	@Id				SetPassword
//	@Summary		Set password
//	@Description	Request for setting password. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.SetPasswordInput	true	"Set password request body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found user"
//	@Failure		409	{object}	dto.ErrorResponse	"Password is set"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account/password [post]
func (h *handler) SetPassword(ctx *gin.Context) {
	log.Debug().Msg("handle set password")

	req := dto.SetPasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.SetPassword(ctx, principal.ID, req); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrPasswordIsSet):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// UpdatePassword godoc
//
//	@Id				UpdatePassword
//	@Summary		Update password
//	@Description	Request for updating password. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UpdatePasswordInput	true	"Update password request body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found user"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account/password [patch]
func (h *handler) UpdatePassword(ctx *gin.Context) {
	log.Debug().Msg("handle update password")

	req := dto.UpdatePasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.UpdatePassword(ctx, principal.ID, req); err != nil {
		switch {
		case errors.Is(err, service.ErrIncorrectOldPassword):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// RestoreAccount godoc
//
//	@Id				RestoreAccount
//	@Summary		Restore service
//	@Description	Request for restoring service. User must be logged in. User must be deleted. In response will be returned restored service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.AccountOutput	"Account info"
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found user"
//	@Failure		409	{object}	dto.ErrorResponse	"User is not deleted"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account/restore [get]
func (h *handler) RestoreAccount(ctx *gin.Context) {
	log.Debug().Msg("handle restore service")

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.svc.RestoreAccount(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrUserNotDeleted):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// DeleteAccount godoc
//
//	@Id				DeleteAccount
//	@Summary		Delete service
//	@Description	Request for deleting service. User must be logged in. User must not be deleted.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404	{object}	dto.ErrorResponse	"Not found user"
//	@Failure		409	{object}	dto.ErrorResponse	"User is deleted"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/account [delete]
func (h *handler) DeleteAccount(ctx *gin.Context) {
	log.Debug().Msg("handle delete service")

	principal, err := middleware2.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.DeleteAccount(ctx, principal.ID); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrUserAlreadyDeleted):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}
