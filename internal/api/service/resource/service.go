package resource

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	dto2 "github.com/mandarine-io/Backend/internal/api/service/resource/dto"
	dto3 "github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	"github.com/mandarine-io/Backend/pkg/storage/s3/dto"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"os"
)

var (
	ErrResourceNotUploaded = dto3.NewI18nError("resource not uploaded", "errors.resource_not_uploaded")
)

type Service struct {
	minioClient s3.Client
}

func NewService(minioClient s3.Client) *Service {
	return &Service{
		minioClient: minioClient,
	}
}

////////// Upload resource //////////

func (s *Service) UploadResource(ctx context.Context, input *dto2.UploadResourceInput) (dto2.UploadResourceOutput, error) {
	log.Info().Msg("upload resource")
	factoryErr := func(err error) (dto2.UploadResourceOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to upload resource")
		return dto2.UploadResourceOutput{}, err
	}

	file := input.Resource

	// File is nil
	if file == nil {
		return dto2.UploadResourceOutput{}, ErrResourceNotUploaded
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		return dto2.UploadResourceOutput{}, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Warn().Stack().Err(err).Msg("failed to close file")
		}
	}()

	// Calculate hash
	hash, err := calculateHash(f)
	if err != nil {
		return factoryErr(err)
	}

	// Upload to S3
	fileData := &dto.FileData{
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
	output := dto2.UploadResourceOutput{
		ObjectID: createDto.ObjectID,
	}
	return output, createDto.Error
}

func (s *Service) UploadResources(ctx context.Context, input *dto2.UploadResourcesInput) (dto2.UploadResourcesOutput, error) {
	log.Info().Msg("upload resources")
	factoryErr := func(err error) (dto2.UploadResourcesOutput, error) {
		log.Error().Stack().Err(err).Msg("failed to upload resources")
		return dto2.UploadResourcesOutput{}, err
	}

	files := input.Resources

	// Open files
	fileDatas := make([]*dto.FileData, 0)
	defer func() {
		for _, fileData := range fileDatas {
			err := fileData.Reader.(multipart.File).Close()
			if err != nil {
				log.Warn().Stack().Err(err).Msg("failed to close file")
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
			return factoryErr(err)
		}

		// Calculate hash
		hash, err := calculateHash(f)
		if err != nil {
			return dto2.UploadResourcesOutput{}, err
		}

		fileData := &dto.FileData{
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
		return dto2.UploadResourcesOutput{Count: 0, Data: make(map[string]dto2.UploadResourceOutput)}, nil
	}

	// Upload to S3
	createDtoMap := s.minioClient.CreateMany(ctx, fileDatas)

	// Map output
	data := make(map[string]dto2.UploadResourceOutput)
	for fileName, createDto := range createDtoMap {
		if createDto.Error != nil {
			log.Error().Stack().Err(createDto.Error).Msg("failed to upload resource")
			continue
		}
		data[fileName] = dto2.UploadResourceOutput{
			ObjectID: createDto.ObjectID,
		}
	}

	return dto2.UploadResourcesOutput{Count: len(data), Data: data}, nil
}

////////// Download resource //////////

func (s *Service) DownloadResource(ctx context.Context, objectID string) (*dto.FileData, error) {
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
			log.Warn().Stack().Err(err).Msg("failed to remove temp file")
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
