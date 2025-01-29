package masterprofile

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/converter"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	gormType "github.com/mandarine-io/backend/internal/persistence/types"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
)

var (
	updateMasterProfileMatcherFactory = func(input v0.UpdateMasterProfileInput) interface{} {
		return mock.MatchedBy(
			func(p *entity.MasterProfile) bool {
				return p.DisplayName == input.DisplayName &&
					p.Job == input.Job &&
					p.Description == input.Description &&
					p.Address == input.Address &&
					p.Point.Lat == converter.MapPointStringToPoint(input.Point).Lat &&
					p.Point.Lng == converter.MapPointStringToPoint(input.Point).Lng &&
					p.AvatarID == input.AvatarID
			},
		)
	}
)

type UpdateMasterProfileSuite struct {
	suite.Suite
}

func (s *UpdateMasterProfileSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Master profile service")
	t.Feature("UpdateMasterProfile")
	t.Tags("Positive")

	userID := uuid.New()
	e := &entity.MasterProfile{
		UserID:      userID,
		DisplayName: "test",
		Job:         "test",
		Point:       *gormType.NewPoint(0, 0),
		Address:     lo.ToPtr("test"),
		Description: lo.ToPtr("test"),
		AvatarID:    lo.ToPtr("test"),
		IsEnabled:   true,
	}
	input := v0.UpdateMasterProfileInput{
		DisplayName: "test1",
		Job:         "test1",
		Point:       "1,1",
		Address:     lo.ToPtr("test1"),
		Description: lo.ToPtr("test1"),
		AvatarID:    lo.ToPtr("test1"),
	}
	updatedEntity := &entity.MasterProfile{
		UserID:      userID,
		DisplayName: "test1",
		Job:         "test1",
		Point:       *gormType.NewPoint(1, 1),
		Address:     lo.ToPtr("test1"),
		Description: lo.ToPtr("test1"),
		AvatarID:    lo.ToPtr("test1"),
		IsEnabled:   true,
	}

	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(e, nil).Once()
	masterProfileRepoMock.On("UpdateMasterProfile", ctx, updateMasterProfileMatcherFactory(input)).Return(
		updatedEntity,
		nil,
	).Once()

	resp, err := svc.UpdateMasterProfile(ctx, userID, input)

	t.Require().NoError(err)
	t.Require().Equal(updatedEntity.DisplayName, resp.DisplayName)
	t.Require().Equal(updatedEntity.Job, resp.Job)

	actualLat := resp.Point.Latitude
	t.Require().NoError(err)
	t.Require().Equal(updatedEntity.Point.Lat, actualLat)

	actualLng := resp.Point.Longitude
	t.Require().NoError(err)
	t.Require().Equal(updatedEntity.Point.Lng, actualLng)

	t.Require().Equal(updatedEntity.Address, resp.Address)
	t.Require().Equal(updatedEntity.Description, resp.Description)
	t.Require().Equal(updatedEntity.AvatarID, resp.AvatarID)
}

func (s *UpdateMasterProfileSuite) Test_ErrFindMasterProfileByUserID(t provider.T) {
	t.Title("Returns finding master profile by user ID error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("UpdateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	expectedErr := errors.New("failed to find master profile by user id")
	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(nil, expectedErr).Once()

	_, err := svc.UpdateMasterProfile(ctx, userID, v0.UpdateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *UpdateMasterProfileSuite) Test_ErrMasterProfileNotFound(t provider.T) {
	t.Title("Returns master profile not found error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("UpdateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(nil, nil).Once()

	_, err := svc.UpdateMasterProfile(ctx, userID, v0.UpdateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(domain.ErrMasterProfileNotExist, err)
}

func (s *UpdateMasterProfileSuite) Test_ErrUpdateMasterProfile(t provider.T) {
	t.Title("Returns updating master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("UpdateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	expectedErr := errors.New("failed to update master profile")
	masterProfileRepoMock.On("FindMasterProfileByUserID", ctx, userID).Return(&entity.MasterProfile{}, nil).Once()
	masterProfileRepoMock.On("UpdateMasterProfile", ctx, mock.Anything).Return(nil, expectedErr).Once()

	_, err := svc.UpdateMasterProfile(ctx, userID, v0.UpdateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
