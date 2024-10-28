package job

import (
	"context"
	"github.com/mandarine-io/Backend/internal/api/persistence/repo"
	"github.com/mandarine-io/Backend/pkg/scheduler"
)

func deleteExpiredTokensJob(bannedTokensRepo repo.BannedTokenRepository) scheduler.Job {
	return scheduler.Job{
		Ctx:            context.Background(),
		Name:           "delete-expired-tokens",
		CronExpression: "0 * * * *",
		Action: func(ctx context.Context) error {
			return bannedTokensRepo.DeleteExpiredBannedToken(ctx)
		},
	}
}
