package s3

import (
	"context"
	dto2 "github.com/mandarine-io/baselib/pkg/transport/http/model"
	"io"
)

const (
	OriginalFilenameMetadata = "x-amz-meta-original-filename"
)

var (
	ErrObjectNotFound = dto2.NewI18nError("object not found", "errors.object_not_found")
)

type (
	FileData struct {
		ID           string
		Size         int64
		ContentType  string
		Reader       io.ReadCloser
		UserMetadata map[string]string
	}

	CreateDto struct {
		ObjectID string
		Error    error
	}

	GetDto struct {
		Data  *FileData
		Error error
	}

	Client interface {
		CreateOne(ctx context.Context, file *FileData) *CreateDto
		CreateMany(ctx context.Context, files []*FileData) map[string]*CreateDto
		GetOne(ctx context.Context, objectID string) *GetDto
		GetMany(ctx context.Context, objectIDs []string) map[string]*GetDto
		DeleteOne(ctx context.Context, objectID string) error
		DeleteMany(ctx context.Context, objectIDs []string) map[string]error
	}
)
