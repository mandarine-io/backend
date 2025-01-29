package converter

import (
	"github.com/mandarine-io/backend/internal/persistence/entity"
	timeutil "github.com/mandarine-io/backend/internal/util/time"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"time"
)

func MapCreateMasterServiceInputToEntity(input v0.CreateMasterServiceInput) *entity.MasterService {
	var (
		minInterval *time.Duration = nil
		maxInterval *time.Duration = nil
	)

	if input.MinInterval != nil {
		minIntervalTmp, err := time.ParseDuration(*input.MinInterval)
		if err == nil {
			minInterval = &minIntervalTmp
		}
	}

	if input.MaxInterval != nil {
		maxIntervalTmp, err := time.ParseDuration(*input.MaxInterval)
		if err == nil {
			maxInterval = &maxIntervalTmp
		}
	}

	return &entity.MasterService{
		Name:        input.Name,
		Description: input.Description,
		MinInterval: minInterval,
		MaxInterval: maxInterval,
		MinPrice:    input.MinPrice,
		MaxPrice:    input.MaxPrice,
		AvatarID:    input.AvatarID,
	}
}

func MapUpdateMasterServiceInputToEntity(
	entity *entity.MasterService,
	input v0.UpdateMasterServiceInput,
) *entity.MasterService {
	var (
		minInterval *time.Duration = nil
		maxInterval *time.Duration = nil
	)

	if input.MinInterval != nil {
		minIntervalTmp, err := time.ParseDuration(*input.MinInterval)
		if err == nil {
			minInterval = &minIntervalTmp
		}
	}

	if input.MaxInterval != nil {
		maxIntervalTmp, err := time.ParseDuration(*input.MaxInterval)
		if err == nil {
			maxInterval = &maxIntervalTmp
		}
	}

	entity.Name = input.Name
	entity.Description = input.Description
	entity.MinInterval = minInterval
	entity.MaxInterval = maxInterval
	entity.MinPrice = input.MinPrice
	entity.MaxPrice = input.MaxPrice
	entity.AvatarID = input.AvatarID
	return entity
}

func MapEntityToMasterServiceOutput(entity *entity.MasterService) v0.MasterServiceOutput {
	var (
		minInterval *string = nil
		maxInterval *string = nil
	)

	if entity.MinInterval != nil {
		minIntervalTmp := timeutil.FormatDuration(*entity.MinInterval)
		minInterval = &minIntervalTmp
	}

	if entity.MaxInterval != nil {
		maxIntervalTmp := timeutil.FormatDuration(*entity.MaxInterval)
		maxInterval = &maxIntervalTmp
	}

	return v0.MasterServiceOutput{
		ID:          entity.ID.String(),
		Name:        entity.Name,
		Description: entity.Description,
		MinInterval: minInterval,
		MaxInterval: maxInterval,
		MinPrice:    entity.MinPrice,
		MaxPrice:    entity.MaxPrice,
		AvatarID:    entity.AvatarID,
	}
}

func MapEntitiesToMasterServiceOutputs(entities []*entity.MasterService) []v0.MasterServiceOutput {
	outputs := make([]v0.MasterServiceOutput, len(entities))
	for i, e := range entities {
		outputs[i] = MapEntityToMasterServiceOutput(e)
	}
	return outputs
}

func MapEntitiesToMasterServicesOutput(entities []*entity.MasterService, count int) v0.MasterServicesOutput {
	return v0.MasterServicesOutput{
		Count: count,
		Data:  MapEntitiesToMasterServiceOutputs(entities),
	}
}
