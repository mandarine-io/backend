package resource

import (
	"errors"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type DownloadResourceSuite struct {
	suite.Suite
}

func (s *DownloadResourceSuite) Test_NotFound(t provider.T) {
	t.Title("Returns NotFound")
	t.Severity(allure.CRITICAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Negative")

	expectError := errors.New("s3 error")
	s3ManagerMock.On("GetOne", ctx, "test").Return(s3.GetResult{Data: nil, Error: expectError}).Once()

	output, err := svc.DownloadResource(ctx, "test")

	t.Require().Error(err)
	t.Require().Equal(expectError, err)
	t.Require().Nil(output)
}

func (s *DownloadResourceSuite) Test_Success(t provider.T) {
	t.Title("Returns Success")
	t.Severity(allure.NORMAL)
	t.Epic("Resource service")
	t.Feature("Resource")
	t.Tags("Positive")

	s3ManagerMock.On("GetOne", ctx, "test").Return(s3.GetResult{Data: &s3.FileData{}, Error: nil}).Once()

	output, err := svc.DownloadResource(ctx, "test")

	t.Require().NoError(err)
	t.Require().NotNil(output)
	t.Require().Equal(&s3.FileData{}, output)
}
