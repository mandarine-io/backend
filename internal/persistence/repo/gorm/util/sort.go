package util

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ColumnSortScope(column string, asc bool) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Order(
			clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{Column: clause.Column{Name: column}, Desc: !asc},
				},
			},
		)

		return tx
	}
}
