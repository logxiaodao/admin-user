package repositories

import (
	"gorm.io/gorm"
)

const DefaultLimit = 10000

// PageDefault  默认分页处理
func PageDefault(currentPage, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if currentPage == 0 {
			currentPage = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (currentPage - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}
