package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	recourceSvc "github.com/mandarine-io/Backend/internal/api/service/resource"
	dto2 "github.com/mandarine-io/Backend/internal/api/service/resource/dto"
	"github.com/mandarine-io/Backend/internal/api/transport/http/handler"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	"github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

var (
	ErrResourceNotUploaded = dto.NewI18nError("resource not uploaded", "errors.resource_not_uploaded")
)

type Handler struct {
	svc *recourceSvc.Service
}

func NewHandler(svc *recourceSvc.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, middlewares handler.RouteMiddlewares) {
	log.Debug().Msg("register resource routes")

	router.GET("v0/resources/:objectId", h.DownloadResource)

	router.POST(
		"v0/resources/one",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.UploadResource,
	)
	router.POST(
		"v0/resources/many",
		middlewares.Auth,
		middlewares.BannedUser,
		middlewares.DeletedUser,
		h.UploadResources,
	)
}

// UploadResource godoc
//
//	@Id				UploadResource
//	@Summary		Upload resource
//	@Description	Request for uploading resource. Return the object id in S3 storage.
//	@Tags			Resource API
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			resource formData		file	true	"File to upload"
//	@Success		201		{object}	dto.UploadResourceOutput
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v0/resources/one [post]
func (h *Handler) UploadResource(ctx *gin.Context) {
	log.Debug().Msg("handle upload resource")

	var req dto2.UploadResourceInput
	if err := ctx.ShouldBind(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrResourceNotUploaded)
		return
	}

	res, err := h.svc.UploadResource(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, recourceSvc.ErrResourceNotUploaded):
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

// UploadResources godoc
//
//	@Id				UploadResources
//	@Summary		Upload resources
//	@Description	Request for uploading resources. Return the array of object ids in S3 storage for successful uploaded files.
//	@Tags			Resource API
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			resources	formData	[]file	true	"Files to upload"
//	@Success		201		{object}	dto.UploadResourcesOutput
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v0/resources/many [post]
func (h *Handler) UploadResources(ctx *gin.Context) {
	log.Debug().Msg("handle upload resources")

	var req dto2.UploadResourcesInput
	if err := ctx.ShouldBind(&req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, ErrResourceNotUploaded)
		return
	}

	res, err := h.svc.UploadResources(ctx, &req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

// DownloadResource godoc
//
//	@Id				DownloadResource
//	@Summary		Download resource
//	@Description	Request for getting resource. Return the resource in S3 storage.
//	@Tags			Resource API
//	@Produce		*/*
//	@Param			objectId	path		string	true	"Object id"
//	@Success		200
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/v0/resources/{objectId} [get]
func (h *Handler) DownloadResource(ctx *gin.Context) {
	log.Debug().Msg("handle download resource")

	objectId := ctx.Param("objectId")

	data, err := h.svc.DownloadResource(ctx, objectId)
	defer func() {
		if data == nil {
			return
		}
		err := data.Reader.Close()
		if err != nil {
			log.Warn().Stack().Err(err).Msg("failed to close file")
		}
	}()
	if err != nil {
		switch {
		case errors.Is(err, s3.ErrObjectNotFound):
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		default:
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.DataFromReader(
		http.StatusOK, data.Size, data.ContentType, data.Reader,
		map[string]string{"Content-Dispositon": "attachment; filename=" + data.ID},
	)
}
