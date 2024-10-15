package repo

import (
	"context"
	"github.com/google/uuid"
	"mandarine/internal/api/persistence/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error)
	UpdateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error)
	FindUserById(ctx context.Context, id uuid.UUID, rolePreload bool) (*model.UserEntity, error)
	FindUserByUsername(ctx context.Context, username string, rolePreload bool) (*model.UserEntity, error)
	FindUserByEmail(ctx context.Context, email string, rolePreload bool) (*model.UserEntity, error)
	FindUserByUsernameOrEmail(ctx context.Context, login string, rolePreload bool) (*model.UserEntity, error)
	ExistsUserById(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error)
	DeleteExpiredUser(ctx context.Context) (*model.UserEntity, error)
}

type BannedTokenRepository interface {
	CreateOrUpdateBannedToken(ctx context.Context, bannedToken *model.BannedTokenEntity) (*model.BannedTokenEntity, error)
	ExistsBannedTokenByJTI(ctx context.Context, jti string) (bool, error)
	DeleteExpiredBannedToken(ctx context.Context) error
}

type Repositories struct {
	User        UserRepository
	BannedToken BannedTokenRepository
}
