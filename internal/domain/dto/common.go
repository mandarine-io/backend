package dto

const (
	MaxPageSize = 100
)

type PaginationInput struct {
	Page     int `json:"page" form:"page,default=0" binding:"min=0"`
	PageSize int `json:"pageSize" form:"pageSize,default=10" binding:"min=1,max=100"`
}

type SortInput struct {
	Field string `json:"field" form:"field" binding:"required"`
	Order string `json:"order" form:"order" binding:"omitempty,oneof=asc desc"`
}

type PointOutput struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
