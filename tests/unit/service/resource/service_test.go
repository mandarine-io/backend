package resource_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mandarine-io/Backend/internal/api/service/resource"
	"github.com/mandarine-io/Backend/internal/api/service/resource/dto"
	dto2 "github.com/mandarine-io/Backend/pkg/storage/s3/dto"
	mock2 "github.com/mandarine-io/Backend/pkg/storage/s3/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

var (
	ctx      = context.Background()
	s3Client *mock2.S3ClientMock
	svc      *resource.Service
)

func TestMain(m *testing.M) {
	// Setup logger
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.Level(10000),
			},
		),
	)
	slog.SetDefault(logger)

	// Setup mocks
	s3Client = new(mock2.S3ClientMock)
	svc = resource.NewService(s3Client)

	os.Exit(m.Run())
}

func Test_ResourceService_UploadResource(t *testing.T) {
	t.Run("File is nil", func(t *testing.T) {
		output, err := svc.UploadResource(ctx, &dto.UploadResourceInput{Resource: nil})

		assert.Equal(t, resource.ErrResourceNotUploaded, err)
		assert.Equal(t, dto.UploadResourceOutput{}, output)
	})

	t.Run("Error upload file", func(t *testing.T) {
		// Create temp file
		file, err := os.CreateTemp("", "test-")
		require.Nil(t, err)
		defer os.Remove(file.Name())

		// Create FileHeader
		fileHeader := createMultipartFileHeader(file.Name())
		require.NotNil(t, fileHeader)

		// Mock
		expectedErr := errors.New("s3 error")
		s3Client.On("CreateOne", ctx, mock.Anything).Return(&dto2.CreateDto{Error: expectedErr}).Once()

		// Call service
		output, err := svc.UploadResource(ctx, &dto.UploadResourceInput{Resource: fileHeader})

		assert.Equal(t, err, expectedErr)
		assert.Equal(t, dto.UploadResourceOutput{}, output)
	})

	t.Run("Success", func(t *testing.T) {
		// Create temp file
		file, err := os.CreateTemp("", "test-")
		require.Nil(t, err)
		defer os.Remove(file.Name())

		// Create FileHeader
		fileHeader := createMultipartFileHeader(file.Name())
		require.NotNil(t, fileHeader)

		// Mock
		s3Client.On("CreateOne", ctx, mock.Anything).Return(&dto2.CreateDto{ObjectID: "test"}).Once()

		// Call service
		output, err := svc.UploadResource(ctx, &dto.UploadResourceInput{Resource: fileHeader})

		assert.Nil(t, err)
		assert.Equal(t, output.ObjectID, "test")
	})
}

func Test_ResourceService_UploadResources(t *testing.T) {
	t.Run("File is nil", func(t *testing.T) {
		fileHeaders := make([]*multipart.FileHeader, 3)

		output, err := svc.UploadResources(ctx, &dto.UploadResourcesInput{Resources: fileHeaders})

		assert.NoError(t, err)
		assert.Equal(t, 0, output.Count)
		assert.Equal(t, map[string]dto.UploadResourceOutput{}, output.Data)
	})

	t.Run("Upload error", func(t *testing.T) {
		var err error
		files := make([]*os.File, 3)
		fileHeaders := make([]*multipart.FileHeader, 3)
		for i := 0; i < 3; i++ {
			// Create temp file
			files[i], err = os.CreateTemp("", fmt.Sprintf("test-%d-", i))
			require.Nil(t, err)

			// Create FileHeader
			fileHeaders[i] = createMultipartFileHeader(files[i].Name())
			require.NotNil(t, fileHeaders[i])
		}
		defer func() {
			for _, file := range files {
				os.Remove(file.Name())
			}
		}()

		t.Run("All files error", func(t *testing.T) {
			// Mock
			s3Output := make(map[string]*dto2.CreateDto)
			for i := 0; i < 3; i++ {
				s3Output[files[i].Name()] = &dto2.CreateDto{Error: errors.New("s3 error")}
			}
			s3Client.On("CreateMany", ctx, mock.Anything).Return(s3Output).Once()

			output, err := svc.UploadResources(ctx, &dto.UploadResourcesInput{Resources: fileHeaders})

			assert.NoError(t, err)
			assert.Equal(t, 0, output.Count)
			assert.Equal(t, map[string]dto.UploadResourceOutput{}, output.Data)
		})

		t.Run("Have one files error", func(t *testing.T) {
			// Mock
			s3Output := make(map[string]*dto2.CreateDto)
			for i := 0; i < 3; i++ {
				if i == 0 {
					s3Output[files[i].Name()] = &dto2.CreateDto{ObjectID: "test", Error: nil}
					continue
				}
				s3Output[files[i].Name()] = &dto2.CreateDto{Error: errors.New("s3 error")}
			}
			s3Client.On("CreateMany", ctx, mock.Anything).Return(s3Output).Once()

			output, err := svc.UploadResources(ctx, &dto.UploadResourcesInput{Resources: fileHeaders})

			assert.NoError(t, err)
			assert.Equal(t, 1, output.Count)

			resourceOutput, ok := output.Data[files[0].Name()]
			assert.True(t, ok)
			assert.Equal(t, "test", resourceOutput.ObjectID)
		})

		t.Run("No files error", func(t *testing.T) {
			// Mock
			s3Output := make(map[string]*dto2.CreateDto)
			for i := 0; i < 3; i++ {
				s3Output[files[i].Name()] = &dto2.CreateDto{ObjectID: fmt.Sprintf("test-%d", i), Error: nil}
			}
			s3Client.On("CreateMany", ctx, mock.Anything).Return(s3Output).Once()

			output, err := svc.UploadResources(ctx, &dto.UploadResourcesInput{Resources: fileHeaders})

			assert.NoError(t, err)
			assert.Equal(t, 3, output.Count)

			for i := 0; i < 3; i++ {
				resourceOutput, ok := output.Data[files[i].Name()]
				assert.True(t, ok)
				assert.Equal(t, fmt.Sprintf("test-%d", i), resourceOutput.ObjectID)
			}
		})
	})
}

func Test_ResourceService_DownloadResource(t *testing.T) {
	t.Run("Not found", func(t *testing.T) {
		// Mock
		expectError := errors.New("s3 error")
		s3Client.On("GetOne", ctx, "test").Return(&dto2.GetDto{Data: nil, Error: expectError}).Once()

		output, err := svc.DownloadResource(ctx, "test")

		assert.Error(t, err)
		assert.Equal(t, expectError, err)
		assert.Nil(t, output)
	})

	t.Run("Success", func(t *testing.T) {
		// Mock
		s3Client.On("GetOne", ctx, "test").Return(&dto2.GetDto{Data: &dto2.FileData{}, Error: nil}).Once()

		output, err := svc.DownloadResource(ctx, "test")

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, &dto2.FileData{}, output)
	})
}

func createMultipartFileHeader(filePath string) *multipart.FileHeader {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		return nil
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil
	}

	return files[0]
}
