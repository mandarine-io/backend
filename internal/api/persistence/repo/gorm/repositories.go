package gorm

import (
	"gorm.io/gorm"
	"mandarine/internal/api/persistence/repo"
)

func NewRepositories(db *gorm.DB) *repo.Repositories {
	return &repo.Repositories{
		User:        NewUserRepository(db),
		BannedToken: NewBannedTokenRepository(db),
	}
}
