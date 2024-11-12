package geocoding

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	handler2 "github.com/mandarine-io/Backend/internal/transport/http/handler"
	"github.com/mandarine-io/Backend/pkg/transport/http/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"net/http"
)

type handler struct {
	svc service.GeocodingService
}

func NewHandler(svc service.GeocodingService) handler2.ApiHandler {
	return &handler{svc: svc}
}

func (h *handler) RegisterRoutes(router *gin.Engine, middlewares handler2.RouteMiddlewares) {
	log.Debug().Msg("register geocoding routes")

	router.GET(
		"v0/geocode/forward",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.Geocode,
	)
	router.GET(
		"v0/geocode/reverse",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
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
//	@Accept			json
//	@Produce		json
//	@Param			params	query		dto.GeocodingInput	true	"Geocoding query parameters"
//	@Success		200		{object}	dto.GeocodingOutput	"Geocoded coordinates"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		503		{object}	dto.ErrorResponse	"Geocoding service is unavailable"
//	@Router			/v0/geocode/forward [get]
func (h *handler) Geocode(ctx *gin.Context) {
	log.Debug().Msg("handle geocode")

	input := dto.GeocodingInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	lang := language.English
	langRaw, ok := ctx.Get(middleware.LangKey)
	if ok {
		parsedLang, err := language.Parse(langRaw.(string))
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse language")
		} else {
			lang = parsedLang
		}
	}

	res, err := h.svc.Geocode(ctx, input, lang)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGeocodeProvidersUnavailable):
			_ = ctx.AbortWithError(http.StatusServiceUnavailable, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
//	@Accept			json
//	@Produce		json
//	@Param			params	query		dto.ReverseGeocodingInput	true	"Reverse geocoding query parameters"
//	@Success		200		{object}	dto.ReverseGeocodingOutput	"Reverse geocoded address"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		503		{object}	dto.ErrorResponse	"Geocoding service is unavailable"
//	@Router			/v0/geocode/reverse [get]
func (h *handler) ReverseGeocode(ctx *gin.Context) {
	log.Debug().Msg("handle reverse geocode")

	input := dto.ReverseGeocodingInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	lang := language.English
	langRaw, ok := ctx.Get(middleware.LangKey)
	if ok {
		parsedLang, err := language.Parse(langRaw.(string))
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse language")
		} else {
			lang = parsedLang
		}
	}

	res, err := h.svc.ReverseGeocode(ctx, input, lang)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusServiceUnavailable, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
