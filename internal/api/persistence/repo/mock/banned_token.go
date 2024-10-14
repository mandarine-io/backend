package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mandarine/internal/api/persistence/model"
)

type BannedTokenRepositoryMock struct {
	mock.Mock
}

func (b *BannedTokenRepositoryMock) CreateOrUpdateBannedToken(ctx context.Context, bannedToken *model.BannedTokenEntity) (*model.BannedTokenEntity, error) {
	args := b.Called(ctx, bannedToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BannedTokenEntity), args.Error(1)
}

func (b *BannedTokenRepositoryMock) ExistsBannedTokenByJTI(ctx context.Context, jti string) (bool, error) {
	args := b.Called(ctx, jti)
	return args.Get(0).(bool), args.Error(1)
}

func (b *BannedTokenRepositoryMock) DeleteExpiredBannedToken(ctx context.Context) error {
	args := b.Called(ctx)
	return args.Error(0)
}
