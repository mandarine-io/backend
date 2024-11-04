package registry

import (
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	gormRepo "github.com/mandarine-io/Backend/internal/api/persistence/repo/gorm"
	"github.com/rs/zerolog/log"
)

type Repositories struct {
	User        repo.UserRepository
	BannedToken repo.BannedTokenRepository
}

func setupGormRepositories(c *Container) {
	log.Debug().Msg("setup gorm repositories")
	c.Repos = &Repositories{
		User:        gormRepo.NewUserRepository(c.DB),
		BannedToken: gormRepo.NewBannedTokenRepository(c.DB),
	}
}
