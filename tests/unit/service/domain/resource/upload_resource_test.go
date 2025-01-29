package resource

import (
	"errors"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"os"
)

type UploadResourceSuite struct {
	suite.Suite
}

func (s *UploadResourceSuite) Test_FileIsNil(t provider.T) {
	t.Title("Returns FileIsNil")
	t.Severity(allure.CRITICAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Negative")

	_, err := svc.UploadResource(ctx, &v0.UploadResourceInput{Resource: nil})

	t.Require().Error(err)
	t.Require().Equal(domain.ErrResourceNotUploaded, err)
}

func (s *UploadResourceSuite) Test_ErrUploadFile(t provider.T) {
	t.Title("Returns ErrUploadFile")
	t.Severity(allure.CRITICAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Negative")

	// Create temp file
	file, err := os.CreateTemp("", "test-")
	t.Require().Nil(err)
	defer func(name string) {
		err := os.Remove(name)
		t.Require().NoError(err)
	}(file.Name())

	// Create FileHeader
	fileHeader := createMultipartFileHeader(file.Name())
	t.Require().NotNil(fileHeader)

	// Mock
	expectedErr := errors.New("s3 error")
	s3ManagerMock.On("CreateOne", ctx, mock.Anything).Return(s3.CreateResult{Error: expectedErr}).Once()

	// Call service
	_, err = svc.UploadResource(ctx, &v0.UploadResourceInput{Resource: fileHeader})

	t.Require().Error(err)
	t.Require().Equal(err, expectedErr)
}

func (s *UploadResourceSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Positive")

	// Create temp file
	file, err := os.CreateTemp("", "test-")
	t.Require().Nil(err)
	defer func(name string) {
		err := os.Remove(name)
		t.Require().NoError(err)
	}(file.Name())

	// Create FileHeader
	fileHeader := createMultipartFileHeader(file.Name())
	t.Require().NotNil(fileHeader)

	// Mock
	s3ManagerMock.On("CreateOne", ctx, mock.Anything).Return(s3.CreateResult{ObjectID: "test"}).Once()

	// Call service
	output, err := svc.UploadResource(ctx, &v0.UploadResourceInput{Resource: fileHeader})

	t.Require().Nil(err)
	t.Require().Equal(output.ObjectID, "test")
}
