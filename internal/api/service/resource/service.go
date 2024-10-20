package resource

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	dto2 "mandarine/internal/api/service/resource/dto"
	"mandarine/pkg/logging"
	dto3 "mandarine/pkg/rest/dto"
	"mandarine/pkg/storage/s3"
	"mandarine/pkg/storage/s3/dto"
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
	slog.Info("Upload resource")
	factoryErr := func(err error) (dto2.UploadResourceOutput, error) {
		slog.Error("Upload resource error", logging.ErrorAttr(err))
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
			slog.Warn("Upload resource error: File close error", logging.ErrorAttr(err))
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
	slog.Info("Upload resource")
	factoryErr := func(err error) (dto2.UploadResourcesOutput, error) {
		slog.Error("Upload resource error", logging.ErrorAttr(err))
		return dto2.UploadResourcesOutput{}, err
	}

	files := input.Resources

	// Open files
	fileDatas := make([]*dto.FileData, 0)
	defer func() {
		for _, fileData := range fileDatas {
			err := fileData.Reader.(multipart.File).Close()
			if err != nil {
				slog.Warn("Upload resource error: File close error", logging.ErrorAttr(err))
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
			slog.Error("Upload resource error", logging.ErrorAttr(createDto.Error))
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
	slog.Info("Download resource")
	getDto := s.minioClient.GetOne(ctx, objectID)
	return getDto.Data, getDto.Error
}

////////// Helpers //////////

func calculateHash(f multipart.File) (string, error) {
	tmpFile, err := os.CreateTemp("", "tmp_")
	if err != nil {
		return "", err
	}
	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			slog.Warn("Upload resource error: File remove error", logging.ErrorAttr(err))
		}
	}()

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
