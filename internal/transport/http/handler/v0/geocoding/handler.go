package geocoding

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/service/domain"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/internal/transport/http/middleware"
	"github.com/mandarine-io/backend/internal/transport/http/util"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"net/http"
)

type handler struct {
	svc    domain.GeocodingService
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.GeocodingService, opts ...Option) apihandler.APIHandler {
	s := &handler{
		svc:    svc,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	h.logger.Debug().Msg("register geocoding routes")

	router.GET(
		"v0/geocode/forward",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.Geocode,
	)
	router.GET(
		"v0/geocode/reverse",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.ReverseGeocode,
	)
}

// Geocode godoc
//
//	@Id				Geocode
//	@Summary		Geocode
//	@Description	Request for geocoding. User must be logged in. In response will be returned coordinates.
//	@Security		BearerAuth
//	@Tags			Geocoding API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			params	query		v0.GeocodingInput	true	"Geocoding query parameters"
//	@Success		200		{object}	v0.GeocodingOutput	"Geocoded coordinates"
//	@Failure		400		{object}	v0.ErrorOutput		"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput		"Unauthorized error"
//	@Failure		403		{object}	v0.ErrorOutput		"User is blocked or deleted"
//	@Failure		503		{object}	v0.ErrorOutput		"Geocoding service is unavailable"
//	@Router			/v0/geocode/forward [get]
func (h *handler) Geocode(ctx *gin.Context) {
	h.logger.Debug().Msg("handle geocode")

	input := v0.GeocodingInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	lang := language.English
	langRaw, ok := ctx.Get(middleware.LangKey)
	if ok {
		parsedLang, err := language.Parse(langRaw.(string))
		if err != nil {
			h.logger.Warn().Err(err).Msg("failed to parse language")
		} else {
			lang = parsedLang
		}
	}

	res, err := h.svc.Geocode(ctx, input, lang)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrGeocodeProvidersUnavailable):
			_ = util.ErrorWithStatus(ctx, http.StatusServiceUnavailable, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// ReverseGeocode godoc
//
//	@Id				ReverseGeocode
//	@Summary		Reverse geocode
//	@Description	Request for reverse geocoding. User must be logged in. In response will be returned address.
//	@Security		BearerAuth
//	@Tags			Geocoding API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			params	query		v0.ReverseGeocodingInput		true	"Reverse geocoding query parameters"
//	@Success		200		{object}	v0.ReverseGeocodingOutput	"Reverse geocoded address"
//	@Failure		400		{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403		{object}	v0.ErrorOutput				"User is blocked or deleted"
//	@Failure		503		{object}	v0.ErrorOutput				"Geocoding service is unavailable"
//	@Router			/v0/geocode/reverse [get]
func (h *handler) ReverseGeocode(ctx *gin.Context) {
	h.logger.Debug().Msg("handle reverse geocode")

	input := v0.ReverseGeocodingInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	lang := language.English
	langRaw, ok := ctx.Get(middleware.LangKey)
	if ok {
		parsedLang, err := language.Parse(langRaw.(string))
		if err != nil {
			h.logger.Warn().Err(err).Msg("failed to parse language")
		} else {
			lang = parsedLang
		}
	}

	res, err := h.svc.ReverseGeocode(ctx, input, lang)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusServiceUnavailable, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
