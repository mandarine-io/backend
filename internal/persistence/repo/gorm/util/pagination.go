package util

import (
	"gorm.io/gorm"
)

func PaginationScope(page, pageSize int) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(page * pageSize).Limit(pageSize)
	}
}
