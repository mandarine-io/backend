package minio

import (
	"context"
	"github.com/mandarine-io/backend/internal/infrastructure/s3"
	minio2 "github.com/mandarine-io/backend/internal/infrastructure/s3/minio"
	"github.com/mandarine-io/backend/tests/integration"
	"github.com/minio/minio-go/v7"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	ctx     = context.Background()
	cfg     minio2.Config
	client  *minio.Client
	manager s3.Manager
)

type MinioS3ManagerSuite struct {
	suite.Suite
}

func TestMinioS3ManagerSuite(t *testing.T) {
	cfg = integration.Cfg.GetMinioConfig()

	var err error
	client, err = minio2.NewClient(cfg)
	require.NoError(t, err)

	manager, err = minio2.NewManager(client, cfg.BucketName)
	require.NoError(t, err)

	suite.RunSuite(t, new(MinioS3ManagerSuite))
}

func (s *MinioS3ManagerSuite) Test(t provider.T) {
	s.RunSuite(t, new(CreateOneSuite))
	s.RunSuite(t, new(CreateManySuite))
	s.RunSuite(t, new(DeleteOneSuite))
	s.RunSuite(t, new(DeleteManySuite))
	s.RunSuite(t, new(GetOneSuite))
	s.RunSuite(t, new(GetManySuite))
}
