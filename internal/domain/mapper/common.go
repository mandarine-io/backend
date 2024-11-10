package mapper

import (
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"github.com/pkg/errors"
	"strings"
)

func MapDomainPaginationToPagination(pagination *dto.PaginationInput) *model.Pagination {
	if pagination == nil {
		return &model.Pagination{
			Limit:  10,
			Offset: 0,
		}
	}

	pageSize := minInt(pagination.PageSize, dto.MaxPageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	page := pagination.Page
	if page < 0 {
		page = 0
	}

	return &model.Pagination{
		Limit:  pageSize,
		Offset: page * pageSize,
	}
}

func MapDomainSortToSort(sort *dto.SortInput) *model.Sort {
	if sort == nil || sort.Field == "" {
		return nil
	}

	order := model.SortOrderAsc
	if strings.ToLower(sort.Order) == "desc" {
		order = model.SortOrderDesc
	}
	return &model.Sort{Field: sort.Field, Order: order}
}

func MapDomainSortsToSorts(sorts []*dto.SortInput) []*model.Sort {
	mappedSorts := make([]*model.Sort, len(sorts))
	for i, sort := range sorts {
		if sort == nil {
			continue
		}
		mappedSorts[i] = MapDomainSortToSort(sort)
	}
	return mappedSorts
}

func MapDomainSortsToSortsWithAvailables(sorts []*dto.SortInput, availableFields []string) ([]*model.Sort, error) {
	availableFiledMap := make(map[string]interface{})
	for _, field := range availableFields {
		availableFiledMap[field] = true
	}

	mappedSorts := make([]*model.Sort, len(sorts))
	for i, sort := range sorts {
		if sort == nil {
			continue
		}
		if _, ok := availableFiledMap[sort.Field]; !ok {
			return nil, errors.New("invalid sort field")
		}
		mappedSorts[i] = MapDomainSortToSort(sort)
	}
	return mappedSorts, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
