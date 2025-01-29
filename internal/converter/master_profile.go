package converter

import (
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/types"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func MapLngLatToPoint(lng, lat decimal.Decimal) *types.Point {
	return types.NewPoint(lat, lng)
}

func MapPointToDomainPoint(point types.Point) v0.PointOutput {
	return v0.PointOutput{
		Latitude:  point.Lat,
		Longitude: point.Lng,
	}
}

func MapCreateMasterProfileInputToEntity(userID uuid.UUID, input v0.CreateMasterProfileInput) *entity.MasterProfile {
	return &entity.MasterProfile{
		UserID:      userID,
		DisplayName: input.DisplayName,
		Job:         input.Job,
		Description: input.Description,
		Address:     input.Address,
		Point:       *MapLngLatToPoint(input.Longitude, input.Latitude),
		AvatarID:    input.AvatarID,
		IsEnabled:   true,
	}
}

func MapUpdateMasterProfileInputToEntity(
	entity *entity.MasterProfile,
	input v0.UpdateMasterProfileInput,
) *entity.MasterProfile {
	entity.DisplayName = input.DisplayName
	entity.Job = input.Job
	entity.Description = input.Description
	entity.Address = input.Address
	entity.Point = *MapLngLatToPoint(input.Longitude, input.Latitude)
	entity.AvatarID = input.AvatarID
	return entity
}

func MapEntityToMasterProfileOutput(entity *entity.MasterProfile) v0.MasterProfileOutput {
	return v0.MasterProfileOutput{
		DisplayName: entity.DisplayName,
		Job:         entity.Job,
		Description: entity.Description,
		Address:     entity.Address,
		Point:       MapPointToDomainPoint(entity.Point),
		AvatarID:    entity.AvatarID,
		IsEnabled:   lo.ToPtr(entity.IsEnabled),
	}
}

func MapEntitiesToMasterProfileOutputs(entities []*entity.MasterProfile) []v0.MasterProfileOutput {
	outputs := make([]v0.MasterProfileOutput, len(entities))
	for i, e := range entities {
		outputs[i] = MapEntityToMasterProfileOutput(e)
	}
	return outputs
}

func MapEntitiesToMasterProfilesOutput(entities []*entity.MasterProfile, count int) v0.MasterProfilesOutput {
	return v0.MasterProfilesOutput{
		Count: count,
		Data:  MapEntitiesToMasterProfileOutputs(entities),
	}
}
