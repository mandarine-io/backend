package v0

import "github.com/shopspring/decimal"

type FindMasterProfilesFilterInput struct {
	DisplayName *string          `form:"displayName" binding:"omitempty"`
	Job         *string          `form:"job" binding:"omitempty"`
	Lng         *decimal.Decimal `form:"lng" format:"decimal" binding:"omitempty"`
	Lat         *decimal.Decimal `form:"lat" format:"decimal" binding:"omitempty"`
	Radius      *decimal.Decimal `form:"radius" format:"decimal" binding:"omitempty"`
}

type FindMasterProfilesInput struct {
	*FindMasterProfilesFilterInput
	*SortInput
	*PaginationInput
}

type CreateMasterProfileInput struct {
	DisplayName string          `json:"displayName" binding:"required"`
	Job         string          `json:"job" binding:"required"`
	Description *string         `json:"description" binding:"omitempty"`
	Address     *string         `json:"address" binding:"omitempty"`
	Longitude   decimal.Decimal `json:"longitude" format:"decimal" binding:"required"`
	Latitude    decimal.Decimal `json:"latitude" format:"decimal" binding:"required"`
	AvatarID    *string         `json:"avatarId" binding:"omitempty"`
}

type UpdateMasterProfileInput struct {
	DisplayName string          `json:"displayName" binding:"required"`
	Job         string          `json:"job" binding:"required"`
	Description *string         `json:"description" binding:"omitempty"`
	Address     *string         `json:"address" binding:"omitempty"`
	Longitude   decimal.Decimal `json:"longitude" format:"decimal" binding:"required"`
	Latitude    decimal.Decimal `json:"latitude" format:"decimal" binding:"required"`
	AvatarID    *string         `json:"avatarId" binding:"omitempty"`
}

type SwitchMasterProfileInput struct {
	IsEnabled bool `form:"isEnabled" binding:"required"`
}

type MasterProfileOutput struct {
	DisplayName string      `json:"displayName" binding:"required"`
	Job         string      `json:"job" binding:"required"`
	Description *string     `json:"description" binding:"omitempty"`
	Address     *string     `json:"address" binding:"omitempty"`
	Point       PointOutput `json:"point" binding:"required"`
	AvatarID    *string     `json:"avatarId" binding:"omitempty"`
	IsEnabled   *bool       `json:"isEnabled" binding:"omitempty"`
}

type MasterProfilesOutput struct {
	Count int                   `json:"count" binding:"required"`
	Data  []MasterProfileOutput `json:"data" binding:"required"`
}
