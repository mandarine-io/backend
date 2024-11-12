package dto

type FindMasterProfilesFilterInput struct {
	DisplayName *string `form:"displayName" binding:"omitempty"`
	Job         *string `form:"job" binding:"omitempty"`
	Point       *string `form:"point" format:"lng,lat" binding:"omitempty,point"`
	Radius      *string `form:"radius" format:"float" binding:"omitempty,numeric"`
}

type FindMasterProfilesInput struct {
	*FindMasterProfilesFilterInput
	*SortInput
	*PaginationInput
}

type CreateMasterProfileInput struct {
	DisplayName string  `json:"displayName" binding:"required"`
	Job         string  `json:"job" binding:"required"`
	Description *string `json:"description" binding:"omitempty"`
	Address     *string `json:"address" binding:"omitempty"`
	Point       string  `json:"point" format:"lng,lat" binding:"required,point"`
	AvatarID    *string `json:"avatarId" binding:"omitempty"`
}

type UpdateMasterProfileInput struct {
	DisplayName string  `json:"displayName" binding:"required"`
	Job         string  `json:"job" binding:"required"`
	Description *string `json:"description" binding:"omitempty"`
	Address     *string `json:"address" binding:"omitempty"`
	Point       string  `json:"point" format:"lng,lat" binding:"required,point"`
	AvatarID    *string `json:"avatarId" binding:"omitempty"`
	IsEnabled   *bool   `json:"isEnabled" binding:"required"`
}

type MasterProfileOutput struct {
	DisplayName string      `json:"displayName" binding:"required"`
	Job         string      `json:"job" binding:"required"`
	Description *string     `json:"description" binding:"omitempty"`
	Address     *string     `json:"address" binding:"omitempty"`
	Point       PointOutput `json:"point" binding:"required"`
	AvatarID    *string     `json:"avatarId" binding:"omitempty"`
}

type MasterProfilesOutput struct {
	Count int                   `json:"count" binding:"required"`
	Data  []MasterProfileOutput `json:"data" binding:"required"`
}

type OwnMasterProfileOutput struct {
	DisplayName string      `json:"displayName" binding:"required"`
	Job         string      `json:"job" binding:"required"`
	Description *string     `json:"description" binding:"omitempty"`
	Address     *string     `json:"address" binding:"omitempty"`
	Point       PointOutput `json:"point" binding:"required"`
	AvatarID    *string     `json:"avatarId" binding:"omitempty"`
	IsEnabled   bool        `json:"isEnabled" binding:"required"`
}
