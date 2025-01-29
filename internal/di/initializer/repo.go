package initializer

import (
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/persistence/repo/gorm"
	"github.com/rs/zerolog/log"
)

func GormRepositories(c *di.Container) di.Initializer {
	return func() error {
		log.Debug().Msg("setup gorm repositories")

		c.Repos = di.Repositories{
			MasterProfile: gorm.NewMasterProfileRepository(
				c.Infrastructure.DB,
				gorm.WithMasterProfileRepoLogger(c.Logger.With().Str("repo", "master_profile").Logger()),
			),
			MasterService: gorm.NewMasterServiceRepository(
				c.Infrastructure.DB,
				gorm.WithMasterServiceRepoLogger(c.Logger.With().Str("repo", "master_service").Logger()),
			),
			User: gorm.NewUserRepository(
				c.Infrastructure.DB,
				gorm.WithUserRepoLogger(c.Logger.With().Str("repo", "user").Logger()),
			),
		}

		return nil
	}
}
