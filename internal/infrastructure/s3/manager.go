package s3

import (
	"context"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"io"
)

const (
	OriginalFilenameMetadata = "x-amz-meta-original-filename"
)

var (
	ErrObjectNotFound = v0.NewI18nError("object not found", "errors.object_not_found")
)

type (
	FileData struct {
		ID           string
		Size         int64
		ContentType  string
		Reader       io.ReadCloser
		UserMetadata map[string]string
	}

	CreateResult struct {
		ObjectID string
		Error    error
	}

	GetResult struct {
		Data  *FileData
		Error error
	}

	Manager interface {
		CreateOne(ctx context.Context, file *FileData) CreateResult
		CreateMany(ctx context.Context, files []*FileData) map[string]CreateResult
		GetOne(ctx context.Context, objectID string) GetResult
		GetMany(ctx context.Context, objectIDs []string) map[string]GetResult
		DeleteOne(ctx context.Context, objectID string) error
		DeleteMany(ctx context.Context, objectIDs []string) map[string]error
	}
)
