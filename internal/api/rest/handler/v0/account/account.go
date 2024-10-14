package account

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mandarine/internal/api/service/account"
	"mandarine/internal/api/service/account/dto"
	"mandarine/pkg/rest/middleware"
	"net/http"
)

type AccountHandler struct {
	accountService *account.AccountService
}

func NewAccountHandler(accountService *account.AccountService) AccountHandler {
	return AccountHandler{accountService: accountService}
}

// GetAccount godoc
//
//	@Id				GetAccount
//	@Summary		Get account
//	@Description	Request for receiving own account. User must be logged in. In response will be returned own account info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.AccountOutput
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/account [get]
func (h AccountHandler) GetAccount(ctx *gin.Context) {
	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.accountService.GetAccount(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
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
//	@Description	Request for updating username. User must be logged in. In response will be returned updated account info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateUsernameInput	true	"Update username request body"
//	@Success		200		{object}	dto.AccountOutput
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v0/account/username [patch]
func (h AccountHandler) UpdateUsername(ctx *gin.Context) {
	req := dto.UpdateUsernameInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.accountService.UpdateUsername(ctx, principal.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, account.ErrDuplicateUsername):
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
//	@Description	Request for updating email. User must be logged in. In process will be sent verification email. In response will be returned updated account info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateEmailInput	true	"Update email request body"
//	@Success		200		{object}	dto.AccountOutput
//	@Success		202		{object}	dto.AccountOutput
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v0/account/email [patch]
func (h AccountHandler) UpdateEmail(ctx *gin.Context) {
	req := dto.UpdateEmailInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.accountService.UpdateEmail(ctx, principal.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, account.ErrDuplicateEmail):
			_ = ctx.AbortWithError(http.StatusConflict, err)
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
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/account/email/verify [post]
func (h AccountHandler) VerifyEmail(ctx *gin.Context) {
	req := dto.VerifyEmailInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.accountService.VerifyEmail(ctx, principal.ID, req); err != nil {
		switch {
		case errors.Is(err, account.ErrInvalidOrExpiredOtp):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, account.ErrUserNotFound):
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
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/account/password [post]
func (h AccountHandler) SetPassword(ctx *gin.Context) {
	req := dto.SetPasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.accountService.SetPassword(ctx, principal.ID, req); err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, account.ErrPasswordIsSet):
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
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/account/password [patch]
func (h AccountHandler) UpdatePassword(ctx *gin.Context) {
	req := dto.UpdatePasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.accountService.UpdatePassword(ctx, principal.ID, req); err != nil {
		switch {
		case errors.Is(err, account.ErrIncorrectOldPassword):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, account.ErrUserNotFound):
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
//	@Summary		Restore account
//	@Description	Request for restoring account. User must be logged in. User must be deleted. In response will be returned restored account info.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.AccountOutput
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/account/restore [get]
func (h AccountHandler) RestoreAccount(ctx *gin.Context) {
	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	res, err := h.accountService.RestoreAccount(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, account.ErrUserNotDeleted):
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
//	@Summary		Delete account
//	@Description	Request for deleting account. User must be logged in. User must not be deleted.
//	@Security		BearerAuth
//	@Tags			Account API
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		409	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/account [delete]
func (h AccountHandler) DeleteAccount(ctx *gin.Context) {
	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if err := h.accountService.DeleteAccount(ctx, principal.ID); err != nil {
		switch {
		case errors.Is(err, account.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, account.ErrUserAlreadyDeleted):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}
