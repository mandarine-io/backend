package resource

import (
	"context"
	mock1 "github.com/mandarine-io/backend/internal/infrastructure/s3/mock"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/internal/service/domain/resource"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	ctx = context.Background()

	s3ManagerMock *mock1.ManagerMock
	svc           domain.ResourceService
)

func init() {
	s3ManagerMock = new(mock1.ManagerMock)
	svc = resource.NewService(s3ManagerMock)
}

type ResourceServiceSuite struct {
	suite.Suite
}

func TestResourceServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(ResourceServiceSuite))
}

func (s *ResourceServiceSuite) Test(t provider.T) {
	s.RunSuite(t, new(DownloadResourceSuite))
	s.RunSuite(t, new(UploadResourceSuite))
	s.RunSuite(t, new(UploadResourcesSuite))
}
