package v0

import "github.com/shopspring/decimal"

type PaginationInput struct {
	Page     int `json:"page" form:"page,default=0" binding:"min=0"`
	PageSize int `json:"pageSize" form:"pageSize,default=10" binding:"min=1,max=100"`
}

type SortInput struct {
	Field string `json:"field" form:"field"`
	Order string `json:"order" form:"order" binding:"oneof=asc desc"`
}

type PointOutput struct {
	Longitude decimal.Decimal `json:"longitude"`
	Latitude  decimal.Decimal `json:"latitude"`
}
