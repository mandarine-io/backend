package minio

import (
	"context"
	"errors"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"sync"
)

type Option func(*manager) error

func WithLogger(logger zerolog.Logger) Option {
	return func(m *manager) error {
		m.logger = logger
		return nil
	}
}

type manager struct {
	client     *minio.Client
	logger     zerolog.Logger
	bucketName string
}

func NewManager(client *minio.Client, bucketName string, opts ...Option) (s3.Manager, error) {
	m := &manager{
		client:     client,
		logger:     zerolog.Nop(),
		bucketName: bucketName,
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return m, nil
}

func (m *manager) CreateOne(ctx context.Context, file *s3.FileData) s3.CreateResult {
	m.logger.Debug().Msg("create one object")

	// Check nil
	if file == nil {
		return s3.CreateResult{Error: errors.New("file data is nil")}
	}

	// Upload
	info, err := m.client.PutObject(
		ctx, m.bucketName, file.ID, file.Reader, file.Size,
		minio.PutObjectOptions{
			SendContentMd5:        true,
			PartSize:              10 * 1024 * 1024,
			ConcurrentStreamParts: true,
			ContentType:           file.ContentType,
			UserMetadata:          file.UserMetadata,
		},
	)
	if err != nil {
		return s3.CreateResult{Error: err}
	}
	return s3.CreateResult{ObjectID: info.Key}
}

func (m *manager) CreateMany(ctx context.Context, files []*s3.FileData) map[string]s3.CreateResult {
	m.logger.Debug().Msg("create many object")

	type entry struct {
		filename string
		model    s3.CreateResult
	}

	modelCh := make(chan *entry, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			modelCh <- &entry{filename: file.UserMetadata[s3.OriginalFilenameMetadata], model: m.CreateOne(ctx, file)}
		}()
	}

	go func() {
		wg.Wait()
		close(modelCh)
	}()

	modelMap := make(map[string]s3.CreateResult)
	for e := range modelCh {
		modelMap[e.filename] = e.model
	}

	return modelMap
}

func (m *manager) GetOne(ctx context.Context, objectID string) s3.GetResult {
	m.logger.Debug().Msg("get one object")

	object, err := m.client.GetObject(ctx, m.bucketName, objectID, minio.GetObjectOptions{})
	if err != nil {
		minioErr := minio.ErrorResponse{}
		if errors.As(err, &minioErr) && minioErr.Code == "NoSuchKey" {
			return s3.GetResult{Error: s3.ErrObjectNotFound}
		}

		return s3.GetResult{Error: err}
	}
	if object == nil {
		return s3.GetResult{Error: s3.ErrObjectNotFound}
	}

	stat, err := object.Stat()
	if err != nil {
		minioErr := minio.ErrorResponse{}
		if errors.As(err, &minioErr) && minioErr.Code == "NoSuchKey" {
			return s3.GetResult{Error: s3.ErrObjectNotFound}
		}

		return s3.GetResult{Error: err}
	}

	return s3.GetResult{
		Data: &s3.FileData{
			Reader:      object,
			ID:          stat.Key,
			Size:        stat.Size,
			ContentType: stat.ContentType,
		},
	}
}

func (m *manager) GetMany(ctx context.Context, objectIDs []string) map[string]s3.GetResult {
	m.logger.Debug().Msg("get many object")

	type entry struct {
		objectID string
		model    s3.GetResult
	}

	modelCh := make(chan *entry, len(objectIDs))
	var wg sync.WaitGroup

	for _, objectID := range objectIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			modelCh <- &entry{objectID: objectID, model: m.GetOne(ctx, objectID)}
		}()
	}

	go func() {
		wg.Wait()
		close(modelCh)
	}()

	modelMap := make(map[string]s3.GetResult)
	for e := range modelCh {
		modelMap[e.objectID] = e.model
	}

	return modelMap
}

func (m *manager) DeleteOne(ctx context.Context, objectID string) error {
	m.logger.Debug().Msg("delete one object")
	return m.client.RemoveObject(ctx, m.bucketName, objectID, minio.RemoveObjectOptions{})
}

func (m *manager) DeleteMany(ctx context.Context, objectIDs []string) map[string]error {
	m.logger.Debug().Msg("delete many object")
	objectIDCh := make(chan minio.ObjectInfo, len(objectIDs))
	for _, objectID := range objectIDs {
		objectIDCh <- minio.ObjectInfo{Key: objectID}
	}
	close(objectIDCh)

	objCh := m.client.RemoveObjects(ctx, m.bucketName, objectIDCh, minio.RemoveObjectsOptions{})

	errMap := make(map[string]error)
	for obj := range objCh {
		errMap[obj.ObjectName] = obj.Err
	}

	return errMap
}
