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

var (
	ErrInvalidProvider     = dto2.NewI18nError("invalid provider", "errors.invalid_provider")
	ErrRedirectUrlNotFound = dto2.NewI18nError("not found redirect url", "errors.redirect_url_not_found")
	ErrInvalidState        = dto2.NewI18nError("invalid state", "errors.invalid_state")
	ErrInvalidCode         = dto2.NewI18nError("invalid code", "errors.invalid_code")

	stateCookieName = "OAuthState"

	stateCookieMaxAge = 20 * 60
)

type SocialLoginHandler struct {
	socialLoginServices map[string]*auth.SocialLoginService
	cfg                 *config.Config
}

func NewSocialLoginHandler(socialLoginServices map[string]*auth.SocialLoginService, cfg *config.Config) SocialLoginHandler {
	return SocialLoginHandler{
		socialLoginServices: socialLoginServices,
		cfg:                 cfg,
	}
}

// SocialLogin godoc
//
//	@Id				SocialLogin
//	@Summary		Social login
//	@Description	Request for redirecting to OAuth consent page. After authorization, it will redirect to redirectUrl with authorization code and state
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			provider	path	string	true	"Social login provider (yandex, google, mailru)"
//	@Param			redirectUrl	query	string	true	"Redirect URL"
//	@Success		302
//	@Header			302	{string}	Set-Cookie	"OAuthGoogleState=; HttpOnly; Max-Age=1200; Secure"
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/social/{provider} [get]
func (h SocialLoginHandler) SocialLogin(ctx *gin.Context) {
	// Get provider
	provider := ctx.Param("provider")
	svc, ok := h.socialLoginServices[provider]
	if !ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrInvalidProvider)
		return
	}

	// Get redirect url
	redirectUrl, ok := ctx.GetQuery("redirectUrl")
	if !ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrRedirectUrlNotFound)
		return
	}

	// Get consent page url
	output := svc.GetConsentPageUrl(ctx, redirectUrl)

	// Set cookies amd redirect
	ctx.SetCookie(stateCookieName, output.OauthState, stateCookieMaxAge, "", "", false, true)
	ctx.Redirect(http.StatusFound, output.ConsentPageUrl)
}

// SocialLoginCallback godoc
//
//	@Id				SocialLoginCallback
//	@Summary		Social login callback
//	@Description	Request for exchanging authorization code to token pairs. In process, it will exchange code to user info and register new user or login existing user. In response will be new access token in body and new refresh tokens in http-only cookie.
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Param			provider	path		string							true	"Social login provider (yandex, google, mailru)"
//	@Param			body		body		dto.SocialLoginCallbackInput	true	"Social login callback request body"
//	@Success		200			{object}	dto.JwtTokensOutput
//	@Header			200			{string}	Set-Cookie	"RefreshToken=; HttpOnly; Max-Age=86400; Secure"
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		403			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/v0/auth/social/{provider}/callback [post]
func (h SocialLoginHandler) SocialLoginCallback(ctx *gin.Context) {
	// Get provider
	provider := ctx.Param("provider")

	// Check if provider is valid
	svc, ok := h.socialLoginServices[provider]
	if !ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrInvalidProvider)
		return
	}

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
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrInvalidState)
		return
	}
	if cookieState != req.State {
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrInvalidState)
		return
	}

	// Fetch user info
	userInfo, err := svc.FetchUserInfo(ctx, dto.FetchUserInfoInput{Code: req.Code})
	if err != nil {
		_ = ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Register or login
	res, err := svc.RegisterOrLogin(ctx, userInfo)
	if err != nil {
		switch {
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
