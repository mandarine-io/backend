package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/config"
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

var (
	errRedirectURLNotFound = v0.NewI18nError("not found redirect url", "errors.redirect_url_not_found")
	errInvalidState        = v0.NewI18nError("invalid state", "errors.invalid_state")

	stateCookieName   = "OAuthState"
	stateCookieMaxAge = 20 * 60
)

type handler struct {
	svc    domain.AuthService
	cfg    config.Config
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.AuthService, cfg config.Config, opts ...Option) apihandler.APIHandler {
	h := &handler{
		svc:    svc,
		cfg:    cfg,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	h.logger.Debug().Msg("register service routes")

	authRouter := router.Group("/v0/auth")
	{
		authRouter.POST("/login", h.Login)
		authRouter.POST("/refresh", h.RefreshTokens)
		authRouter.GET("/social/:provider", h.SocialLogin)
		authRouter.POST("/social/:provider/callback", h.SocialLoginCallback)
		authRouter.POST("/register", h.Register)
		authRouter.POST("/register/confirm", h.RegisterConfirm)
		authRouter.POST("/recovery-password", h.RecoveryPassword)
		authRouter.POST("/recovery-password/verify", h.VerifyRecoveryCode)
		authRouter.POST("/reset-password", h.ResetPassword)

		authRouter.GET(
			"/logout",
			middleware.Registry.Auth,
			h.Logout,
		)
	}
}

// Login godoc
//
//	@Id				Login
//	@Summary		Sign in
//	@Description	Request for serviceentication. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body		v0.LoginInput		true	"Login request body"
//	@Header			200		{string}	Set-Cookie				"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Success		200		{object}	v0.JwtTokensOutput	"JWT tokens"
//	@Failure		400		{object}	v0.ErrorOutput		"Validation error"
//	@Failure		403		{object}	v0.ErrorOutput		"User is blocked"
//	@Failure		404		{object}	v0.ErrorOutput		"User not found"
//	@Failure		500		{object}	v0.ErrorOutput		"Internal server error"
//	@Router			/v0/auth/login [post]
func (h *handler) Login(ctx *gin.Context) {
	h.logger.Debug().Msg("handle login")

	input := v0.LoginInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.Login(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrBadCredentials):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domain.ErrUserIsBlocked):
			_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// RefreshTokens godoc
//
//	@Id				RefreshTokens
//	@Summary		Refresh tokens
//	@Description	Request for refreshing tokens. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body		v0.RefreshTokensInput	true	"Refresh token body"
//	@Success		200		{object}	v0.JwtTokensOutput		"JWT tokens"
//	@Failure		400		{object}	v0.ErrorOutput			"Validation error"
//	@Failure		403		{object}	v0.ErrorOutput			"User is blocked"
//	@Failure		404		{object}	v0.ErrorOutput			"User not found"
//	@Failure		500		{object}	v0.ErrorOutput			"Internal server error"
//	@Router			/v0/auth/refresh [post]
func (h *handler) RefreshTokens(ctx *gin.Context) {
	h.logger.Debug().Msg("handle refresh tokens")

	input := v0.RefreshTokensInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.svc.RefreshTokens(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrInvalidJWTToken):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrUserIsBlocked):
			_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// Logout godoc
//
//	@Id				Logout
//	@Summary		Logout
//	@Description	Request for logout. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Authentication and Authorization API
//	@Accept			application/json
//	@Produce		application/json
//	@Success		200
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/logout [get]
func (h *handler) Logout(c *gin.Context) {
	h.logger.Debug().Msg("handle logout")

	principal, err := middleware.GetAuthUser(c)
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
//	@Id				Register
//	@Summary		Register
//	@Description	Request for creating new user. At the end will be sent confirmation email with code
//	@Tags			Authentication and Authorization API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.RegisterInput	true	"Register request body"
//	@Success		202
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		409	{object}	v0.ErrorOutput	"User already exists"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/register [post]
func (h *handler) Register(ctx *gin.Context) {
	h.logger.Debug().Msg("handle register")

	input := v0.RegisterInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	h.logger.Debug().Msg("get localizer")
	localizer := ctx.Value(middleware.LocalizerKey).(locale.Localizer)

	if err := h.svc.Register(ctx, input, localizer); err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateUser):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		case errors.Is(err, domain.ErrSendEmail):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.RegisterConfirmInput	true	"Register confirm body"
//	@Success		200
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		409	{object}	v0.ErrorOutput	"User already exists"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/register/confirm [post]
func (h *handler) RegisterConfirm(ctx *gin.Context) {
	h.logger.Debug().Msg("handle register confirm")

	input := v0.RegisterConfirmInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.RegisterConfirm(ctx, input); err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrInvalidOrExpiredOTP):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		case errors.Is(err, domain.ErrDuplicateUser):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.RecoveryPasswordInput	true	"Recovery password body"
//	@Success		202
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		404	{object}	v0.ErrorOutput	"User not found"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/recovery-password [post]
func (h *handler) RecoveryPassword(ctx *gin.Context) {
	h.logger.Debug().Msg("handle recovery password")

	input := v0.RecoveryPasswordInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	h.logger.Debug().Msg("get localizer")
	localizer := ctx.Value(middleware.LocalizerKey).(locale.Localizer)

	if err := h.svc.RecoveryPassword(ctx, input, localizer); err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrSendEmail):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.VerifyRecoveryCodeInput	true	"Verify recovery code body"
//	@Success		200
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/recovery-password/verify [post]
func (h *handler) VerifyRecoveryCode(ctx *gin.Context) {
	h.logger.Debug().Msg("handle verify recovery code")

	input := v0.VerifyRecoveryCodeInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.VerifyRecoveryCode(ctx, input); err != nil {
		switch {
		case errors.Is(err, infrastructure.ErrInvalidOrExpiredOTP):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
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
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body	v0.ResetPasswordInput	true	"Reset password body"
//	@Success		200
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		404	{object}	v0.ErrorOutput	"User not found"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/reset-password [post]
func (h *handler) ResetPassword(ctx *gin.Context) {
	h.logger.Debug().Msg("handle reset password")

	input := v0.ResetPasswordInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.ResetPassword(ctx, input); err != nil {
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

// SocialLogin godoc
//
//	@Id				SocialLogin
//	@Summary		Social login
//	@Description	Request for redirecting to OAuth consent page. After serviceorization, it will redirect to redirectURL with serviceorization code and state
//	@Tags			Authentication and Authorization API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			provider	path		string		true	"Social login provider (yandex, google, mailru)"
//	@Param			redirectURL	query		string		true	"Redirect URL"
//	@Header			302			{string}	Set-Cookie	"OAuthGoogleState=; HttpOnly; Max-Age=1200; Secure"
//	@Success		302
//	@Failure		404	{object}	v0.ErrorOutput	"Provider not found"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/auth/social/{provider} [get]
func (h *handler) SocialLogin(ctx *gin.Context) {
	h.logger.Debug().Msg("handle social login")

	// Get provider
	p := ctx.Param("provider")

	// Get redirect url
	redirectURL, ok := ctx.GetQuery("redirectURL")
	if !ok {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, errRedirectURLNotFound)
		return
	}

	// Get consent page url
	output, err := h.svc.GetConsentPageURL(ctx, p, redirectURL)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidProvider):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Set cookies amd redirect
	ctx.SetCookie(stateCookieName, output.OauthState, stateCookieMaxAge, "", "", false, true)
	ctx.Redirect(http.StatusFound, output.ConsentPageURL)
}

// SocialLoginCallback godoc
//
//	@Id				SocialLoginCallback
//	@Summary		Social login callback
//	@Description	Request for exchanging serviceorization code to token pairs. In process, it will exchange code to user info and register new user or login existing user. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			provider	path		string							true	"Social login provider (yandex, google, mailru)"
//	@Param			input		body		v0.SocialLoginCallbackInput	true	"Social login callback request body"
//	@Success		200			{object}	v0.JwtTokensOutput			"JWT tokens"
//	@Failure		400			{object}	v0.ErrorOutput				"Validation error"
//	@Failure		403			{object}	v0.ErrorOutput				"User already exists"
//	@Failure		404			{object}	v0.ErrorOutput				"User not found"
//	@Failure		500			{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/auth/social/{provider}/callback [post]
func (h *handler) SocialLoginCallback(ctx *gin.Context) {
	h.logger.Debug().Msg("handle social login callback")

	// Get provider
	p := ctx.Param("provider")

	// Bind request
	input := v0.SocialLoginCallbackInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	// Get and check state
	cookieState, err := ctx.Cookie(stateCookieName)
	ctx.SetCookie(stateCookieName, "", -1, "", "", true, true)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, errInvalidState)
		return
	}
	if cookieState != input.State {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, errInvalidState)
		return
	}

	// Fetch user info
	userInfo, err := h.svc.FetchUserInfo(ctx, p, v0.FetchUserInfoInput{Code: input.Code})
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		return
	}

	// Register or login
	res, err := h.svc.RegisterOrLogin(ctx, userInfo)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserIsBlocked):
			_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}
