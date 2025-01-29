package masterprofile

import (
	"context"
	"github.com/mandarine-io/backend/internal/persistence/repo/mock"
	"github.com/mandarine-io/backend/internal/service/domain"
	masterprofile "github.com/mandarine-io/backend/internal/service/domain/master/profile"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var (
	ctx = context.Background()

	masterProfileRepoMock *mock.MasterProfileRepositoryMock
	svc                   domain.MasterProfileService
)

func init() {
	masterProfileRepoMock = new(mock.MasterProfileRepositoryMock)
	svc = masterprofile.NewService(masterProfileRepoMock)
}

type MasterProfileServiceSuite struct {
	suite.Suite
}

func TestMasterProfileServiceSuite(t *testing.T) {
	suite.RunSuite(t, new(MasterProfileServiceSuite))
}

func (s *MasterProfileServiceSuite) Test(t provider.T) {
	s.RunSuite(t, new(CreateMasterProfileSuite))
	s.RunSuite(t, new(FindMasterProfilesSuite))
	s.RunSuite(t, new(GetMasterProfileByUsernameSuite))
	s.RunSuite(t, new(GetOwnMasterProfileSuite))
	s.RunSuite(t, new(UpdateMasterProfileSuite))
}
