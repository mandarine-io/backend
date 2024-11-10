package util

import (
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ColumnSortScope(sorts []model.Sort) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		for _, sort := range sorts {
			tx.Order(clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{Column: clause.Column{Name: sort.Field}, Desc: sort.Order == model.SortOrderDesc},
				},
			})
		}

		return tx
	}
}
