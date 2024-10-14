package mock

import (
	"context"
	"io"
	"mandarine/pkg/storage/s3"

	"github.com/stretchr/testify/mock"
)

// S3ClientMock is a mock for the S3Client interface
type S3ClientMock struct {
	mock.Mock
}

// CreateOne mocks the CreateOne method
func (m *S3ClientMock) CreateOne(ctx context.Context, file *s3.FileData) (string, error) {
	args := m.Called(ctx, file)
	return args.String(0), args.Error(1)
}

// CreateMany mocks the CreateMany method
func (m *S3ClientMock) CreateMany(ctx context.Context, files []*s3.FileData) ([]string, error) {
	args := m.Called(ctx, files)
	return args.Get(0).([]string), args.Error(1)
}

// GetOne mocks the GetOne method
func (m *S3ClientMock) GetOne(ctx context.Context, objectID string) (io.Reader, error) {
	args := m.Called(ctx, objectID)
	return args.Get(0).(io.Reader), args.Error(1)
}

// GetMany mocks the GetMany method
func (m *S3ClientMock) GetMany(ctx context.Context, objectIDs []string) ([]io.Reader, error) {
	args := m.Called(ctx, objectIDs)
	return args.Get(0).([]io.Reader), args.Error(1)
}

// DeleteOne mocks the DeleteOne method
func (m *S3ClientMock) DeleteOne(ctx context.Context, objectID string) error {
	args := m.Called(ctx, objectID)
	return args.Error(0)
}

// DeleteMany mocks the DeleteMany method
func (m *S3ClientMock) DeleteMany(ctx context.Context, objectIDs []string) error {
	args := m.Called(ctx, objectIDs)
	return args.Error(0)
}
