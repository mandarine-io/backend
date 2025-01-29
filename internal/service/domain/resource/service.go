package resource

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog"
	"io"
	"mime/multipart"
	"os"
)

type svc struct {
	s3Manager s3.Manager
	logger    zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(s3Manager s3.Manager, opts ...Option) domain.ResourceService {
	s := &svc{
		s3Manager: s3Manager,
		logger:    zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

////////// Upload resource //////////

func (s *svc) UploadResource(ctx context.Context, input *v0.UploadResourceInput) (
	v0.UploadResourceOutput,
	error,
) {
	s.logger.Info().Msg("upload resource")

	file := input.Resource

	// File is nil
	if file == nil {
		s.logger.Error().Stack().Err(domain.ErrResourceNotUploaded).Msg("failed to upload resource")
		return v0.UploadResourceOutput{}, domain.ErrResourceNotUploaded
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to upload resource")
		return v0.UploadResourceOutput{}, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to close file")
		}
	}()

	// Calculate hash
	hash, err := s.calculateHash(f)
	if err != nil {
		s.logger.Error().Stack().Err(err).Msg("failed to upload resource")
		return v0.UploadResourceOutput{}, err
	}

	// Upload to S3
	fileData := &s3.FileData{
		Reader:      f,
		ID:          fmt.Sprintf("%s-%s", hash, file.Filename),
		Size:        file.Size,
		ContentType: file.Header.Get("Content-Type"),
		UserMetadata: map[string]string{
			s3.OriginalFilenameMetadata: file.Filename,
		},
	}
	createDto := s.s3Manager.CreateOne(ctx, fileData)

	// Map output
	output := v0.UploadResourceOutput{
		ObjectID: createDto.ObjectID,
	}
	return output, createDto.Error
}

func (s *svc) UploadResources(ctx context.Context, input *v0.UploadResourcesInput) (
	v0.UploadResourcesOutput,
	error,
) {
	s.logger.Info().Msg("upload resources")

	files := input.Resources

	// Open files
	fileDatas := make([]*s3.FileData, 0)
	defer func() {
		for _, fileData := range fileDatas {
			err := fileData.Reader.(multipart.File).Close()
			if err != nil {
				s.logger.Warn().Err(err).Msg("failed to close file")
			}
		}
	}()

	for _, file := range files {
		// File is nil
		if file == nil {
			continue
		}

		f, err := file.Open()
		if err != nil {
			s.logger.Error().Stack().Err(err).Msg("failed to upload resources")
			return v0.UploadResourcesOutput{}, err
		}

		// Calculate hash
		hash, err := s.calculateHash(f)
		if err != nil {
			return v0.UploadResourcesOutput{}, err
		}

		fileData := &s3.FileData{
			Reader:      f,
			ID:          fmt.Sprintf("%s-%s", hash, file.Filename),
			Size:        file.Size,
			ContentType: file.Header.Get("Content-Type"),
			UserMetadata: map[string]string{
				s3.OriginalFilenameMetadata: file.Filename,
			},
		}
		fileDatas = append(fileDatas, fileData)
	}

	if len(fileDatas) == 0 {
		return v0.UploadResourcesOutput{Count: 0, Data: make(map[string]v0.UploadResourceOutput)}, nil
	}

	// Upload to S3
	createDtoMap := s.s3Manager.CreateMany(ctx, fileDatas)

	// Map output
	data := make(map[string]v0.UploadResourceOutput)
	for fileName, createDto := range createDtoMap {
		if createDto.Error != nil {
			s.logger.Error().Stack().Err(createDto.Error).Msg("failed to upload resource")
			continue
		}
		data[fileName] = v0.UploadResourceOutput{
			ObjectID: createDto.ObjectID,
		}
	}

	return v0.UploadResourcesOutput{Count: len(data), Data: data}, nil
}

////////// Download resource //////////

func (s *svc) DownloadResource(ctx context.Context, objectID string) (*s3.FileData, error) {
	s.logger.Info().Msg("download resource")
	getDto := s.s3Manager.GetOne(ctx, objectID)
	return getDto.Data, getDto.Error
}

////////// Helpers //////////

func (s *svc) calculateHash(f multipart.File) (string, error) {
	// Create temp file
	s.logger.Debug().Msg("create temp file")
	tmpFile, err := os.CreateTemp("", "tmp_")
	if err != nil {
		return "", err
	}
	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to remove temp file")
		}
	}()

	// Calculate hash
	s.logger.Debug().Msg("calculate hash during write to temp file")
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
