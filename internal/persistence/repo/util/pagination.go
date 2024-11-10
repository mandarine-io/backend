package util

import (
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"gorm.io/gorm"
)

func PaginationScope(pagination *model.Pagination) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if pagination == nil {
			return tx
		}
		return tx.Offset(pagination.Offset).Limit(pagination.Limit)
	}
}
