package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	svc    domain.MasterServiceService
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.MasterServiceService, opts ...Option) apihandler.APIHandler {
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
	log.Debug().Msg("register master service routes")

	router.POST(
		"v0/masters/profiles/:username/services",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.CreateMasterService,
	)
	router.PATCH(
		"v0/masters/profiles/:username/services/:id",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.UpdateMasterService,
	)
	router.DELETE(
		"v0/masters/profiles/:username/services/:id",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.DeleteMasterService,
	)
	router.GET(
		"v0/masters/profiles/-/services",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.FindMasterServices,
	)
	router.GET(
		"v0/masters/profiles/:username/services",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.FindMasterServicesByUsername,
	)
	router.GET(
		"v0/masters/profiles/:username/services/:id",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.GetMasterServiceByUsername,
	)
}

// CreateMasterService godoc
//
//	@Id				CreateMasterService
//	@Summary		Create master service
//	@Description	Request for creating master service. User must be logged in. In response will be returned created master service.
//	@Security		BearerAuth
//	@Tags			Master Service API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			username	path		string							true	"Username"
//	@Param			input		body		v0.CreateMasterServiceInput	true	"Create master service request body"
//	@Success		201			{object}	v0.MasterServiceOutput		"Created master service"
//	@Failure		400			{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401			{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput				"User is blocked or deleted; Master profile is disabled; Cannot create master service for another master profile"
//	@Failure		404			{object}	v0.ErrorOutput				"Master service not found"
//	@Failure		500			{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles/{username}/services [post]
func (h *handler) CreateMasterService(ctx *gin.Context) {
	log.Debug().Msg("handle create master service")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")
	if principal.Username != username {
		_ = util.ErrorWithStatus(ctx, http.StatusForbidden, domain.ErrMasterServiceCreation)
		return
	}

	input := v0.CreateMasterServiceInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	resp, err := h.svc.CreateMasterService(ctx, principal.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMasterProfileNotExist):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrMasterProfileDisabled):
			_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UpdateMasterService godoc
