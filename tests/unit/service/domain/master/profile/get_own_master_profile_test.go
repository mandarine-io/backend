package masterprofile

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	gormType "github.com/mandarine-io/backend/internal/persistence/types"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type GetOwnMasterProfileSuite struct {
	suite.Suite
}

func (s *GetOwnMasterProfileSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Master profile service")
	t.Feature("GetOwnMasterProfile")
	t.Tags("Positive")

	userID := uuid.New()
	e := &entity.MasterProfile{
		UserID:      userID,
		DisplayName: "test",
		Job:         "test",
		Point:       *gormType.NewPoint(decimal.NewFromFloat(0), decimal.NewFromFloat(0)),
		Address:     lo.ToPtr("test"),
		Description: lo.ToPtr("test"),
		AvatarID:    lo.ToPtr("test"),
		IsEnabled:   true,
	}
	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(e, nil).Once()

	resp, err := svc.GetOwnMasterProfile(ctx, userID)

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
	t.Require().NotNil(resp.IsEnabled)
	t.Require().Equal(e.IsEnabled, *resp.IsEnabled)
}

func (s *GetOwnMasterProfileSuite) Test_ErrFindMasterProfileByUserID(t provider.T) {
	t.Title("Returns finding master profile by user ID error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("GetOwnMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	expectedErr := errors.New("failed to find master profile by user id")
	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(nil, expectedErr).Once()

	_, err := svc.GetOwnMasterProfile(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *GetOwnMasterProfileSuite) Test_ErrMasterProfileNotExists(t provider.T) {
	t.Title("Returns master profile not exists error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("GetOwnMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	expectedErr := errors.New("failed to find master profile by user id")
	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(nil, expectedErr).Once()

	_, err := svc.GetOwnMasterProfile(ctx, userID)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
