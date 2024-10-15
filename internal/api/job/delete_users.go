package job

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"mandarine/internal/api/persistence/repo"
	"mandarine/pkg/logging"
	"mandarine/pkg/scheduler"
)

func deleteExpiredDeletedUsersJob(usersRepo repo.UserRepository) scheduler.Job {
	return scheduler.Job{
		Name:       "delete-expired-users",
		Definition: gocron.CronJob("0 0 1 * *", false),
		Task: gocron.NewTask(
			func() {
				_, err := usersRepo.DeleteExpiredUser(context.Background())
				if err != nil {
					slog.Error("Delete expired users error", logging.ErrorAttr(err))
				}
			},
		),
	}
}
