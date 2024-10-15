package job

import (
	"mandarine/internal/api/registry"
	"mandarine/pkg/scheduler"
)

func SetupJobs(container *registry.Container) []scheduler.Job {
	return []scheduler.Job{
		deleteExpiredTokensJob(container.Repositories.BannedToken),
		deleteExpiredDeletedUsersJob(container.Repositories.User),
	}
}
