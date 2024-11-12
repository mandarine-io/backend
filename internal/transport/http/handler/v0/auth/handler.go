package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	apihandler "github.com/mandarine-io/Backend/internal/transport/http/handler"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	middleware2 "github.com/mandarine-io/Backend/pkg/transport/http/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	refreshTokenCookieName = "RefreshToken"
)

var (
	errSessionExpired      = dto2.NewI18nError("session expired", "errors.session_expired")
	errRedirectUrlNotFound = dto2.NewI18nError("not found redirect url", "errors.redirect_url_not_found")
	errInvalidState        = dto2.NewI18nError("invalid state", "errors.invalid_state")

	stateCookieName   = "OAuthState"
	stateCookieMaxAge = 20 * 60
)

type handler struct {
	svc service.AuthService
	cfg *config.Config
}

func NewHandler(svc service.AuthService, cfg *config.Config) apihandler.ApiHandler {
	return &handler{
		svc: svc,
		cfg: cfg,
	}
}

func (h *handler) RegisterRoutes(router *gin.Engine, middlewares apihandler.RouteMiddlewares) {
	log.Debug().Msg("register service routes")

	router.POST("v0/auth/login", h.Login)
	router.GET("v0/auth/refresh", h.RefreshTokens)
	router.GET("v0/auth/social/:provider", h.SocialLogin)
	router.POST("v0/auth/social/:provider/callback", h.SocialLoginCallback)
	router.POST("v0/auth/register", h.Register)
	router.POST("v0/auth/register/confirm", h.RegisterConfirm)
	router.POST("v0/auth/recovery-password", h.RecoveryPassword)
	router.POST("v0/auth/recovery-password/verify", h.VerifyRecoveryCode)
	router.POST("v0/auth/reset-password", h.ResetPassword)

	router.GET("v0/auth/logout", middlewares.Auth, h.Logout)
}

// Login godoc
//
//	@Id				Login
//	@Summary		Sign in
//	@Description	Request for serviceentication. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginInput	true	"Login request body"
//	@Header			200		{string}	Set-Cookie	"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Success		200		{object}	dto.JwtTokensOutput	"JWT access token"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked"
//	@Failure		404		{object}	dto.ErrorResponse	"User not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/login [post]
func (h *handler) Login(ctx *gin.Context) {
	log.Debug().Msg("handle login")

	req := dto.LoginInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrBadCredentials):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, service.ErrUserIsBlocked):
			_ = ctx.AbortWithError(http.StatusForbidden, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.SetCookie(refreshTokenCookieName, res.RefreshToken, h.cfg.Security.JWT.RefreshTokenTTL, "/v0/auth/refresh", "", true, true)
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
//	@Header			200	{string}	Set-Cookie	"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Success		200	{object}	dto.JwtTokensOutput	"JWT access token"
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		403	{object}	dto.ErrorResponse	"User is blocked"
//	@Failure		404	{object}	dto.ErrorResponse	"User not found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/refresh [get]
func (h *handler) RefreshTokens(ctx *gin.Context) {
	log.Debug().Msg("handle refresh tokens")

	refreshToken, err := ctx.Cookie(refreshTokenCookieName)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, errSessionExpired)
		return
	}

	res, err := h.svc.RefreshTokens(ctx, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidJwtToken):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrUserIsBlocked):
			_ = ctx.AbortWithError(http.StatusForbidden, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.SetCookie(refreshTokenCookieName, res.RefreshToken, h.cfg.Security.JWT.RefreshTokenTTL, "", "", true, true)
	ctx.JSON(http.StatusOK, res)
}

