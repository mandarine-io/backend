package masterprofile

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	gormType "github.com/mandarine-io/backend/internal/persistence/types"
	"github.com/mandarine-io/backend/internal/service/domain"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type CreateMasterProfileSuite struct {
	suite.Suite
}

func (s *CreateMasterProfileSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Master profile service")
	t.Feature("CreateMasterProfile")
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
	masterProfileRepoMock.On("ExistsMasterProfileByUserID", ctx, userID).Return(false, nil).Once()
	masterProfileRepoMock.On("CreateMasterProfile", ctx, mock.Anything).Return(e, nil).Once()

	resp, err := svc.CreateMasterProfile(ctx, userID, v0.CreateMasterProfileInput{})

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

func (s *CreateMasterProfileSuite) Test_ErrExistsMasterProfile(t provider.T) {
	t.Title("Returns exists master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("CreateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	expectedErr := errors.New("failed to check if master profile exists")
	masterProfileRepoMock.On("ExistsMasterProfileByUserID", ctx, userID).Return(false, expectedErr).Once()

	_, err := svc.CreateMasterProfile(ctx, userID, v0.CreateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *CreateMasterProfileSuite) Test_ErrUserAlreadyHasMasterProfile(t provider.T) {
	t.Title("Returns user already has master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("CreateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	masterProfileRepoMock.On("ExistsMasterProfileByUserID", ctx, userID).Return(true, nil).Once()

	_, err := svc.CreateMasterProfile(ctx, userID, v0.CreateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(domain.ErrDuplicateMasterProfile, err)
}

func (s *CreateMasterProfileSuite) Test_ErrDuplicateMasterProfile(t provider.T) {
	t.Title("Returns duplicate master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("CreateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	masterProfileRepoMock.On("ExistsMasterProfileByUserID", ctx, userID).Return(false, nil).Once()
	masterProfileRepoMock.On("CreateMasterProfile", ctx, mock.Anything).Return(
		nil,
		repo.ErrDuplicateMasterProfile,
	).Once()

	_, err := svc.CreateMasterProfile(ctx, userID, v0.CreateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(domain.ErrDuplicateMasterProfile, err)
}

func (s *CreateMasterProfileSuite) Test_ErrUserNotFound(t provider.T) {
	t.Title("Returns user not found error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("CreateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	masterProfileRepoMock.On("ExistsMasterProfileByUserID", ctx, userID).Return(false, nil).Once()
	masterProfileRepoMock.On("CreateMasterProfile", ctx, mock.Anything).Return(
		nil,
		repo.ErrUserForMasterProfileNotExist,
	).Once()

	_, err := svc.CreateMasterProfile(ctx, userID, v0.CreateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUserNotFound, err)
}

func (s *CreateMasterProfileSuite) Test_ErrCreateMasterProfile(t provider.T) {
	t.Title("Returns create master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("CreateMasterProfile")
	t.Tags("Negative")

	userID := uuid.New()
	expectedErr := errors.New("failed to create master profile")
	masterProfileRepoMock.On("ExistsMasterProfileByUserID", ctx, userID).Return(false, nil).Once()
	masterProfileRepoMock.On("CreateMasterProfile", ctx, mock.Anything).Return(nil, expectedErr).Once()

	_, err := svc.CreateMasterProfile(ctx, userID, v0.CreateMasterProfileInput{})

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
