package v0

import "mime/multipart"

type UploadResourceInput struct {
	Resource *multipart.FileHeader `form:"redis" binding:"required"`
}

type UploadResourceOutput struct {
	ObjectID string `json:"object_id"`
}

type UploadResourcesInput struct {
	Resources []*multipart.FileHeader `form:"resources" binding:"required"`
}

type UploadResourcesOutput struct {
	Count int                             `json:"count"`
	Data  map[string]UploadResourceOutput `json:"data"`
}
