package masterprofile

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	gormType "github.com/mandarine-io/backend/internal/persistence/types"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/samber/lo"
)

type GetMasterProfileByUsernameSuite struct {
	suite.Suite
}

func (s *GetMasterProfileByUsernameSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Master profile service")
	t.Feature("GetMasterProfileByUsername")
	t.Tags("Positive")

	username := "test"
	e := &entity.MasterProfile{
		UserID:      uuid.New(),
		DisplayName: "test",
		Job:         "test",
		Point:       *gormType.NewPoint(0, 0),
		Address:     lo.ToPtr("test"),
		Description: lo.ToPtr("test"),
		AvatarID:    lo.ToPtr("test"),
		IsEnabled:   true,
	}
	masterProfileRepoMock.On("FindEnabledMasterProfileByUsername", ctx, username).Return(e, nil).Once()

	resp, err := svc.GetMasterProfileByUsername(ctx, username)

	t.Require().NoError(err)
	t.Require().Equal(e.DisplayName, resp.DisplayName)
	t.Require().Equal(e.Job, resp.Job)

	actualLat := resp.Point.Latitude
	t.Require().NoError(err)
	t.Require().Equal(e.Point.Lat, actualLat)

	actualLng := resp.Point.Longitude
	t.Require().NoError(err)
	t.Require().Equal(e.Point.Lng, actualLng)

	t.Require().Equal(e.Address, resp.Address)
	t.Require().Equal(e.Description, resp.Description)
	t.Require().Equal(e.AvatarID, resp.AvatarID)
	t.Require().Nil(resp.IsEnabled)
}

func (s *GetMasterProfileByUsernameSuite) Test_ErrFindMasterProfile(t provider.T) {
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("GetMasterProfileByUsername")
	t.Tags("Negative")

	username := "test"
	expectedErr := errors.New("failed to find master profile by username")
	masterProfileRepoMock.On("FindEnabledMasterProfileByUsername", ctx, username).Return(nil, expectedErr).Once()

	_, err := svc.GetMasterProfileByUsername(ctx, username)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *GetMasterProfileByUsernameSuite) Test_ErrMasterProfileNotExists(t provider.T) {
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("GetMasterProfileByUsername")
	t.Tags("Negative")

	username := "test"
	masterProfileRepoMock.On("FindEnabledMasterProfileByUsername", ctx, username).Return(nil, nil).Once()

	_, err := svc.GetMasterProfileByUsername(ctx, username)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrMasterProfileNotFound, err)
}
