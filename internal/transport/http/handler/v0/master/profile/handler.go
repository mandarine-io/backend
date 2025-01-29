package profile

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/service/domain"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/internal/transport/http/middleware"
	"github.com/mandarine-io/backend/internal/transport/http/util"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type handler struct {
	svc    domain.MasterProfileService
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.MasterProfileService, opts ...Option) apihandler.APIHandler {
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
	log.Debug().Msg("register master profile routes")

	router.POST(
		"v0/masters/profiles",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.CreateMasterProfile,
	)
	router.PATCH(
		"v0/masters/profiles",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.UpdateMasterProfile,
	)
	router.GET(
		"v0/masters/profiles",
		//middleware.Registry.Auth,
		//middleware.Registry.BannedUser,
		//middleware.Registry.DeletedUser,
		h.FindMasterProfiles,
	)
	router.GET(
		"v0/masters/profiles/:username",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.GetMasterProfile,
	)
}

// CreateMasterProfile godoc
//
//	@Id				CreateMasterProfile
//	@Summary		Create master profile
//	@Description	Request for creating master profile. User must be logged in. In response will be returned created master profile.
//	@Security		BearerAuth
//	@Tags			Master Profile API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body		v0.CreateMasterProfileInput	true	"Create master profile request body"
//	@Success		201		{object}	v0.MasterProfileOutput		"Created master profile"
//	@Failure		400		{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403		{object}	v0.ErrorOutput				"User is blocked or deleted"
//	@Failure		404		{object}	v0.ErrorOutput				"Master profile not found"
//	@Failure		409		{object}	v0.ErrorOutput				"Master profile already exists"
//	@Failure		500		{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles [post]
func (h *handler) CreateMasterProfile(ctx *gin.Context) {
	log.Debug().Msg("handle create master profile")

	input := v0.CreateMasterProfileInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	resp, err := h.svc.CreateMasterProfile(ctx, principal.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateMasterProfile):
			_ = util.ErrorWithStatus(ctx, http.StatusConflict, err)
		case errors.Is(err, domain.ErrUserNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateMasterProfile godoc
//
//	@Id				UpdateMasterProfile
//	@Summary		Update master profile
//	@Description	Request for updating master profile. User must be logged in. In response will be returned updated master profile.
//	@Security		BearerAuth
//	@Tags			Master Profile API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	body		v0.UpdateMasterProfileInput	true	"Update master profile request body"
//	@Success		200		{object}	v0.MasterProfileOutput		"Updated master profile"
//	@Failure		400		{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403		{object}	v0.ErrorOutput				"User is blocked or deleted"
//	@Failure		404		{object}	v0.ErrorOutput				"Master profile not found"
//	@Failure		500		{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles [patch]
func (h *handler) UpdateMasterProfile(ctx *gin.Context) {
	log.Debug().Msg("handle update master profile")

	input := v0.UpdateMasterProfileInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	resp, err := h.svc.UpdateMasterProfile(ctx, principal.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMasterProfileNotExist):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// FindMasterProfiles godoc
//
//	@Id				FindMasterProfiles
//	@Summary		Find master profiles
//	@Description	Request for finding master profiles. User must be logged in. In response will be returned found master profiles.
//	@Security		BearerAuth
//	@Tags			Master Profile API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	query		v0.FindMasterProfilesInput	true	"Query parameters for finding master profiles"
//	@Success		200		{object}	v0.MasterProfilesOutput		"Found master profiles"
//	@Failure		400		{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403		{object}	v0.ErrorOutput				"User is blocked or deleted"
//	@Failure		404		{object}	v0.ErrorOutput				"Master profile not found"
//	@Failure		500		{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles [get]
func (h *handler) FindMasterProfiles(ctx *gin.Context) {
	log.Debug().Msg("handle find master profiles")

	input := v0.FindMasterProfilesInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	resp, err := h.svc.FindMasterProfiles(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnavailableSortField):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetMasterProfile godoc
//
//	@Id				GetMasterProfile
//	@Summary		Get master profile
//	@Description	Request for getting master profile. User must be logged in. In response will be returned found master profile.
//	@Security		BearerAuth
//	@Tags			Master Profile API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			username	path		string						true	"Username"
//	@Success		200			{object}	v0.MasterProfileOutput	"Found master profile"
//	@Failure		401			{object}	v0.ErrorOutput			"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput			"User is blocked or deleted; Own master profile is disabled"
//	@Failure		404			{object}	v0.ErrorOutput			"Master profile not found"
//	@Failure		500			{object}	v0.ErrorOutput			"Internal server error"
//	@Router			/v0/masters/profiles/{username} [get]
func (h *handler) GetMasterProfile(ctx *gin.Context) {
	log.Debug().Msg("handle get master profile")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")

	if principal.Username == username {
		resp, err := h.svc.GetOwnMasterProfile(ctx, principal.ID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrMasterProfileDisabled):
				_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
			case errors.Is(err, domain.ErrMasterProfileNotExist):
				_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
			default:
				_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp, err := h.svc.GetMasterProfileByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMasterProfileNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
