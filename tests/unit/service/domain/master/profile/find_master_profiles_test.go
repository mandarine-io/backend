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
	"gorm.io/gorm"
)

type FindMasterProfilesSuite struct {
	suite.Suite
}

func (s *FindMasterProfilesSuite) Test_Success(t provider.T) {
	t.Title("Returns success")
	t.Severity(allure.NORMAL)
	t.Epic("Master profile service")
	t.Feature("FindMasterProfiles")
	t.Tags("Positive")

	pagination := v0.PaginationInput{
		Page:     0,
		PageSize: 10,
	}
	filter := v0.FindMasterProfilesFilterInput{
		DisplayName: lo.ToPtr("test"),
		Job:         lo.ToPtr("test"),
		Lng:         lo.ToPtr(decimal.NewFromFloat(0)),
		Lat:         lo.ToPtr(decimal.NewFromFloat(0)),
		Radius:      lo.ToPtr(decimal.NewFromFloat(0)),
	}
	sorts := v0.SortInput{
		Field: "display_name",
		Order: "asc",
	}
	input := v0.FindMasterProfilesInput{
		PaginationInput:               &pagination,
		FindMasterProfilesFilterInput: &filter,
		SortInput:                     &sorts,
	}

	entities := []*entity.MasterProfile{
		{
			UserID:      uuid.New(),
			DisplayName: "test",
			Job:         "test",
			Point:       *gormType.NewPoint(decimal.NewFromFloat(0), decimal.NewFromFloat(0)),
			Address:     lo.ToPtr("test"),
			Description: lo.ToPtr("test"),
			AvatarID:    lo.ToPtr("test"),
			IsEnabled:   true,
		},
	}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	masterProfileRepoMock.On("WithPagination", mock.Anything, mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("WithColumnSort", mock.Anything, mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("WithPointSort", mock.Anything, mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("WithDisplayNameFilter", mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("WithJobFilter", mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("WithAddressFilter", mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("WithPointFilter", mock.Anything, mock.Anything, mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On(
		"FindMasterProfiles",
		ctx,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).
		Return(entities, nil).Once()
	masterProfileRepoMock.On(
		"CountMasterProfiles",
		ctx,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).
		Return(int64(1), nil).Once()

	resp, err := svc.FindMasterProfiles(ctx, input)

	t.Require().NoError(err)
	t.Require().Equal(resp.Count, 1)
	t.Require().Equal(entities[0].DisplayName, resp.Data[0].DisplayName)
	t.Require().Equal(entities[0].Job, resp.Data[0].Job)

	actualLat := resp.Data[0].Point.Latitude
	t.Require().NoError(err)
	t.Require().Equal(entities[0].Point.Lat, actualLat)

	actualLng := resp.Data[0].Point.Longitude
	t.Require().NoError(err)
	t.Require().Equal(entities[0].Point.Lng, actualLng)

	t.Require().Equal(entities[0].Address, resp.Data[0].Address)
	t.Require().Equal(entities[0].Description, resp.Data[0].Description)
	t.Require().Equal(entities[0].AvatarID, resp.Data[0].AvatarID)
}

func (s *FindMasterProfilesSuite) Test_ErrUnavailableSortField(t provider.T) {
	t.Title("Returns unavailable sort field error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("FindMasterProfiles")
	t.Tags("Negative")

	sorts := v0.SortInput{
		Field: "unavailable",
		Order: "asc",
	}
	input := v0.FindMasterProfilesInput{
		SortInput: &sorts,
	}

	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }
	masterProfileRepoMock.On("WithPagination", mock.Anything, mock.Anything).Once().Return(scope)

	_, err := svc.FindMasterProfiles(ctx, input)

	t.Require().Error(err)
	t.Require().Equal(domain.ErrUnavailableSortField, err)
}

func (s *FindMasterProfilesSuite) Test_ErrFindMasterProfile(t provider.T) {
	t.Title("Returns finding master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("FindMasterProfiles")
	t.Tags("Negative")

	pagination := v0.PaginationInput{
		Page:     0,
		PageSize: 10,
	}
	input := v0.FindMasterProfilesInput{
		PaginationInput: &pagination,
	}

	expectedErr := errors.New("failed to find master profiles")
	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }

	masterProfileRepoMock.On("WithPagination", mock.Anything, mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On("FindMasterProfiles", ctx, mock.Anything, mock.Anything, mock.Anything).Return(
		nil,
		expectedErr,
	).Once()
	masterProfileRepoMock.On("CountMasterProfiles", ctx, mock.Anything).Return(int64(0), nil).Once()

	_, err := svc.FindMasterProfiles(ctx, input)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}

func (s *FindMasterProfilesSuite) Test_ErrCountMasterProfile(t provider.T) {
	t.Title("Returns counting master profile error")
	t.Severity(allure.CRITICAL)
	t.Epic("Master profile service")
	t.Feature("FindMasterProfiles")
	t.Tags("Negative")

	pagination := v0.PaginationInput{
		Page:     0,
		PageSize: 10,
	}
	input := v0.FindMasterProfilesInput{
		PaginationInput: &pagination,
	}

	expectedErr := errors.New("failed to count master profiles")
	var scope repo.Scope = func(db *gorm.DB) *gorm.DB { return db }

	masterProfileRepoMock.On("WithPagination", mock.Anything, mock.Anything).Once().Return(scope)
	masterProfileRepoMock.On(
		"FindMasterProfiles",
		ctx,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return([]*entity.MasterProfile{}, nil).Once()
	masterProfileRepoMock.On("CountMasterProfiles", ctx, mock.Anything).Return(int64(0), expectedErr).Once()

	_, err := svc.FindMasterProfiles(ctx, input)

	t.Require().Error(err)
	t.Require().Equal(expectedErr, err)
}
