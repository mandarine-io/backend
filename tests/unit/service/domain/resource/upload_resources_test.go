package resource

import (
	"errors"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"os"
)

type UploadResourcesSuite struct {
	suite.Suite
}

func (s *UploadResourcesSuite) Test_FileIsNil(t provider.T) {
	t.Title("Returns FileIsNil")
	t.Severity(allure.CRITICAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Negative")

	fileHeaders := make([]*multipart.FileHeader, 3)

	output, err := svc.UploadResources(ctx, &v0.UploadResourcesInput{Resources: fileHeaders})

	t.Require().NoError(err)
	t.Require().Equal(0, output.Count)
	t.Require().Equal(map[string]v0.UploadResourceOutput{}, output.Data)
}

func (s *UploadResourcesSuite) Test_ErrUploadFiles(t provider.T) {
	t.Title("Returns FileIsNil")
	t.Severity(allure.CRITICAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Negative")

	var err error
	files := make([]*os.File, 3)
	fileHeaders := make([]*multipart.FileHeader, 3)
	for i := 0; i < 3; i++ {
		// Create temp file
		files[i], err = os.CreateTemp("", fmt.Sprintf("test-%d-", i))
		t.Require().Nil(err)

		// Create FileHeader
		fileHeaders[i] = createMultipartFileHeader(files[i].Name())
		t.Require().NotNil(fileHeaders[i])
	}
	defer func() {
		for _, file := range files {
			err = os.Remove(file.Name())
			t.Require().Nil(err)
		}
	}()

	t.Run(
		"All files error", func(t provider.T) {
			// Mock
			s3Output := make(map[string]s3.CreateResult)
			for i := 0; i < 3; i++ {
				s3Output[files[i].Name()] = s3.CreateResult{Error: errors.New("s3 error")}
			}
			s3ManagerMock.On("CreateMany", ctx, mock.Anything).Return(s3Output).Once()

			output, err := svc.UploadResources(ctx, &v0.UploadResourcesInput{Resources: fileHeaders})

			t.Require().NoError(err)
			t.Require().Equal(0, output.Count)
			t.Require().Equal(map[string]v0.UploadResourceOutput{}, output.Data)
		},
	)

	t.Run(
		"Have one files error", func(t provider.T) {
			// Mock
			s3Output := make(map[string]s3.CreateResult)
			for i := 0; i < 3; i++ {
				if i == 0 {
					s3Output[files[i].Name()] = s3.CreateResult{ObjectID: "test", Error: nil}
					continue
				}
				s3Output[files[i].Name()] = s3.CreateResult{Error: errors.New("s3 error")}
			}
			s3ManagerMock.On("CreateMany", ctx, mock.Anything).Return(s3Output).Once()

			output, err := svc.UploadResources(ctx, &v0.UploadResourcesInput{Resources: fileHeaders})

			t.Require().NoError(err)
			t.Require().Equal(1, output.Count)

			resourceOutput, ok := output.Data[files[0].Name()]
			t.Require().True(ok)
			t.Require().Equal("test", resourceOutput.ObjectID)
		},
	)

	t.Run(
		"No files error", func(t provider.T) {
			// Mock
			s3Output := make(map[string]s3.CreateResult)
			for i := 0; i < 3; i++ {
				s3Output[files[i].Name()] = s3.CreateResult{ObjectID: fmt.Sprintf("test-%d", i), Error: nil}
			}
			s3ManagerMock.On("CreateMany", ctx, mock.Anything).Return(s3Output).Once()

			output, err := svc.UploadResources(ctx, &v0.UploadResourcesInput{Resources: fileHeaders})

			t.Require().NoError(err)
			t.Require().Equal(3, output.Count)

			for i := 0; i < 3; i++ {
				resourceOutput, ok := output.Data[files[i].Name()]
				t.Require().True(ok)
				t.Require().Equal(fmt.Sprintf("test-%d", i), resourceOutput.ObjectID)
			}
		},
	)
}
