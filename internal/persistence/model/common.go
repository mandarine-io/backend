package model

type Pagination struct {
	Offset int
	Limit  int
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type Sort struct {
	Field string
	Order SortOrder
}
