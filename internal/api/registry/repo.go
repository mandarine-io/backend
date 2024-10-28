package registry

import (
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	gormRepo "github.com/mandarine-io/Backend/internal/api/persistence/repo/gorm"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repositories struct {
	User        repo.UserRepository
	BannedToken repo.BannedTokenRepository
}

func newGormRepositories(db *gorm.DB) *Repositories {
	log.Debug().Msg("setup gorm repositories")
	return &Repositories{
		User:        gormRepo.NewUserRepository(db),
		BannedToken: gormRepo.NewBannedTokenRepository(db),
	}
}
