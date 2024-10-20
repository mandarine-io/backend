package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	recourceSvc "mandarine/internal/api/service/resource"
	dto2 "mandarine/internal/api/service/resource/dto"
	"mandarine/pkg/logging"
	"mandarine/pkg/rest/dto"
	"mandarine/pkg/rest/middleware"
	"mandarine/pkg/storage/s3"
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

func (h *Handler) RegisterRoutes(router *gin.Engine, requireAuth middleware.RequireAuth, _ middleware.RequireRoleFactory) {
	router.GET("v0/resources/:objectId", h.DownloadResource)

	router.POST("v0/resources/one", requireAuth, h.UploadResource)
	router.POST("v0/resources/many", requireAuth, h.UploadResources)
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
	objectId := ctx.Param("objectId")

	data, err := h.svc.DownloadResource(ctx, objectId)
	defer func() {
		if data == nil {
			return
		}
		err := data.Reader.Close()
		if err != nil {
			slog.Warn("Get resource error: File close error", logging.ErrorAttr(err))
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
