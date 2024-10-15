package job

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"mandarine/internal/api/persistence/repo"
	"mandarine/pkg/logging"
	"mandarine/pkg/scheduler"
)

func deleteExpiredTokensJob(bannedTokensRepo repo.BannedTokenRepository) scheduler.Job {
	return scheduler.Job{
		Name:       "delete-expired-tokens",
		Definition: gocron.CronJob("0 * * * *", false),
		Task: gocron.NewTask(
			func() {
				err := bannedTokensRepo.DeleteExpiredBannedToken(context.Background())
				if err != nil {
					slog.Error("Delete expired tokens error", logging.ErrorAttr(err))
				}
			},
		),
	}
}
