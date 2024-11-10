package repo

import (
	"context"
	"github.com/google/uuid"
	model2 "github.com/mandarine-io/Backend/internal/persistence/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model2.UserEntity) (*model2.UserEntity, error)
	UpdateUser(ctx context.Context, user *model2.UserEntity) (*model2.UserEntity, error)
	FindUserById(ctx context.Context, id uuid.UUID, rolePreload bool) (*model2.UserEntity, error)
	FindUserByUsername(ctx context.Context, username string, rolePreload bool) (*model2.UserEntity, error)
	FindUserByEmail(ctx context.Context, email string, rolePreload bool) (*model2.UserEntity, error)
	FindUserByUsernameOrEmail(ctx context.Context, login string, rolePreload bool) (*model2.UserEntity, error)
	ExistsUserById(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error)
	DeleteExpiredUser(ctx context.Context) (*model2.UserEntity, error)
}

type BannedTokenRepository interface {
	CreateOrUpdateBannedToken(ctx context.Context, bannedToken *model2.BannedTokenEntity) (*model2.BannedTokenEntity, error)
	ExistsBannedTokenByJTI(ctx context.Context, jti string) (bool, error)
	DeleteExpiredBannedToken(ctx context.Context) error
}

type MasterProfileRepository interface {
	CreateMasterProfile(ctx context.Context, masterProfile *model2.MasterProfileEntity) (*model2.MasterProfileEntity, error)
	UpdateMasterProfile(ctx context.Context, masterProfile *model2.MasterProfileEntity) (*model2.MasterProfileEntity, error)
	FindMasterProfiles(
		ctx context.Context,
		filter map[model2.MasterProfileFilter]interface{},
		pagination *model2.Pagination,
		sort []*model2.Sort,
	) ([]*model2.MasterProfileEntity, error)
	CountMasterProfiles(
		ctx context.Context,
		filter map[model2.MasterProfileFilter]interface{},
	) (int64, error)
	FindMasterProfileByUserId(ctx context.Context, id uuid.UUID) (*model2.MasterProfileEntity, error)
	FindEnabledMasterProfileByUsername(ctx context.Context, username string) (*model2.MasterProfileEntity, error)
	ExistsMasterProfileByUserId(ctx context.Context, id uuid.UUID) (bool, error)
}
