package account

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/infrastructure/locale"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/internal/transport/http/middleware"
	"github.com/mandarine-io/backend/internal/transport/http/util"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"net/http"
)

type handler struct {
	svc    domain.AccountService
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.AccountService, opts ...Option) apihandler.APIHandler {
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
	h.logger.Debug().Msg("register account routes")

	accountRouter := router.Group("/v0/account")
	{
		accountRouter.GET(
			"",
			middleware.Registry.Auth,
			h.getAccount,
		)
		accountRouter.PATCH(
			"/username",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			middleware.Registry.DeletedUser,
			h.updateUsername,
		)
		accountRouter.PATCH(
			"/email",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			middleware.Registry.DeletedUser,
			h.updateEmail,
		)
		accountRouter.POST(
			"/email/verify",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			middleware.Registry.DeletedUser,
			h.verifyEmail,
		)
		accountRouter.POST(
			"/password",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			middleware.Registry.DeletedUser,
			h.setPassword,
		)
		accountRouter.PATCH(
			"/password",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			middleware.Registry.DeletedUser,
			h.updatePassword,
		)
		accountRouter.DELETE(
			"",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			h.deleteAccount,
		)
		accountRouter.GET(
			"/restore",
			middleware.Registry.Auth,
			middleware.Registry.BannedUser,
			h.restoreAccount,
		)
	}
}

// getAccount godoc
//
//	@Id				GetAccount
//	@Summary		Get service
//	@Description	Request for receiving own domain. User must be logged in. In response will be returned own service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200	{object}	v0.AccountOutput	"Account info"
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted"
//	@Failure		404	{object}	v0.ErrorOutput	"Not found user"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/account [get]
func (h *handler) getAccount(ctx *gin.Context) {
	h.logger.Debug().Msg("handle get service")
	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	res, err := h.svc.GetAccount(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// updateUsername godoc
//
//	@Id				UpdateUsername
//	@Summary		Update username
//	@Description	Request for updating username. User must be logged in. In response will be returned updated service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body		v0.UpdateUsernameInput	true	"Update username request body"
//	@Success		200		{object}	v0.AccountOutput			"Account info"
//	@Failure		400		{object}	v0.ErrorOutput			"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput			"Unauthorized"
//	@Failure		403		{object}	v0.ErrorOutput			"User is blocked or deleted"
//	@Failure		404		{object}	v0.ErrorOutput			"Not found user"
//	@Failure		409		{object}	v0.ErrorOutput			"Duplicate username"
//	@Failure		500		{object}	v0.ErrorOutput			"Internal server error"
//	@Router			/v0/account/username [patch]
func (h *handler) updateUsername(ctx *gin.Context) {
	h.logger.Debug().Msg("handle update username")

	input := v0.UpdateUsernameInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	res, err := h.svc.UpdateUsername(ctx, principal.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrDuplicateUsername):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// updateEmail godoc
//
//	@Id				UpdateEmail
//	@Summary		Update email
//	@Description	Request for updating email. User must be logged in. In process will be sent verification email. In response will be returned updated service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body		v0.UpdateEmailInput	true	"Update email request body"
//	@Success		200		{object}	v0.AccountOutput		"Account info (email is verified)"
//	@Success		202		{object}	v0.AccountOutput		"Account info (email is not verified)"
//	@Failure		400		{object}	v0.ErrorOutput		"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput		"Unauthorized"
//	@Failure		403		{object}	v0.ErrorOutput		"User is blocked or deleted"
//	@Failure		404		{object}	v0.ErrorOutput		"Not found user"
//	@Failure		409		{object}	v0.ErrorOutput		"Duplicate email"
//	@Failure		500		{object}	v0.ErrorOutput		"Internal server error"
//	@Router			/v0/account/email [patch]
func (h *handler) updateEmail(ctx *gin.Context) {
	h.logger.Debug().Msg("handle update email")

	input := v0.UpdateEmailInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	h.logger.Debug().Msg("get localizer")
	localizer := ctx.Value(middleware.LocalizerKey).(locale.Localizer)

	res, err := h.svc.UpdateEmail(ctx, principal.ID, input, localizer)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrDuplicateEmail):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		case errors.Is(err, domain.ErrSendEmail):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if res.IsEmailVerified {
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(http.StatusAccepted, res)
	}
}

// verifyEmail godoc
//
//	@Id				VerifyEmail
//	@Summary		Verify email
//	@Description	Request for verify email. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.VerifyEmailInput	true	"Verify email request body"
//	@Success		200
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted"
//	@Failure		404	{object}	v0.ErrorOutput	"Not found user"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/account/email/verify [post]
func (h *handler) verifyEmail(ctx *gin.Context) {
	h.logger.Debug().Msg("handle verify email")

	input := v0.VerifyEmailInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.VerifyEmail(ctx, principal.ID, input); err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrInvalidOrExpiredOTP):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// setPassword godoc
//
//	@Id				SetPassword
//	@Summary		Set password
//	@Description	Request for setting password. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.SetPasswordInput	true	"Set password request body"
//	@Success		200
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted"
//	@Failure		404	{object}	v0.ErrorOutput	"Not found user"
//	@Failure		409	{object}	v0.ErrorOutput	"Password is set"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/account/password [post]
func (h *handler) setPassword(ctx *gin.Context) {
	h.logger.Debug().Msg("handle set password")

	input := v0.SetPasswordInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.SetPassword(ctx, principal.ID, input); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrPasswordIsSet):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// updatePassword godoc
//
//	@Id				UpdatePassword
//	@Summary		Update password
//	@Description	Request for updating password. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.UpdatePasswordInput	true	"Update password request body"
//	@Success		200
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted"
//	@Failure		404	{object}	v0.ErrorOutput	"Not found user"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/account/password [patch]
func (h *handler) updatePassword(ctx *gin.Context) {
	h.logger.Debug().Msg("handle update password")

	input := v0.UpdatePasswordInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.UpdatePassword(ctx, principal.ID, input); err != nil {
		switch {
		case errors.Is(err, domain.ErrIncorrectOldPassword):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// restoreAccount godoc
//
//	@Id				RestoreAccount
//	@Summary		Restore service
//	@Description	Request for restoring domain. User must be logged in. User must be deleted. In response will be returned restored service info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200	{object}	v0.AccountOutput	"Account info"
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted"
//	@Failure		404	{object}	v0.ErrorOutput	"Not found user"
//	@Failure		409	{object}	v0.ErrorOutput	"User is not deleted"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/account/restore [get]
func (h *handler) restoreAccount(ctx *gin.Context) {
	h.logger.Debug().Msg("handle restore service")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	res, err := h.svc.RestoreAccount(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrUserNotDeleted):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// deleteAccount godoc
//
//	@Id				DeleteAccount
//	@Summary		Delete service
//	@Description	Request for deleting domain. User must be logged in. User must not be deleted.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		204
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted"
//	@Failure		404	{object}	v0.ErrorOutput	"Not found user"
//	@Failure		409	{object}	v0.ErrorOutput	"User is deleted"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/account [delete]
func (h *handler) deleteAccount(ctx *gin.Context) {
	h.logger.Debug().Msg("handle delete service")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	if err := h.svc.DeleteAccount(ctx, principal.ID); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrUserAlreadyDeleted):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
