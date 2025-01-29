package v0

import (
	"github.com/shopspring/decimal"
)

type CreateMasterServiceInput struct {
	Name        string           `json:"name" binding:"required"`
	Description *string          `json:"description" binding:"omitempty"`
	MinPrice    *decimal.Decimal `json:"minPrice" binding:"omitempty"`
	MaxPrice    *decimal.Decimal `json:"maxPrice" binding:"omitempty"`
	MinInterval *string          `json:"minInterval" format:"hh:mm:ss" binding:"omitempty,duration"`
	MaxInterval *string          `json:"maxInterval" format:"hh:mm:ss" binding:"omitempty,duration"`
	AvatarID    *string          `json:"avatarId" binding:"omitempty"`
}

type UpdateMasterServiceInput struct {
	Name        string           `json:"name" binding:"required"`
	Description *string          `json:"description" binding:"omitempty"`
	MinPrice    *decimal.Decimal `json:"minPrice" binding:"omitempty"`
	MaxPrice    *decimal.Decimal `json:"maxPrice" binding:"omitempty"`
	MinInterval *string          `json:"minInterval" format:"hh:mm:ss" binding:"omitempty,duration"`
	MaxInterval *string          `json:"maxInterval" format:"hh:mm:ss" binding:"omitempty,duration"`
	AvatarID    *string          `json:"avatarId" binding:"omitempty"`
}

type FindMasterServicesFilterInput struct {
	Name        *string          `form:"name" binding:"omitempty"`
	MinPrice    *decimal.Decimal `form:"minPrice" binding:"omitempty"`
	MaxPrice    *decimal.Decimal `form:"maxPrice" binding:"omitempty"`
	MinInterval *string          `form:"minInterval" binding:"omitempty,duration"`
	MaxInterval *string          `form:"maxInterval" binding:"omitempty,duration"`
}

type FindMasterServicesInput struct {
	*FindMasterServicesFilterInput
	*SortInput
	*PaginationInput
}

type MasterServiceOutput struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	MinPrice    *decimal.Decimal `json:"minPrice,omitempty"`
	MaxPrice    *decimal.Decimal `json:"maxPrice,omitempty"`
	MinInterval *string          `json:"minInterval"`
	MaxInterval *string          `json:"maxInterval"`
	AvatarID    *string          `json:"avatarId,omitempty"`
}

type MasterServicesOutput struct {
	Count int                   `json:"count"`
	Data  []MasterServiceOutput `json:"data"`
}
