package profile_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/domain/service"
	"github.com/mandarine-io/Backend/internal/domain/service/master/profile"
	"github.com/mandarine-io/Backend/internal/domain/service/master/profile/mapper"
	"github.com/mandarine-io/Backend/internal/helper/ref"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	mock2 "github.com/mandarine-io/Backend/internal/persistence/repo/mock"
	gormType "github.com/mandarine-io/Backend/internal/persistence/type"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	masterProfileRepo = new(mock2.MasterProfileRepositoryMock)
	svc               = profile.NewService(masterProfileRepo)
	ctx               = context.Background()
)

func Test_MasterProfileService_CreateMasterProfile(t *testing.T) {
	userId := uuid.New()

	t.Run("Failed to check if master profile exists", func(t *testing.T) {
		expectedErr := errors.New("failed to check if master profile exists")
		masterProfileRepo.On("ExistsMasterProfileByUserId", ctx, userId).Return(false, expectedErr).Once()

		_, err := svc.CreateMasterProfile(ctx, userId, dto.CreateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("User already has a master profile", func(t *testing.T) {
		masterProfileRepo.On("ExistsMasterProfileByUserId", ctx, userId).Return(true, nil).Once()

		_, err := svc.CreateMasterProfile(ctx, userId, dto.CreateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, service.ErrDuplicateMasterProfile, err)
	})

	t.Run("Duplicate master profile", func(t *testing.T) {
		masterProfileRepo.On("ExistsMasterProfileByUserId", ctx, userId).Return(false, nil).Once()
		masterProfileRepo.On("CreateMasterProfile", ctx, mock.Anything).Return(nil, repo.ErrDuplicateMasterProfile).Once()

		_, err := svc.CreateMasterProfile(ctx, userId, dto.CreateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, service.ErrDuplicateMasterProfile, err)
	})

	t.Run("User for master profile not exist", func(t *testing.T) {
		masterProfileRepo.On("ExistsMasterProfileByUserId", ctx, userId).Return(false, nil).Once()
		masterProfileRepo.On("CreateMasterProfile", ctx, mock.Anything).Return(nil, repo.ErrUserForMasterProfileNotExist).Once()

		_, err := svc.CreateMasterProfile(ctx, userId, dto.CreateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, service.ErrUserNotFound, err)
	})

	t.Run("Failed to create master profile", func(t *testing.T) {
		expectedErr := errors.New("failed to create master profile")
		masterProfileRepo.On("ExistsMasterProfileByUserId", ctx, userId).Return(false, nil).Once()
		masterProfileRepo.On("CreateMasterProfile", ctx, mock.Anything).Return(nil, expectedErr).Once()

		_, err := svc.CreateMasterProfile(ctx, userId, dto.CreateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Success", func(t *testing.T) {
		entity := &model.MasterProfileEntity{
			UserID:      userId,
			DisplayName: "test",
			Job:         "test",
			Point:       gormType.NewPoint(0, 0),
			Address:     ref.SafeRef("test"),
			Description: ref.SafeRef("test"),
			AvatarID:    ref.SafeRef("test"),
			IsEnabled:   true,
		}
		masterProfileRepo.On("ExistsMasterProfileByUserId", ctx, userId).Return(false, nil).Once()
		masterProfileRepo.On("CreateMasterProfile", ctx, mock.Anything).Return(entity, nil).Once()

		resp, err := svc.CreateMasterProfile(ctx, userId, dto.CreateMasterProfileInput{})

		assert.NoError(t, err)
		assert.Equal(t, entity.DisplayName, resp.DisplayName)
		assert.Equal(t, entity.Job, resp.Job)

		actualLat := resp.Point.Latitude
		assert.NoError(t, err)
		assert.Equal(t, entity.Point.Lat, actualLat)

		actualLng := resp.Point.Longitude
		assert.NoError(t, err)
		assert.Equal(t, entity.Point.Lng, actualLng)

		assert.Equal(t, entity.Address, resp.Address)
		assert.Equal(t, entity.Description, resp.Description)
		assert.Equal(t, entity.AvatarID, resp.AvatarID)
		assert.Equal(t, entity.IsEnabled, resp.IsEnabled)
	})
}

func Test_MasterProfileService_UpdateMasterProfile(t *testing.T) {
	userId := uuid.New()

	updateMasterProfileMatcherFactory := func(input dto.UpdateMasterProfileInput) interface{} {
		return mock.MatchedBy(func(p *model.MasterProfileEntity) bool {
			return p.DisplayName == input.DisplayName &&
				p.Job == input.Job &&
				p.Description == input.Description &&
				p.Address == input.Address &&
				p.Point.Lat == mapper.MapPointStringToPoint(input.Point).Lat &&
				p.Point.Lng == mapper.MapPointStringToPoint(input.Point).Lng &&
				p.AvatarID == input.AvatarID
		})
	}

	t.Run("Failed to find master profile by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to find master profile by user id")
		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(nil, expectedErr).Once()

		_, err := svc.UpdateMasterProfile(ctx, userId, dto.UpdateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Master profile not exist", func(t *testing.T) {
		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(nil, nil).Once()

		_, err := svc.UpdateMasterProfile(ctx, userId, dto.UpdateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, service.ErrMasterProfileNotExist, err)
	})

	t.Run("Failed to update master profile", func(t *testing.T) {
		expectedErr := errors.New("failed to update master profile")
		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(&model.MasterProfileEntity{}, nil).Once()
		masterProfileRepo.On("UpdateMasterProfile", ctx, mock.Anything).Return(nil, expectedErr).Once()

		_, err := svc.UpdateMasterProfile(ctx, userId, dto.UpdateMasterProfileInput{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Success", func(t *testing.T) {
		entity := &model.MasterProfileEntity{
			UserID:      userId,
			DisplayName: "test",
			Job:         "test",
			Point:       gormType.NewPoint(0, 0),
			Address:     ref.SafeRef("test"),
			Description: ref.SafeRef("test"),
			AvatarID:    ref.SafeRef("test"),
			IsEnabled:   true,
		}
		input := dto.UpdateMasterProfileInput{
			DisplayName: "test1",
			Job:         "test1",
			Point:       "1,1",
			Address:     ref.SafeRef("test1"),
			Description: ref.SafeRef("test1"),
			AvatarID:    ref.SafeRef("test1"),
		}
		updatedEntity := &model.MasterProfileEntity{
			UserID:      userId,
			DisplayName: "test1",
			Job:         "test1",
			Point:       gormType.NewPoint(1, 1),
			Address:     ref.SafeRef("test1"),
			Description: ref.SafeRef("test1"),
			AvatarID:    ref.SafeRef("test1"),
			IsEnabled:   true,
		}

		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(entity, nil).Once()
		masterProfileRepo.On("UpdateMasterProfile", ctx, updateMasterProfileMatcherFactory(input)).Return(updatedEntity, nil).Once()

		resp, err := svc.UpdateMasterProfile(ctx, userId, input)

		assert.NoError(t, err)
		assert.Equal(t, updatedEntity.DisplayName, resp.DisplayName)
		assert.Equal(t, updatedEntity.Job, resp.Job)

		actualLat := resp.Point.Latitude
		assert.NoError(t, err)
		assert.Equal(t, updatedEntity.Point.Lat, actualLat)

		actualLng := resp.Point.Longitude
		assert.NoError(t, err)
		assert.Equal(t, updatedEntity.Point.Lng, actualLng)

		assert.Equal(t, updatedEntity.Address, resp.Address)
		assert.Equal(t, updatedEntity.Description, resp.Description)
		assert.Equal(t, updatedEntity.AvatarID, resp.AvatarID)
	})
}

func Test_MasterProfileService_GetOwnMasterProfile(t *testing.T) {
	userId := uuid.New()

	t.Run("Failed to find master profile by user id", func(t *testing.T) {
		expectedErr := errors.New("failed to find master profile by user id")
		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(nil, expectedErr).Once()

		_, err := svc.GetOwnMasterProfile(ctx, userId)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Master profile not exist", func(t *testing.T) {
		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(nil, nil).Once()

		_, err := svc.GetOwnMasterProfile(ctx, userId)

		assert.Error(t, err)
		assert.Equal(t, service.ErrMasterProfileNotExist, err)
	})

	t.Run("Success", func(t *testing.T) {
		entity := &model.MasterProfileEntity{
			UserID:      userId,
			DisplayName: "test",
			Job:         "test",
			Point:       gormType.NewPoint(0, 0),
			Address:     ref.SafeRef("test"),
			Description: ref.SafeRef("test"),
			AvatarID:    ref.SafeRef("test"),
			IsEnabled:   true,
		}
		masterProfileRepo.On("FindMasterProfileByUserId", ctx, userId).Return(entity, nil).Once()

		resp, err := svc.GetOwnMasterProfile(ctx, userId)

		assert.NoError(t, err)
		assert.Equal(t, entity.DisplayName, resp.DisplayName)
		assert.Equal(t, entity.Job, resp.Job)

		actualLat := resp.Point.Latitude
		assert.NoError(t, err)
		assert.Equal(t, entity.Point.Lat, actualLat)

		actualLng := resp.Point.Longitude
		assert.NoError(t, err)
		assert.Equal(t, entity.Point.Lng, actualLng)

		assert.Equal(t, entity.Address, resp.Address)
		assert.Equal(t, entity.Description, resp.Description)
		assert.Equal(t, entity.AvatarID, resp.AvatarID)
		assert.Equal(t, entity.IsEnabled, resp.IsEnabled)
	})
}

func Test_MasterProfileService_GetMasterProfileByUsername(t *testing.T) {
	username := "test"

	t.Run("Failed to find master profile by username", func(t *testing.T) {
		expectedErr := errors.New("failed to find master profile by username")
		masterProfileRepo.On("FindEnabledMasterProfileByUsername", ctx, username).Return(nil, expectedErr).Once()

		_, err := svc.GetMasterProfileByUsername(ctx, username)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Master profile not exist", func(t *testing.T) {
		masterProfileRepo.On("FindEnabledMasterProfileByUsername", ctx, username).Return(nil, nil).Once()

		_, err := svc.GetMasterProfileByUsername(ctx, username)

		assert.Error(t, err)
		assert.Equal(t, service.ErrMasterProfileNotFound, err)
	})

	t.Run("Success", func(t *testing.T) {
		entity := &model.MasterProfileEntity{
			UserID:      uuid.New(),
			DisplayName: "test",
			Job:         "test",
			Point:       gormType.NewPoint(0, 0),
			Address:     ref.SafeRef("test"),
			Description: ref.SafeRef("test"),
			AvatarID:    ref.SafeRef("test"),
			IsEnabled:   true,
		}
		masterProfileRepo.On("FindEnabledMasterProfileByUsername", ctx, username).Return(entity, nil).Once()

		resp, err := svc.GetMasterProfileByUsername(ctx, username)

		assert.NoError(t, err)
		assert.Equal(t, entity.DisplayName, resp.DisplayName)
		assert.Equal(t, entity.Job, resp.Job)

		actualLat := resp.Point.Latitude
		assert.NoError(t, err)
		assert.Equal(t, entity.Point.Lat, actualLat)

		actualLng := resp.Point.Longitude
		assert.NoError(t, err)
		assert.Equal(t, entity.Point.Lng, actualLng)

		assert.Equal(t, entity.Address, resp.Address)
		assert.Equal(t, entity.Description, resp.Description)
		assert.Equal(t, entity.AvatarID, resp.AvatarID)
	})
}

func Test_MasterProfileService_FindMasterProfiles(t *testing.T) {
	t.Run("Unavailable sort field", func(t *testing.T) {
		pagination := dto.PaginationInput{
			Page:     0,
			PageSize: 10,
		}
		sorts := dto.SortInput{
			Field: "unavailable",
			Order: "asc",
		}
		filter := dto.FindMasterProfilesFilterInput{
			DisplayName: ref.SafeRef("test"),
		}
		input := dto.FindMasterProfilesInput{
			FindMasterProfilesFilterInput: &filter,
			SortInput:                     &sorts,
			PaginationInput:               &pagination,
		}

		_, err := svc.FindMasterProfiles(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, service.ErrUnavailableSortField, err)
	})

	t.Run("Failed to find master profiles", func(t *testing.T) {
		pagination := dto.PaginationInput{
			Page:     0,
			PageSize: 10,
		}
		input := dto.FindMasterProfilesInput{
			PaginationInput: &pagination,
		}

		expectedErr := errors.New("failed to find master profiles")
		masterProfileRepo.On("FindMasterProfiles", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedErr).Once()
		masterProfileRepo.On("CountMasterProfiles", ctx, mock.Anything).Return(int64(0), nil).Once()

		_, err := svc.FindMasterProfiles(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Failed to count master profiles", func(t *testing.T) {
		pagination := dto.PaginationInput{
			Page:     0,
			PageSize: 10,
		}
		input := dto.FindMasterProfilesInput{
			PaginationInput: &pagination,
		}

		expectedErr := errors.New("failed to count master profiles")
		masterProfileRepo.On("FindMasterProfiles", ctx, mock.Anything, mock.Anything, mock.Anything).Return([]*model.MasterProfileEntity{}, nil).Once()
		masterProfileRepo.On("CountMasterProfiles", ctx, mock.Anything).Return(int64(0), expectedErr).Once()

		_, err := svc.FindMasterProfiles(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Success", func(t *testing.T) {
		pagination := dto.PaginationInput{
			Page:     0,
			PageSize: 10,
		}
		input := dto.FindMasterProfilesInput{
			PaginationInput: &pagination,
		}

		entities := []*model.MasterProfileEntity{
			{
				UserID:      uuid.New(),
				DisplayName: "test",
				Job:         "test",
				Point:       gormType.NewPoint(0, 0),
				Address:     ref.SafeRef("test"),
				Description: ref.SafeRef("test"),
				AvatarID:    ref.SafeRef("test"),
				IsEnabled:   true,
			},
		}

		masterProfileRepo.On("FindMasterProfiles", ctx, mock.Anything, mock.Anything, mock.Anything).Return(entities, nil).Once()
		masterProfileRepo.On("CountMasterProfiles", ctx, mock.Anything).Return(int64(1), nil).Once()

		resp, err := svc.FindMasterProfiles(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, resp.Count, 1)
		assert.Equal(t, entities[0].DisplayName, resp.Data[0].DisplayName)
		assert.Equal(t, entities[0].Job, resp.Data[0].Job)

		actualLat := resp.Data[0].Point.Latitude
		assert.NoError(t, err)
		assert.Equal(t, entities[0].Point.Lat, actualLat)

		actualLng := resp.Data[0].Point.Longitude
		assert.NoError(t, err)
		assert.Equal(t, entities[0].Point.Lng, actualLng)

		assert.Equal(t, entities[0].Address, resp.Data[0].Address)
		assert.Equal(t, entities[0].Description, resp.Data[0].Description)
		assert.Equal(t, entities[0].AvatarID, resp.Data[0].AvatarID)
	})
}
