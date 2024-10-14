package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mandarine/internal/api/config"
	"mandarine/internal/api/service/auth"
	"mandarine/internal/api/service/auth/dto"
	dto2 "mandarine/pkg/rest/dto"
	"net/http"
)

const (
	refreshTokenCookieName = "RefreshToken"
)

var (
	errSessionExpired = dto2.NewI18nError("session expired", "errors.session_expired")
)

type LoginHandler struct {
	loginService *auth.LoginService
	cfg          *config.Config
}

func NewLoginHandler(loginService *auth.LoginService, cfg *config.Config) LoginHandler {
	return LoginHandler{
		loginService: loginService,
		cfg:          cfg,
	}
}

// Login godoc
//
//	@Id				Login
//	@Summary		Sign in
//	@Description	Request for authentication. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginInput	true	"Login request body"
//	@Success		200		{object}	dto.JwtTokensOutput
//	@Header			200		{string}	Set-Cookie	"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v0/auth/login [post]
func (h LoginHandler) Login(ctx *gin.Context) {
	req := dto.LoginInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := h.loginService.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, auth.ErrBadCredentials):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, auth.ErrUserIsBlocked):
			_ = ctx.AbortWithError(http.StatusForbidden, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.SetCookie(refreshTokenCookieName, res.RefreshToken, h.cfg.Security.JWT.RefreshTokenTTL, "", "", true, true)
	ctx.JSON(http.StatusOK, res)
}

// RefreshTokens godoc
//
//	@Id				RefreshTokens
//	@Summary		Refresh tokens
//	@Description	Request for refreshing tokens. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.JwtTokensOutput
//	@Header			200	{string}	Set-Cookie	"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/refresh [get]
func (h LoginHandler) RefreshTokens(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(refreshTokenCookieName)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, errSessionExpired)
		return
	}

	res, err := h.loginService.RefreshTokens(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidJwtToken):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, auth.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, auth.ErrUserIsBlocked):
			_ = ctx.AbortWithError(http.StatusForbidden, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.SetCookie(refreshTokenCookieName, res.RefreshToken, h.cfg.Security.JWT.RefreshTokenTTL, "", "", true, true)
	ctx.JSON(http.StatusOK, res)
}
