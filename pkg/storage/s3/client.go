package s3

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"io"
	"strings"
	"sync"
)

type FileData struct {
	FileName    string
	FileSize    int64
	ContentType string
	Reader      io.Reader
}

type Client interface {
	CreateOne(ctx context.Context, file *FileData) (string, error)
	CreateMany(ctx context.Context, files []*FileData) ([]string, error)
	GetOne(ctx context.Context, objectID string) (io.Reader, error)
	GetMany(ctx context.Context, objectIDs []string) ([]io.Reader, error)
	DeleteOne(ctx context.Context, objectID string) error
	DeleteMany(ctx context.Context, objectIDs []string) error
}

type client struct {
	minio      *minio.Client
	bucketName string
}

type operationError struct {
	ObjectID string
	Error    error
}

func NewClient(minio *minio.Client, bucketName string) Client {
	return &client{minio: minio, bucketName: bucketName}
}

func (c *client) CreateOne(ctx context.Context, file *FileData) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	// Upload
	info, err := c.minio.PutObject(ctx, c.bucketName, file.FileName, file.Reader, file.FileSize, minio.PutObjectOptions{
		SendContentMd5:        true,
		ConcurrentStreamParts: true,
		ContentType:           file.ContentType,
	})
	if err != nil {
		return "", err
	}
	return info.Key, nil
}

func (c *client) CreateMany(ctx context.Context, files []*FileData) ([]string, error) {
	// Create channel and sync object
	objectIdCh := make(chan string, len(files))
	errCh := make(chan operationError, len(files))

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			objectId, err := c.CreateOne(ctx, file)
			if err != nil {
				errCh <- operationError{ObjectID: objectId, Error: err}
				cancel()
				return
			}

			objectIdCh <- objectId
		}()
	}

	go func() {
		wg.Wait()
		close(objectIdCh)
		close(errCh)
	}()

	objectIds := make([]string, 0, len(files))
	for objectId := range objectIdCh {
		objectIds = append(objectIds, objectId)
	}

	errs := make([]operationError, 0, len(files))
	for err := range errCh {
		errs = append(errs, err)
	}

	return objectIds, mapOperationErrorsToError(errs)
}

func (c *client) GetOne(ctx context.Context, objectID string) (io.Reader, error) {
	object, err := c.minio.GetObject(ctx, c.bucketName, objectID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	if object == nil {
		return nil, errors.New("object not found")
	}
	return object, nil
}

func (c *client) GetMany(ctx context.Context, objectIDs []string) ([]io.Reader, error) {
	readerCh := make(chan io.Reader, len(objectIDs))
	errCh := make(chan operationError, len(objectIDs))

	var wg sync.WaitGroup
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reader, err := c.GetOne(ctx, objectID)
			if err != nil {
				errCh <- operationError{ObjectID: objectID, Error: err}
				cancel()
				return
			}
			readerCh <- reader
		}()
	}

	go func() {
		wg.Wait()
		close(readerCh)
		close(errCh)
	}()

	readers := make([]io.Reader, 0, len(objectIDs))
	for reader := range readerCh {
		readers = append(readers, reader)
	}

	errs := make([]operationError, 0, len(objectIDs))
	for err := range errCh {
		errs = append(errs, err)
	}

	return readers, mapOperationErrorsToError(errs)
}

func (c *client) DeleteOne(ctx context.Context, objectID string) error {
	return c.minio.RemoveObject(ctx, c.bucketName, objectID, minio.RemoveObjectOptions{})
}

func (c *client) DeleteMany(ctx context.Context, objectIDs []string) error {
	objectIdCh := make(chan minio.ObjectInfo, len(objectIDs))

	for _, objectID := range objectIDs {
		objectIdCh <- minio.ObjectInfo{Key: objectID}
	}
	close(objectIdCh)

	errCh := c.minio.RemoveObjects(ctx, c.bucketName, objectIdCh, minio.RemoveObjectsOptions{})

	errs := make([]operationError, 0, len(objectIDs))
	for err := range errCh {
		errs = append(errs, operationError{ObjectID: err.ObjectName, Error: err.Err})
	}

	return mapOperationErrorsToError(errs)
}

func mapOperationErrorsToError(operationErrs []operationError) error {
	errMsg := strings.Builder{}
	for _, err := range operationErrs {
		if errMsg.Len() == 0 {
			errMsg.WriteString(err.ObjectID)
			errMsg.WriteString(": ")
			errMsg.WriteString(err.Error.Error())
			continue
		}
		errMsg.WriteString("; ")
		errMsg.WriteString(err.ObjectID)
		errMsg.WriteString(": ")
		errMsg.WriteString(err.Error.Error())
	}

	if errMsg.Len() == 0 {
		return nil
	}
	return errors.New(errMsg.String())
}
