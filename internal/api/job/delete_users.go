package job

import (
	"context"
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	"github.com/mandarine-io/Backend/pkg/scheduler"
)

func DeleteExpiredDeletedUsersJob(usersRepo repo.UserRepository) scheduler.Job {
	return scheduler.Job{
		Ctx:            context.Background(),
		Name:           "delete-expired-users",
		CronExpression: "0 0 1 * *",
		Action: func(ctx context.Context) error {
			_, err := usersRepo.DeleteExpiredUser(ctx)
			return err
		},
	}
}
