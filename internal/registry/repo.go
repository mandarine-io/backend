package registry

import (
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/mandarine-io/Backend/internal/persistence/repo/gorm"
	"github.com/rs/zerolog/log"
)

type Repositories struct {
	User          repo.UserRepository
	BannedToken   repo.BannedTokenRepository
	MasterProfile repo.MasterProfileRepository
}

func setupGormRepositories(c *Container) {
	log.Debug().Msg("setup gorm repositories")
	c.Repos = &Repositories{
		User:          gorm.NewUserRepository(c.DB),
		BannedToken:   gorm.NewBannedTokenRepository(c.DB),
		MasterProfile: gorm.NewMasterProfileRepository(c.DB),
	}
}
