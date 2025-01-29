package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/mandarine-io/backend/internal/service/domain"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/internal/transport/http/middleware"
	"github.com/mandarine-io/backend/internal/transport/http/util"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

var (
	ErrResourceNotUploaded = v0.NewI18nError("resource not uploaded", "errors.resource_not_uploaded")
)

type handler struct {
	svc    domain.ResourceService
	logger zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(svc domain.ResourceService, opts ...Option) apihandler.APIHandler {
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
	log.Debug().Msg("register resource routes")

	router.GET("v0/resources/:objectID", h.DownloadResource)

	router.POST(
		"v0/resources/one",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
		h.UploadResource,
	)
	router.POST(
		"v0/resources/many",
		middleware.Registry.Auth,
		middleware.Registry.BannedUser,
		middleware.Registry.DeletedUser,
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
//	@Produce		application/json
//	@Param			resource	formData	string						true	"File to upload"
//	@Success		201			{object}	v0.UploadResourceOutput	"Uploaded resource"
//	@Failure		400			{object}	v0.ErrorOutput			"Validation error"
//	@Failure		401			{object}	v0.ErrorOutput			"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput			"User is blocked or deleted"
//	@Failure		500			{object}	v0.ErrorOutput			"Internal server error"
//	@Router			/v0/resources/one [post]
func (h *handler) UploadResource(ctx *gin.Context) {
	log.Debug().Msg("handle upload resource")

	var input v0.UploadResourceInput
	if err := ctx.ShouldBind(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, ErrResourceNotUploaded)
		return
	}

	res, err := h.svc.UploadResource(ctx, &input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourceNotUploaded):
			_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
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
//	@Produce		application/json
//	@Param			resources	formData	[]string						true	"Files to upload"
//	@Success		201			{object}	v0.UploadResourcesOutput	"Uploaded resources"
//	@Failure		400			{object}	v0.ErrorOutput			"Validation error"
//	@Failure		401			{object}	v0.ErrorOutput			"Unauthorized error"
//	@Failure		403			{object}	v0.ErrorOutput			"User is blocked or deleted"
//	@Failure		500			{object}	v0.ErrorOutput			"Internal server error"
//	@Router			/v0/resources/many [post]
func (h *handler) UploadResources(ctx *gin.Context) {
	log.Debug().Msg("handle upload resources")

	var input v0.UploadResourcesInput
	if err := ctx.ShouldBind(&input); err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusBadRequest, ErrResourceNotUploaded)
		return
	}

	res, err := h.svc.UploadResources(ctx, &input)
	if err != nil {
		_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
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
//	@Param			objectID	path	string	true	"Object id"
//	@Success		200
//	@Failure		404	{object}	v0.ErrorOutput	"Resource not found"
//	@Failure		500	{object}	v0.ErrorOutput	"Internal server error"
//	@Router			/v0/resources/{objectID} [get]
func (h *handler) DownloadResource(ctx *gin.Context) {
	log.Debug().Msg("handle download resource")

	objectID := ctx.Param("objectID")

	data, err := h.svc.DownloadResource(ctx, objectID)
	defer func() {
		if data == nil {
			return
		}
		err := data.Reader.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close file")
		}
	}()
	if err != nil {
		switch {
		case errors.Is(err, s3.ErrObjectNotFound):
			_ = util.ErrorWithStatus(ctx, http.StatusNotFound, err)
		default:
			_ = util.ErrorWithStatus(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if isASCII(data.ID) {
		ctx.Writer.Header().Set("Content-Disposition", `attachment; filename="`+escapeQuotes(data.ID)+`"`)
	} else {
		ctx.Writer.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(data.ID))
	}

	ctx.DataFromReader(
		http.StatusOK, data.Size, data.ContentType, data.Reader,
		map[string]string{},
	)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func escapeQuotes(s string) string {
	quoteEscaper := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	return quoteEscaper.Replace(s)
}
