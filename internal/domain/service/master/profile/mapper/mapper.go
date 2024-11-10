package mapper

import (
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	gormType "github.com/mandarine-io/Backend/internal/persistence/type"
	"strconv"
	"strings"
)

func MapPointStringToPoint(point string) gormType.Point {
	var (
		latitude  float64
		longitude float64
	)

	parts := strings.Split(point, ",")
	if len(parts) >= 2 {
		longitude, _ = strconv.ParseFloat(parts[0], 64)
		latitude, _ = strconv.ParseFloat(parts[1], 64)
	}

	return gormType.NewPoint(latitude, longitude)
}

func MapPointToDomainPoint(point gormType.Point) dto.PointOutput {
	return dto.PointOutput{
		Latitude:  point.Lat,
		Longitude: point.Lng,
	}
}

func MapCreateMasterProfileInputToEntity(userId uuid.UUID, input dto.CreateMasterProfileInput) *model.MasterProfileEntity {
	return &model.MasterProfileEntity{
		UserID:      userId,
		DisplayName: input.DisplayName,
		Job:         input.Job,
		Description: input.Description,
		Address:     input.Address,
		Point:       MapPointStringToPoint(input.Point),
		AvatarID:    input.AvatarID,
		IsEnabled:   true,
	}
}

func MapUpdateMasterProfileInputToEntity(entity *model.MasterProfileEntity, input dto.UpdateMasterProfileInput) *model.MasterProfileEntity {
	entity.DisplayName = input.DisplayName
	entity.Job = input.Job
	entity.Description = input.Description
	entity.Address = input.Address
	entity.Point = MapPointStringToPoint(input.Point)
	entity.AvatarID = input.AvatarID
	if input.IsEnabled != nil {
		entity.IsEnabled = *input.IsEnabled
	}
	return entity
}

func MapEntityToOwnMasterProfileOutput(entity *model.MasterProfileEntity) dto.OwnMasterProfileOutput {
	return dto.OwnMasterProfileOutput{
		DisplayName: entity.DisplayName,
		Job:         entity.Job,
		Description: entity.Description,
		Address:     entity.Address,
		Point:       MapPointToDomainPoint(entity.Point),
		AvatarID:    entity.AvatarID,
		IsEnabled:   entity.IsEnabled,
	}
}

func MapEntityToMasterProfileOutput(entity *model.MasterProfileEntity) dto.MasterProfileOutput {
	return dto.MasterProfileOutput{
		DisplayName: entity.DisplayName,
		Job:         entity.Job,
		Description: entity.Description,
		Address:     entity.Address,
		Point: dto.PointOutput{
			Latitude:  entity.Point.Lat,
			Longitude: entity.Point.Lng,
		},
		AvatarID: entity.AvatarID,
	}
}

func MapEntitiesToMasterProfileOutputs(entities []*model.MasterProfileEntity) []dto.MasterProfileOutput {
	outputs := make([]dto.MasterProfileOutput, len(entities))
	for i, entity := range entities {
		outputs[i] = MapEntityToMasterProfileOutput(entity)
	}
	return outputs
}

func MapEntitiesToMasterProfilesOutput(entities []*model.MasterProfileEntity, count int) dto.MasterProfilesOutput {
	return dto.MasterProfilesOutput{
		Count: count,
		Data:  MapEntitiesToMasterProfileOutputs(entities),
	}
}

func MapFindMasterProfilesToFilterMap(filter *dto.FindMasterProfilesFilterInput) map[model.MasterProfileFilter]interface{} {
	filterMap := make(map[model.MasterProfileFilter]interface{})
	if filter == nil {
		return filterMap
	}

	if filter.DisplayName != nil {
		filterMap[model.MasterProfileFilterDisplayName] = *filter.DisplayName
	}
	if filter.Job != nil {
		filterMap[model.MasterProfileFilterJob] = *filter.Job
	}
	if filter.Radius != nil && filter.Point != nil {
		point := MapPointStringToPoint(*filter.Point)
		radius, _ := strconv.ParseFloat(*filter.Radius, 64)
		filterMap[model.MasterProfileFilterPoint] = model.MasterProfileFilterPointValue{
			Latitude:  point.Lat,
			Longitude: point.Lng,
			Radius:    radius,
		}
	}

	return filterMap
}
