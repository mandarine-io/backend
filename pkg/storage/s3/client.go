package s3

import (
	"context"
	dto2 "github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
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

	client struct {
		minio      *minio.Client
		bucketName string
	}
)

func NewClient(minio *minio.Client, bucketName string) Client {
	return &client{minio: minio, bucketName: bucketName}
}

func (c *client) CreateOne(ctx context.Context, file *FileData) *CreateDto {
	log.Debug().Msg("create one object")
	if file == nil {
		return &CreateDto{Error: errors.New("file is nil")}
	}

	// Upload
	info, err := c.minio.PutObject(
		ctx, c.bucketName, file.ID, file.Reader, file.Size,
		minio.PutObjectOptions{
			SendContentMd5:        true,
			PartSize:              10 * 1024 * 1024,
			ConcurrentStreamParts: true,
			ContentType:           file.ContentType,
			UserMetadata:          file.UserMetadata,
		})
	if err != nil {
		return &CreateDto{Error: err}
	}
	return &CreateDto{ObjectID: info.Key}
}

func (c *client) CreateMany(ctx context.Context, files []*FileData) map[string]*CreateDto {
	log.Debug().Msg("create many object")

	type entry struct {
		filename string
		dto      *CreateDto
	}

	dtoCh := make(chan *entry, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dtoCh <- &entry{filename: file.UserMetadata[OriginalFilenameMetadata], dto: c.CreateOne(ctx, file)}
		}()
	}

	go func() {
		wg.Wait()
		close(dtoCh)
	}()

	dtoMap := make(map[string]*CreateDto)
	for entry := range dtoCh {
		dtoMap[entry.filename] = entry.dto
	}

	return dtoMap
}

func (c *client) GetOne(ctx context.Context, objectID string) *GetDto {
	log.Debug().Msg("get one object")

	object, err := c.minio.GetObject(ctx, c.bucketName, objectID, minio.GetObjectOptions{})
	if err != nil {
		if errors.As(err, &minio.ErrorResponse{}) && err.(minio.ErrorResponse).Code == "NoSuchKey" {
			return &GetDto{Error: ErrObjectNotFound}
		}
		return &GetDto{Error: err}
	}
	if object == nil {
		return &GetDto{Error: ErrObjectNotFound}
	}

	stat, err := object.Stat()
	if err != nil {
		if errors.As(err, &minio.ErrorResponse{}) && err.(minio.ErrorResponse).Code == "NoSuchKey" {
			return &GetDto{Error: ErrObjectNotFound}
		}
		return &GetDto{Error: err}
	}

	return &GetDto{
		Data: &FileData{
			Reader:      object,
			ID:          stat.Key,
			Size:        stat.Size,
			ContentType: stat.ContentType,
		},
	}
}

func (c *client) GetMany(ctx context.Context, objectIDs []string) map[string]*GetDto {
	log.Debug().Msg("get many object")

	type entry struct {
		objectID string
		dto      *GetDto
	}

	dtoCh := make(chan *entry, len(objectIDs))
	var wg sync.WaitGroup

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dtoCh <- &entry{objectID: objectID, dto: c.GetOne(ctx, objectID)}
		}()
	}

	go func() {
		wg.Wait()
		close(dtoCh)
	}()

	dtoMap := make(map[string]*GetDto)
	for entry := range dtoCh {
		dtoMap[entry.objectID] = entry.dto
	}

	return dtoMap
}

func (c *client) DeleteOne(ctx context.Context, objectID string) error {
	log.Debug().Msg("delete one object")
	return c.minio.RemoveObject(ctx, c.bucketName, objectID, minio.RemoveObjectOptions{})
}

func (c *client) DeleteMany(ctx context.Context, objectIDs []string) map[string]error {
	log.Debug().Msg("delete many object")
	objectIdCh := make(chan minio.ObjectInfo, len(objectIDs))
	for _, objectID := range objectIDs {
		objectIdCh <- minio.ObjectInfo{Key: objectID}
	}
	close(objectIdCh)

	objCh := c.minio.RemoveObjects(ctx, c.bucketName, objectIdCh, minio.RemoveObjectsOptions{})

	errMap := make(map[string]error)
	for obj := range objCh {
		errMap[obj.ObjectName] = obj.Err
	}

	return errMap
}
