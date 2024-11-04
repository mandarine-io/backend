package mock

import (
	"context"
	"github.com/mandarine-io/Backend/pkg/storage/s3"
	"github.com/stretchr/testify/mock"
)

// S3ClientMock is a mock for the S3Client interface
type S3ClientMock struct {
	mock.Mock
}

// CreateOne mocks the CreateOne method
func (m *S3ClientMock) CreateOne(ctx context.Context, file *s3.FileData) *s3.CreateDto {
	args := m.Called(ctx, file)
	return args.Get(0).(*s3.CreateDto)
}

// CreateMany mocks the CreateMany method
func (m *S3ClientMock) CreateMany(ctx context.Context, files []*s3.FileData) map[string]*s3.CreateDto {
	args := m.Called(ctx, files)
	return args.Get(0).(map[string]*s3.CreateDto)
}

// GetOne mocks the GetOne method
func (m *S3ClientMock) GetOne(ctx context.Context, objectID string) *s3.GetDto {
	args := m.Called(ctx, objectID)
	return args.Get(0).(*s3.GetDto)
}

// GetMany mocks the GetMany method
func (m *S3ClientMock) GetMany(ctx context.Context, objectIDs []string) map[string]*s3.GetDto {
	args := m.Called(ctx, objectIDs)
	return args.Get(0).(map[string]*s3.GetDto)
}

// DeleteOne mocks the DeleteOne method
func (m *S3ClientMock) DeleteOne(ctx context.Context, objectID string) error {
	args := m.Called(ctx, objectID)
	return args.Error(0)
}

// DeleteMany mocks the DeleteMany method
func (m *S3ClientMock) DeleteMany(ctx context.Context, objectIDs []string) map[string]error {
	args := m.Called(ctx, objectIDs)
	return args.Get(0).(map[string]error)
}
