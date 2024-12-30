package minio

import (
	"context"
	"fmt"
	"github.com/mandarine-io/baselib/pkg/storage/s3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"log/slog"
	"sync"
)

type Config struct {
	Address    string
	AccessKey  string
	SecretKey  string
	BucketName string
}

type client struct {
	minio      *minio.Client
	bucketName string
}

func MustNewMinioClient(cfg *Config) s3.Client {
	// Configure to use MinIO Server
	ctx := context.Background()
	minioClient, err := minio.New(cfg.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to connect to minio")
	}
	log.Info().Msgf("connected to minio host %s", cfg.Address)

	// Check if bucket exists
	log.Info().Msgf("check bucket \"%s\"", cfg.BucketName)
	exists, err := minioClient.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to check minio bucket")
	}
	if !exists {
		slog.Info(fmt.Sprintf("Create bucket \"%s\"", cfg.BucketName))
		err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal().Stack().Err(err).Msg("failed to create minio bucket")
		}
	}

	return &client{minio: minioClient, bucketName: cfg.BucketName}
}

func (c *client) CreateOne(ctx context.Context, file *s3.FileData) *s3.CreateDto {
	log.Debug().Msg("create one object")
	if file == nil {
		return &s3.CreateDto{Error: errors.New("file is nil")}
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
		return &s3.CreateDto{Error: err}
	}
	return &s3.CreateDto{ObjectID: info.Key}
}

func (c *client) CreateMany(ctx context.Context, files []*s3.FileData) map[string]*s3.CreateDto {
	log.Debug().Msg("create many object")

	type entry struct {
		filename string
		dto      *s3.CreateDto
	}

	dtoCh := make(chan *entry, len(files))
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dtoCh <- &entry{filename: file.UserMetadata[s3.OriginalFilenameMetadata], dto: c.CreateOne(ctx, file)}
		}()
	}

	go func() {
		wg.Wait()
		close(dtoCh)
	}()

	dtoMap := make(map[string]*s3.CreateDto)
	for entry := range dtoCh {
		dtoMap[entry.filename] = entry.dto
	}

	return dtoMap
}

func (c *client) GetOne(ctx context.Context, objectID string) *s3.GetDto {
	log.Debug().Msg("get one object")

	object, err := c.minio.GetObject(ctx, c.bucketName, objectID, minio.GetObjectOptions{})
	if err != nil {
		if errors.As(err, &minio.ErrorResponse{}) && err.(minio.ErrorResponse).Code == "NoSuchKey" {
			return &s3.GetDto{Error: s3.ErrObjectNotFound}
		}
		return &s3.GetDto{Error: err}
	}
	if object == nil {
		return &s3.GetDto{Error: s3.ErrObjectNotFound}
	}

	stat, err := object.Stat()
	if err != nil {
		if errors.As(err, &minio.ErrorResponse{}) && err.(minio.ErrorResponse).Code == "NoSuchKey" {
			return &s3.GetDto{Error: s3.ErrObjectNotFound}
		}
		return &s3.GetDto{Error: err}
	}

	return &s3.GetDto{
		Data: &s3.FileData{
			Reader:      object,
			ID:          stat.Key,
			Size:        stat.Size,
			ContentType: stat.ContentType,
		},
	}
}

func (c *client) GetMany(ctx context.Context, objectIDs []string) map[string]*s3.GetDto {
	log.Debug().Msg("get many object")

	type entry struct {
		objectID string
		dto      *s3.GetDto
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

	dtoMap := make(map[string]*s3.GetDto)
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
