package registry

import (
	"gorm.io/gorm"
	"mandarine/internal/api/persistence/repo"
	gormRepo "mandarine/internal/api/persistence/repo/gorm"
)

type Repositories struct {
	User        repo.UserRepository
	BannedToken repo.BannedTokenRepository
}

func newGormRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:        gormRepo.NewUserRepository(db),
		BannedToken: gormRepo.NewBannedTokenRepository(db),
	}
}
