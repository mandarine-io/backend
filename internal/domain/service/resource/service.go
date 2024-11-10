package resource

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"os"
)

type svc struct {
	minioClient s3.Client
}

func NewService(minioClient s3.Client) service.ResourceService {
	return &svc{minioClient: minioClient}
}

////////// Upload resource //////////

func (s *svc) UploadResource(ctx context.Context, input *dto.UploadResourceInput) (dto.UploadResourceOutput, error) {
	log.Info().Msg("upload resource")

	file := input.Resource

	// File is nil
	if file == nil {
		log.Error().Stack().Err(service.ErrResourceNotUploaded).Msg("failed to upload resource")
		return dto.UploadResourceOutput{}, service.ErrResourceNotUploaded
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to upload resource")
		return dto.UploadResourceOutput{}, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close file")
		}
	}()

	// Calculate hash
	hash, err := calculateHash(f)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to upload resource")
		return dto.UploadResourceOutput{}, err
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
	createDto := s.minioClient.CreateOne(ctx, fileData)

	// Map output
	output := dto.UploadResourceOutput{
		ObjectID: createDto.ObjectID,
	}
	return output, createDto.Error
}

func (s *svc) UploadResources(ctx context.Context, input *dto.UploadResourcesInput) (dto.UploadResourcesOutput, error) {
	log.Info().Msg("upload resources")

	files := input.Resources

	// Open files
	fileDatas := make([]*s3.FileData, 0)
	defer func() {
		for _, fileData := range fileDatas {
			err := fileData.Reader.(multipart.File).Close()
			if err != nil {
				log.Warn().Err(err).Msg("failed to close file")
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
			log.Error().Stack().Err(err).Msg("failed to upload resources")
			return dto.UploadResourcesOutput{}, err
		}

		// Calculate hash
		hash, err := calculateHash(f)
		if err != nil {
			return dto.UploadResourcesOutput{}, err
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
		return dto.UploadResourcesOutput{Count: 0, Data: make(map[string]dto.UploadResourceOutput)}, nil
	}

	// Upload to S3
	createDtoMap := s.minioClient.CreateMany(ctx, fileDatas)

	// Map output
	data := make(map[string]dto.UploadResourceOutput)
	for fileName, createDto := range createDtoMap {
		if createDto.Error != nil {
			log.Error().Stack().Err(createDto.Error).Msg("failed to upload resource")
			continue
		}
		data[fileName] = dto.UploadResourceOutput{
			ObjectID: createDto.ObjectID,
		}
	}

	return dto.UploadResourcesOutput{Count: len(data), Data: data}, nil
}

////////// Download resource //////////

func (s *svc) DownloadResource(ctx context.Context, objectID string) (*s3.FileData, error) {
	log.Info().Msg("download resource")
	getDto := s.minioClient.GetOne(ctx, objectID)
	return getDto.Data, getDto.Error
}

////////// Helpers //////////

func calculateHash(f multipart.File) (string, error) {
	// Create temp file
	log.Debug().Msg("create temp file")
	tmpFile, err := os.CreateTemp("", "tmp_")
	if err != nil {
		return "", err
	}
	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			log.Warn().Err(err).Msg("failed to remove temp file")
		}
	}()

	// Calculate hash
	log.Debug().Msg("calculate hash during write to temp file")
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
