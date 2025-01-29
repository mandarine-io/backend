package job

import (
	"context"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/scheduler"
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
