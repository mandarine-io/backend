package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/mandarine-io/backend/internal/persistence/repo/gorm/util"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type masterServiceRepo struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type MasterServiceRepoOption func(*masterServiceRepo)

func WithMasterServiceRepoLogger(logger zerolog.Logger) MasterServiceRepoOption {
	return func(r *masterServiceRepo) {
		r.logger = logger
	}
}

func NewMasterServiceRepository(db *gorm.DB, opts ...MasterServiceRepoOption) repo.MasterServiceRepository {
	r := &masterServiceRepo{
		db:     db,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *masterServiceRepo) CreateMasterService(
	ctx context.Context,
	masterService *entity.MasterService,
) (*entity.MasterService, error) {
	r.logger.Debug().Msg("create master service")

	tx := r.db.WithContext(ctx).Create(masterService)

	err := tx.Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return masterService, repo.ErrDuplicateMasterService
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return masterService, repo.ErrUserForMasterServiceNotExist
	}

	return masterService, err
}

func (r *masterServiceRepo) UpdateMasterService(
	ctx context.Context,
	masterService *entity.MasterService,
) (*entity.MasterService, error) {
	r.logger.Debug().Msg("update master service")

	tx := r.db.WithContext(ctx).Save(masterService)

	return masterService, tx.Error
}

func (r *masterServiceRepo) DeleteMasterServiceByID(
	ctx context.Context,
	masterProfileID uuid.UUID,
	id uuid.UUID,
) error {
	r.logger.Debug().Msg("delete master service")

	tx := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Where("master_profile_id = ?", masterProfileID).
		Delete(&entity.MasterService{})

	return tx.Error
}

func (r *masterServiceRepo) FindMasterServices(ctx context.Context, scopes ...repo.Scope) (
	[]*entity.MasterService,
	error,
) {
	r.logger.Debug().Msg("find master services")

	tx := r.db.WithContext(ctx)

	for _, scope := range scopes {
		tx = tx.Scopes(scope)
	}

	var masterServices []*entity.MasterService
	err := tx.
		Joins("join master_profiles on master_profiles.user_id = master_services.master_profile_id").
		Where("master_profiles.is_enabled = ?", true).
		Find(&masterServices).Error

	if masterServices == nil {
		masterServices = make([]*entity.MasterService, 0)
	}

	return masterServices, err
}

func (r *masterServiceRepo) CountMasterServices(ctx context.Context, scopes ...repo.Scope) (int64, error) {
	r.logger.Debug().Msg("count master services")

	tx := r.db.WithContext(ctx)

	for _, scope := range scopes {
		tx = tx.Scopes(scope)
	}

	var count int64
	err := tx.
		Select("count(*)").
		Find(&count).
		Error

	return count, err
}

func (r *masterServiceRepo) FindMasterServicesByMasterProfileID(
	ctx context.Context,
	masterProfileID uuid.UUID,
	scopes ...repo.Scope,
) ([]*entity.MasterService, error) {
	r.logger.Debug().Msg("find master services by master profile ID")

	tx := r.db.WithContext(ctx)

	for _, scope := range scopes {
		tx = tx.Scopes(scope)
	}

	var masterServices []*entity.MasterService
	err := tx.
		Where("master_profile_id = ?", masterProfileID).
		Find(&masterServices).
		Error

	if masterServices == nil {
		masterServices = make([]*entity.MasterService, 0)
	}

	return masterServices, err
}

func (r *masterServiceRepo) CountMasterServicesByMasterProfileID(
	ctx context.Context,
	masterProfileID uuid.UUID,
	scopes ...repo.Scope,
) (int64, error) {
	r.logger.Debug().Msg("count master services by master profile ID")

	tx := r.db.WithContext(ctx)

	for _, scope := range scopes {
		tx = tx.Scopes(scope)
	}

	var count int64
	err := tx.
		Where("master_profile_id = ?", masterProfileID).
		Select("count(*)").
		Find(&count).
		Error

	return count, err
}

func (r *masterServiceRepo) FindMasterServiceByID(
	ctx context.Context,
	masterProfileID uuid.UUID,
	id uuid.UUID,
) (*entity.MasterService, error) {
	r.logger.Debug().Msg("find master service by id")

	masterService := &entity.MasterService{}

	tx := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Where("master_profile_id = ?", masterProfileID).
		First(masterService)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return masterService, tx.Error
}

func (r *masterServiceRepo) ExistsMasterServiceByMasterID(ctx context.Context, masterID uuid.UUID) (bool, error) {
	r.logger.Debug().Msg("exists master profile by master id")

	var exists bool
	tx := r.db.WithContext(ctx).
		Model(&entity.MasterProfile{}).
		Select("count(*) > 0").
		Where("master_id = ?", masterID).
		Find(&exists)

	return exists, tx.Error
}

func (r *masterServiceRepo) WithPagination(page, pageSize int) repo.Scope {
	return util.PaginationScope(page, pageSize)
}

func (r *masterServiceRepo) WithSort(field string, asc bool) repo.Scope {
	return util.ColumnSortScope(field, asc)
}

func (r *masterServiceRepo) WithNameFilter(name string) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or("master_services.name ILIKE ?", fmt.Sprintf("%%%s%%", name))
		return tx
	}
}

func (r *masterServiceRepo) WithMinPriceFilter(minPrice decimal.Decimal) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or(
			"master_services.min_price IS NULL OR "+
				"(master_services.min_price IS NOT NULL AND master_services.min_price <= ?)",
			minPrice,
		)
		return tx
	}
}

func (r *masterServiceRepo) WithMaxPriceFilter(maxPrice decimal.Decimal) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or(
			"master_services.max_price IS NULL OR "+
				"(master_services.max_price IS NOT NULL AND master_services.max_price >= ?)",
			maxPrice,
		)
		return tx
	}
}

func (r *masterServiceRepo) WithMinIntervalFilter(minInterval time.Duration) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or(
			"master_services.min_interval IS NULL OR "+
				"(master_services.min_interval IS NOT NULL AND master_services.min_interval <= ?)",
			minInterval,
		)
		return tx
	}
}

func (r *masterServiceRepo) WithMaxIntervalFilter(maxInterval time.Duration) repo.Scope {
	return func(tx *gorm.DB) *gorm.DB {
		tx.Or(
			"master_services.max_interval IS NULL OR "+
				"(master_services.max_interval IS NOT NULL AND master_services.max_interval <= ?)",
			maxInterval,
		)
		return tx
	}
}
