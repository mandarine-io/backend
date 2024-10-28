package job

import (
	"github.com/mandarine-io/Backend/internal/api/registry"
	"github.com/mandarine-io/Backend/pkg/scheduler"
)

func SetupJobs(container *registry.Container) []scheduler.Job {
	return []scheduler.Job{
		deleteExpiredTokensJob(container.Repositories.BannedToken),
		deleteExpiredDeletedUsersJob(container.Repositories.User),
	}
}
