package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/persistence/repo/gorm/util"
	"github.com/mandarine-io/backend/internal/persistence/types"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type masterProfileRepo struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type MasterProfileRepoOption func(*masterProfileRepo)

func WithMasterProfileRepoLogger(logger zerolog.Logger) MasterProfileRepoOption {
	return func(r *masterProfileRepo) {
		r.logger = logger
	}
}

func NewMasterProfileRepository(db *gorm.DB, opts ...MasterProfileRepoOption) repo.MasterProfileRepository {
	r := &masterProfileRepo{
		db:     db,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *masterProfileRepo) CreateMasterProfile(
	ctx context.Context,
	masterProfile *entity.MasterProfile,
) (*entity.MasterProfile, error) {
	r.logger.Debug().Msg("create master profile")

	tx := r.db.WithContext(ctx).Create(masterProfile)

	err := tx.Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return masterProfile, repo.ErrDuplicateMasterProfile
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return masterProfile, repo.ErrUserForMasterProfileNotExist
	}

	return masterProfile, err
}

func (r *masterProfileRepo) UpdateMasterProfile(
	ctx context.Context,
	masterProfile *entity.MasterProfile,
) (*entity.MasterProfile, error) {
	r.logger.Debug().Msg("update master profile")

	tx := r.db.WithContext(ctx).Save(masterProfile)

	return masterProfile, tx.Error
}

func (r *masterProfileRepo) FindMasterProfiles(ctx context.Context, scopes ...repo.Scope) (
	[]*entity.MasterProfile,
	error,
) {
	r.logger.Debug().Msg("find master profiles")

	tx := r.db.WithContext(ctx)

	for _, scope := range scopes {
		tx.Scopes(scope)
	}

	var masterProfiles []*entity.MasterProfile
	err := tx.Find(&masterProfiles).Error

	if masterProfiles == nil {
		masterProfiles = make([]*entity.MasterProfile, 0)
	}

	return masterProfiles, err
}

func (r *masterProfileRepo) CountMasterProfiles(ctx context.Context, scopes ...repo.Scope) (int64, error) {
	r.logger.Debug().Msg("count master profiles")

	tx := r.db.WithContext(ctx)

	for _, scope := range scopes {
		tx.Scopes(scope)
	}

	var count int64
	err := tx.
		Model(&entity.MasterProfile{}).
		Select("count(*)").
		Find(&count).
		Error

	return count, err
}

func (r *masterProfileRepo) FindMasterProfileByUserID(ctx context.Context, userID uuid.UUID) (
	*entity.MasterProfile,
	error,
) {
	r.logger.Debug().Msg("find master profile by user id")

	masterProfile := &entity.MasterProfile{}
	tx := r.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		First(masterProfile)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return masterProfile, tx.Error
}

func (r *masterProfileRepo) FindMasterProfileByUsername(ctx context.Context, username string) (
	*entity.MasterProfile,
	error,
) {
	r.logger.Debug().Msg("find master profile by username")

	masterProfile := &entity.MasterProfile{}
	tx := r.db.
		WithContext(ctx).
		Joins("join users on users.id = master_profiles.user_id").
		Where("users.username = ?", username).
		First(masterProfile)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return masterProfile, tx.Error
}

func (r *masterProfileRepo) FindEnabledMasterProfileByUsername(
	ctx context.Context,
	username string,
) (*entity.MasterProfile, error) {
	r.logger.Debug().Msg("find enabled master profile by username")

	masterProfile := &entity.MasterProfile{}
	tx := r.db.
		WithContext(ctx).
		Joins("join users on users.id = master_profiles.user_id").
		Where("users.username = ?", username).
		Where("master_profiles.is_enabled = ?", true).
		First(masterProfile)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return masterProfile, tx.Error
}

func (r *masterProfileRepo) ExistsMasterProfileByUserID(ctx context.Context, id uuid.UUID) (bool, error) {
	r.logger.Debug().Msg("exists master profile by user id")

	var exists bool
	tx := r.db.WithContext(ctx).
		Model(&entity.MasterProfile{}).
		Select("count(*) > 0").
		Where("user_id = ?", id).
		Find(&exists)

	return exists, tx.Error
}

func (r *masterProfileRepo) ExistsMasterProfileByUsername(ctx context.Context, username string) (bool, error) {
	r.logger.Debug().Msg("exists master profile by user id")

	var exists bool
	tx := r.db.
		WithContext(ctx).
		Model(&entity.MasterProfile{}).
		Joins("join users on users.id = master_profiles.user_id").
		Select("count(*) > 0").
		Where("users.username = ?", username).
		Find(&exists)

	return exists, tx.Error
}

func (r *masterProfileRepo) WithPagination(page, pageSize int) repo.Scope {
	return util.PaginationScope(page, pageSize)
}

func (r *masterProfileRepo) WithColumnSort(field string, asc bool) repo.Scope {
	return util.ColumnSortScope(field, asc)
}

func (r *masterProfileRepo) WithPointSort(latitude, longitude decimal.Decimal, asc bool) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		center := types.NewPoint(latitude, longitude)

		orderStr := ""
		if !asc {
			orderStr = " DESC"
		}

		return tx.
			Group("user_id").
			Order(
				clause.OrderBy{
					Expression: clause.Expr{
						SQL:                "st_distance(master_services.point, ?)" + orderStr,
						Vars:               []any{center},
						WithoutParentheses: true,
					},
				},
			)
	}
}

func (r *masterProfileRepo) WithDisplayNameFilter(displayName string) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or("master_services.name ILIKE ?", fmt.Sprintf("%%%s%%", displayName))
		return tx
	}
}

func (r *masterProfileRepo) WithJobFilter(job string) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or("master_services.name ILIKE ?", fmt.Sprintf("%%%s%%", job))
		return tx
	}
}

func (r *masterProfileRepo) WithPointFilter(latitude, longitude decimal.Decimal, radius decimal.Decimal) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		center := types.NewPoint(latitude, longitude)
		tx.Or("ST_DWithin(master_services.point::GEOGRAPHY, ?::GEOGRAPHY, ?)", center, radius)
		return tx
	}
}

func (r *masterProfileRepo) WithAddressFilter(address string) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or("master_services.name ILIKE ?", fmt.Sprintf("%%%s%%", address))
		return tx
	}
}
