package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	"net/http"
)

type ResetPasswordHandler struct {
	resetPasswordService *auth.ResetPasswordService
}

func NewResetPasswordHandler(resetPasswordService *auth.ResetPasswordService) ResetPasswordHandler {
	return ResetPasswordHandler{
		resetPasswordService: resetPasswordService,
	}
}

// RecoveryPassword godoc
//
//	@Id				RecoveryPassword
//	@Summary		Recovery password
//	@Description	Request for recovery password. At the end will be sent email with code
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.RecoveryPasswordInput	true	"Recovery password body"
//	@Success		202
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/recovery-password [post]
func (h ResetPasswordHandler) RecoveryPassword(ctx *gin.Context) {
	req := dto.RecoveryPasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.resetPasswordService.RecoveryPassword(ctx, req); err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusAccepted)
}

// VerifyRecoveryCode godoc
//
//	@Id				VerifyRecoveryCode
//	@Summary		Verify recovery code
//	@Description	Request for verify recovery code. If code is correct will be sent status 200
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.VerifyRecoveryCodeInput	true	"Verify recovery code body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/recovery-password/verify [post]
func (h ResetPasswordHandler) VerifyRecoveryCode(ctx *gin.Context) {
	req := dto.VerifyRecoveryCodeInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.resetPasswordService.VerifyRecoveryCode(ctx, req); err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidOrExpiredOtp):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}

// ResetPassword godoc
//
//	@Id				ResetPassword
//	@Summary		Reset password
//	@Description	Request for reset password. If code is correct will be updated password
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.ResetPasswordInput	true	"Reset password body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/reset-password [post]
func (h ResetPasswordHandler) ResetPassword(ctx *gin.Context) {
	req := dto.ResetPasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.resetPasswordService.ResetPassword(ctx, req); err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidOrExpiredOtp):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, auth.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
}
