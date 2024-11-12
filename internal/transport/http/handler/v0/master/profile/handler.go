package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	apihandler "github.com/mandarine-io/Backend/internal/transport/http/handler"
	"github.com/mandarine-io/Backend/pkg/transport/http/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	svc service.MasterProfileService
}

func NewHandler(svc service.MasterProfileService) apihandler.ApiHandler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, middlewares apihandler.RouteMiddlewares) {
	log.Debug().Msg("register master profile routes")

	router.POST(
		"v0/masters/profile",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.CreateMasterProfile,
	)
	router.PATCH(
		"v0/masters/profile",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.UpdateMasterProfile,
	)
	router.GET(
		"v0/masters/profile",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.FindMasterProfiles,
	)
	router.GET(
		"v0/masters/profile/:username",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
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
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.CreateMasterProfileInput	true	"Create master profile request body"
//	@Success		201		{object}	dto.OwnMasterProfileOutput "Created master profile"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404		{object}	dto.ErrorResponse	"Master profile not found"
//	@Failure		409		{object}	dto.ErrorResponse	"Master profile already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/masters/profile [post]
func (h *Handler) CreateMasterProfile(ctx *gin.Context) {
	log.Debug().Msg("handle create master profile")

	req := dto.CreateMasterProfileInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	resp, err := h.svc.CreateMasterProfile(ctx, principal.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateMasterProfile):
			_ = ctx.AbortWithError(http.StatusConflict, err)
		case errors.Is(err, service.ErrUserNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateMasterProfileInput	true	"Update master profile request body"
//	@Success		200		{object}	dto.OwnMasterProfileOutput	"Updated master profile"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404		{object}	dto.ErrorResponse	"Master profile not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/masters/profile [patch]
func (h *Handler) UpdateMasterProfile(ctx *gin.Context) {
	log.Debug().Msg("handle update master profile")

	req := dto.UpdateMasterProfileInput{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	resp, err := h.svc.UpdateMasterProfile(ctx, principal.ID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMasterProfileNotExist):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
//	@Accept			json
//	@Produce		json
//	@Param			input	query		dto.FindMasterProfilesInput	true	"Query parameters for finding master profiles"
//	@Success		200		{object}	dto.MasterProfilesOutput	"Found master profiles"
//	@Failure		400		{object}	dto.ErrorResponse	"Validation error"
//	@Failure		401		{object}	dto.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404		{object}	dto.ErrorResponse	"Master profile not found"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/masters/profile [get]
func (h *Handler) FindMasterProfiles(ctx *gin.Context) {
	log.Debug().Msg("handle find master profiles")

	input := dto.FindMasterProfilesInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := h.svc.FindMasterProfiles(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnavailableSortField):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string	true	"Username"
//	@Success		200			{object}	dto.MasterProfileOutput	"Found master profile"
//	@Success		200			{object}	dto.OwnMasterProfileOutput	"Found own master profile"
//	@Failure		401			{object}	dto.ErrorResponse	"Unauthorized error"
//	@Failure		403			{object}	dto.ErrorResponse	"User is blocked or deleted"
//	@Failure		404			{object}	dto.ErrorResponse	"Master profile not found"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal server error"
//	@Router			/v0/masters/profile/{username} [get]
func (h *Handler) GetMasterProfile(ctx *gin.Context) {
	log.Debug().Msg("handle get master profile")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")

	if principal.Username != username {
		resp, err := h.svc.GetMasterProfileByUsername(ctx, username)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrMasterProfileNotFound):
				_ = ctx.AbortWithError(http.StatusNotFound, err)
			default:
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp, err := h.svc.GetOwnMasterProfile(ctx, principal.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMasterProfileNotExist):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