//
//	@Id				UpdateMasterService
//	@Summary		Update master service
//	@Description	Request for updating master service. User must be logged in. In response will be returned updated master service.
//	@Security		BearerAuth
//	@Tags			Master Service API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			username	path		string							true	"Username"
//	@Param			id			path		string							true	"Master service ID"
//	@Param			input		body		v0.UpdateMasterServiceInput	true	"Update master service request body"
//	@Success		200			{object}	v0.MasterServiceOutput		"Updated master service"
//	@Failure		400			{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401			{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput				"User is blocked or deleted; Master profile is disabled; Cannot update master service for another master profile"
//	@Failure		404			{object}	v0.ErrorOutput				"Master service not found; Master profile not found"
//	@Failure		500			{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles/{username}/services/{id} [patch]
func (h *handler) UpdateMasterService(ctx *gin.Context) {
	log.Debug().Msg("handle update master service")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")
	if principal.Username != username {
		_ = util.ErrorWithStatus(ctx, http.StatusForbidden, domain.ErrMasterServiceModification)
		return
	}

	input := v0.UpdateMasterServiceInput{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	masterServiceIDRaw := ctx.Param("id")
	masterServiceID, err := uuid.Parse(masterServiceIDRaw)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	resp, err := h.svc.UpdateMasterService(ctx, principal.ID, masterServiceID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMasterProfileNotExist),
			errors.Is(err, domain.ErrMasterServiceNotExist):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrMasterProfileDisabled):
			_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteMasterService godoc
//
//	@Id				DeleteMasterService
//	@Summary		Delete master service
//	@Description	Request for deleting master service. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Master Service API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			username	path	string	true	"Username"
//	@Param			id			path	string	true	"Master service ID"
//	@Success		204
//	@Failure		400	{object}	v0.ErrorOutput	"Validation error"
//	@Failure		401	{object}	v0.ErrorOutput	"Unauthorized error"
//	@Failure		403	{object}	v0.ErrorOutput	"User is blocked or deleted; Master profile is disabled; Cannot delete master service for another master profile"
//	@Failure		404	{object}	v0.ErrorOutput	"Master service not found; Master profile not found"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/masters/profiles/{username}/services/{id} [delete]
func (h *handler) DeleteMasterService(ctx *gin.Context) {
	log.Debug().Msg("handle delete master service")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")
	if principal.Username != username {
		_ = util.ErrorWithStatus(ctx, http.StatusForbidden, domain.ErrMasterServiceDeletion)
		return
	}

	masterServiceIDRaw := ctx.Param("id")
	masterServiceID, err := uuid.Parse(masterServiceIDRaw)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.svc.DeleteMasterService(ctx, principal.ID, masterServiceID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMasterProfileNotExist):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrMasterProfileDisabled):
			_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
		case errors.Is(err, domain.ErrMasterServiceNotExist):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

// FindMasterServices godoc
//
//	@Id				FindMasterServices
//	@Summary		Find master services
//	@Description	Request for finding master services. User must be logged in. In response will be returned found master services.
//	@Security		BearerAuth
//	@Tags			Master Service API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			input	query		v0.FindMasterServicesInput	true	"Query parameters for finding master services"
//	@Success		200		{object}	v0.MasterServicesOutput		"Found master services"
//	@Failure		400		{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401		{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403		{object}	v0.ErrorOutput				"User is blocked or deleted"
//	@Failure		500		{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles/-/services [get]
func (h *handler) FindMasterServices(ctx *gin.Context) {
	log.Debug().Msg("handle find master services")

	input := v0.FindMasterServicesInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	resp, err := h.svc.FindAllMasterServices(ctx, input)
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

// FindMasterServicesByUsername godoc
//
//	@Id				FindMasterServicesByUsername
//	@Summary		Find master services by username
//	@Description	Request for finding master services by username. User must be logged in. In response will be returned found master services.
//	@Security		BearerAuth
//	@Tags			Master Service API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			username	path		string							true	"Username"
//	@Param			input		query		v0.FindMasterServicesInput	true	"Query parameters for finding master services"
//	@Success		200			{object}	v0.MasterServicesOutput		"Found master services"
//	@Failure		400			{object}	v0.ErrorOutput				"Validation error"
//	@Failure		401			{object}	v0.ErrorOutput				"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput				"User is blocked or deleted; Master profile is disabled"
//	@Failure		404			{object}	v0.ErrorOutput				"Master profile not found"
//	@Failure		500			{object}	v0.ErrorOutput				"Internal server error"
//	@Router			/v0/masters/profiles/{username}/services [get]
func (h *handler) FindMasterServicesByUsername(ctx *gin.Context) {
	log.Debug().Msg("handle find master services")

	input := v0.FindMasterServicesInput{}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")

	if principal.Username == username {
		resp, err := h.svc.FindAllMasterServicesByUsername(ctx, username, input)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrUnavailableSortField):
				_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
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

	resp, err := h.svc.FindAllOwnMasterServices(ctx, principal.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnavailableSortField):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
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
}

// GetMasterServiceByUsername godoc
//
//	@Id				GetMasterServiceByUsername
//	@Summary		Get master service
//	@Description	Request for getting master service. User must be logged in. In response will be returned found master service.
//	@Security		BearerAuth
//	@Tags			Master Service API
//	@Accept			application/json
//	@Produce		application/json
//	@Param			username	path		string						true	"Username"
//	@Param			id			path		string						true	"Master profile ID"
//	@Success		200			{object}	v0.MasterServiceOutput	"Found master service"
//	@Failure		401			{object}	v0.ErrorOutput			"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput			"User is blocked or deleted; Master profile is disabled"
//	@Failure		404			{object}	v0.ErrorOutput			"Master service not found; Master profile not found"
//	@Failure		500			{object}	v0.ErrorOutput			"Internal server error"
//	@Router			/v0/masters/profiles/{username}/services/{id} [get]
func (h *handler) GetMasterServiceByUsername(ctx *gin.Context) {
	log.Debug().Msg("handle get master service")

	principal, err := middleware.GetAuthUser(ctx)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusUnauthorized, err)
		return
	}

	username := ctx.Param("username")

	masterServiceIDRaw := ctx.Param("id")
	masterServiceID, err := uuid.Parse(masterServiceIDRaw)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		return
	}

	if principal.Username == username {
		resp, err := h.svc.GetOwnMasterServiceByID(ctx, principal.ID, masterServiceID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrMasterProfileNotExist),
				errors.Is(err, domain.ErrMasterServiceNotExist):
				_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
			case errors.Is(err, domain.ErrMasterProfileDisabled):
				_ = util.ErrorWithStatus(ctx, http.StatusForbidden, err)
			default:
				_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
			}
			return
		}

		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp, err := h.svc.GetMasterServiceByID(ctx, username, masterServiceID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMasterProfileNotExist),
			errors.Is(err, domain.ErrMasterServiceNotExist):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