// Logout godoc
//
//	@Id				Logout
//	@Summary		Logout
//	@Description	Request for logout. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/logout [get]
func (h *handler) Logout(c *gin.Context) {
	log.Debug().Msg("handle logout")

	principal, err := middleware2.GetAuthUser(c)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	err = h.svc.Logout(c, principal.JTI)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// Register godoc
//
//	@Id				RegisterClient
//	@Summary		RegisterClient
//	@Description	Request for creating new user. At the end will be sent confirmation email with code
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.RegisterInput	true	"RegisterClient request body"
//	@Success		202
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		409	{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/register [post]
func (h *handler) Register(ctx *gin.Context) {
	log.Debug().Msg("handle register")

	req := dto.RegisterInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	log.Debug().Msg("get localizer")
	localizer := ctx.Value(middleware2.LocalizerKey).(*i18n.Localizer)

	if err := h.svc.Register(ctx, req, localizer); err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateUser):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		case errors.Is(err, service.ErrSendEmail):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
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
//	@Summary		RegisterClient confirmation
//	@Description	Request for confirming registration. At the end will be created new user
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.RegisterConfirmInput	true	"RegisterClient confirm body"
//	@Success		200
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		409	{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/register/confirm [post]
func (h *handler) RegisterConfirm(ctx *gin.Context) {
	log.Debug().Msg("handle register confirm")

	req := dto.RegisterConfirmInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.svc.RegisterConfirm(ctx, req); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrExpiredOtp):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		case errors.Is(err, service.ErrDuplicateUser):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusOK)
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
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		404	{object}	dto.ErrorResponse	"User not found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/recovery-password [post]
func (h *handler) RecoveryPassword(ctx *gin.Context) {
	log.Debug().Msg("handle recovery password")

	req := dto.RecoveryPasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	log.Debug().Msg("get localizer")
	localizer := ctx.Value(middleware2.LocalizerKey).(*i18n.Localizer)

	if err := h.svc.RecoveryPassword(ctx, req, localizer); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		case errors.Is(err, service.ErrSendEmail):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
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
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/recovery-password/verify [post]
func (h *handler) VerifyRecoveryCode(ctx *gin.Context) {
	log.Debug().Msg("handle verify recovery code")

	req := dto.VerifyRecoveryCodeInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.svc.VerifyRecoveryCode(ctx, req); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrExpiredOtp):
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
//	@Failure		400	{object}	dto.ErrorResponse	"Validation error"
//	@Failure		404	{object}	dto.ErrorResponse	"User not found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/reset-password [post]
func (h *handler) ResetPassword(ctx *gin.Context) {
	log.Debug().Msg("handle reset password")

	req := dto.ResetPasswordInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ResetPassword(ctx, req); err != nil {
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

// SocialLogin godoc
//
//	@Id				SocialLogin
//	@Summary		Social login
//	@Description	Request for redirecting to OAuth consent page. After serviceorization, it will redirect to redirectUrl with serviceorization code and state
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			provider	path	string	true	"Social login provider (yandex, google, mailru)"
//	@Param			redirectUrl	query	string	true	"Redirect URL"
//	@Header			302	{string}	Set-Cookie	"OAuthGoogleState=; HttpOnly; Max-Age=1200; Secure"
//	@Success		302
//	@Failure		404	{object}	dto.ErrorResponse	"Provider not found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/social/{provider} [get]
func (h *handler) SocialLogin(ctx *gin.Context) {
	log.Debug().Msg("handle social login")

	// Get provider
	provider := ctx.Param("provider")

	// Get redirect url
	redirectUrl, ok := ctx.GetQuery("redirectUrl")
	if !ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, errRedirectUrlNotFound)
		return
	}

	// Get consent page url
	output, err := h.svc.GetConsentPageUrl(ctx, provider, redirectUrl)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProvider):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	// Set cookies amd redirect
	ctx.SetCookie(stateCookieName, output.OauthState, stateCookieMaxAge, "", "", false, true)
	ctx.Redirect(http.StatusFound, output.ConsentPageUrl)
}

// SocialLoginCallback godoc
//
//	@Id				SocialLoginCallback
//	@Summary		Social login callback
//	@Description	Request for exchanging serviceorization code to token pairs. In process, it will exchange code to user info and register new user or login existing user. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string							true	"Social login provider (yandex, google, mailru)"
//	@Param			body		body		dto.SocialLoginCallbackInput	true	"Social login callback request body"
//	@Header			200			{string}	Set-Cookie	"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Success		200			{object}	dto.JwtTokensOutput	"JWT access token"
//	@Failure		400			{object}	dto.ErrorResponse	"Validation error"
//	@Failure		403			{object}	dto.ErrorResponse	"User already exists"
//	@Failure		404			{object}	dto.ErrorResponse	"User not found"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/auth/social/{provider}/callback [post]
func (h *handler) SocialLoginCallback(ctx *gin.Context) {
	log.Debug().Msg("handle social login callback")

	// Get provider
	provider := ctx.Param("provider")

	// Bind request
	req := dto.SocialLoginCallbackInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Get and check state
	cookieState, err := ctx.Cookie(stateCookieName)
	ctx.SetCookie(stateCookieName, "", -1, "", "", true, true)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, errInvalidState)
		return
	}
	if cookieState != req.State {
		_ = ctx.AbortWithError(http.StatusBadRequest, errInvalidState)
		return
	}

	// Fetch user info
	userInfo, err := h.svc.FetchUserInfo(ctx, provider, dto.FetchUserInfoInput{Code: req.Code})
	if err != nil {
		_ = ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	// RegisterClient or login
	res, err := h.svc.RegisterOrLogin(ctx, userInfo)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserIsBlocked):
			_ = ctx.AbortWithError(http.StatusForbidden, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.SetCookie(refreshTokenCookieName, res.RefreshToken, h.cfg.Security.JWT.RefreshTokenTTL, "", "", true, true)
	ctx.JSON(http.StatusOK, res)
}
