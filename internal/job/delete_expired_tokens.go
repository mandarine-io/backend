package job

import (
	"context"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/mandarine-io/Backend/pkg/scheduler"
)

func DeleteExpiredTokensJob(bannedTokensRepo repo.BannedTokenRepository) scheduler.Job {
	return scheduler.Job{
		Ctx:            context.Background(),
		Name:           "delete-expired-tokens",
		CronExpression: "0 * * * *",
		Action: func(ctx context.Context) error {
			return bannedTokensRepo.DeleteExpiredBannedToken(ctx)
		},
	}
}
