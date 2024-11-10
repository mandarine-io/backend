package gorm

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	"github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/mandarine-io/Backend/internal/persistence/repo/util"
	gormtype "github.com/mandarine-io/Backend/internal/persistence/type"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type masterProfileRepo struct {
	db *gorm.DB
}

func NewMasterProfileRepository(db *gorm.DB) repo.MasterProfileRepository {
	return &masterProfileRepo{db}
}

func (m *masterProfileRepo) CreateMasterProfile(ctx context.Context, masterProfile *model.MasterProfileEntity) (*model.MasterProfileEntity, error) {
	log.Debug().Msg("create master profile")
	tx := m.db.WithContext(ctx).Create(masterProfile)

	err := tx.Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return masterProfile, repo.ErrDuplicateMasterProfile
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return masterProfile, repo.ErrUserForMasterProfileNotExist
	}

	return masterProfile, err
}

func (m *masterProfileRepo) UpdateMasterProfile(ctx context.Context, masterProfile *model.MasterProfileEntity) (*model.MasterProfileEntity, error) {
	log.Debug().Msg("update master profile")
	tx := m.db.WithContext(ctx).Save(masterProfile)
	return masterProfile, tx.Error
}

func (m *masterProfileRepo) FindMasterProfiles(
	ctx context.Context,
	filter map[model.MasterProfileFilter]interface{},
	pagination *model.Pagination,
	sort []*model.Sort,
) ([]*model.MasterProfileEntity, error) {
	log.Debug().Msg("find master profiles")
	var masterProfiles []*model.MasterProfileEntity

	tx := m.db.
		WithContext(ctx).
		Scopes(util.PaginationScope(pagination)).
		Scopes(filterScope(filter)).
		Scopes(sortScope(filter, sort))

	err := tx.Find(&masterProfiles).Error
	if masterProfiles == nil {
		masterProfiles = make([]*model.MasterProfileEntity, 0)
	}

	return masterProfiles, err
}

func (m *masterProfileRepo) CountMasterProfiles(
	ctx context.Context,
	filter map[model.MasterProfileFilter]interface{},
) (int64, error) {
	log.Debug().Msg("count master profiles")
	var count int64

	err := m.db.
		Model(&model.MasterProfileEntity{}).
		WithContext(ctx).
		Scopes(filterScope(filter)).
		Select("count(*)").
		Find(&count).
		Error
	return count, err
}

func (m *masterProfileRepo) FindMasterProfileByUserId(ctx context.Context, userId uuid.UUID) (*model.MasterProfileEntity, error) {
	log.Debug().Msg("find master profile by user id")
	masterProfile := &model.MasterProfileEntity{}

	tx := m.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		First(masterProfile)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return masterProfile, tx.Error
}

func (m *masterProfileRepo) FindEnabledMasterProfileByUsername(ctx context.Context, username string) (*model.MasterProfileEntity, error) {
	log.Debug().Msg("find master profile by username")
	masterProfile := &model.MasterProfileEntity{}

	tx := m.db.
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

func (m *masterProfileRepo) ExistsMasterProfileByUserId(ctx context.Context, id uuid.UUID) (bool, error) {
	log.Debug().Msg("exists master profile by user id")
	var exists bool
	tx := m.db.WithContext(ctx).
		Model(&model.MasterProfileEntity{}).
		Select("count(*) > 0").
		Where("user_id = ?", id).
		Find(&exists)
	return exists, tx.Error
}

func filterScope(filter map[model.MasterProfileFilter]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		for key, value := range filter {
			switch key {
			case model.MasterProfileFilterDisplayName:
				value, ok := value.(string)
				if !ok {
					continue
				}
				tx.Or("master_profiles.display_name ILIKE ?", fmt.Sprintf("%%%s%%", value))
			case model.MasterProfileFilterJob:
				value, ok := value.(string)
				if !ok {
					continue
				}
				tx.Or("master_profiles.job ILIKE ?", fmt.Sprintf("%%%s%%", value))
			case model.MasterProfileFilterAddress:
				value, ok := value.(string)
				if !ok {
					continue
				}
				tx.Or("master_profiles.address ILIKE ?", fmt.Sprintf("%%%s%%", value))
			case model.MasterProfileFilterPoint:
				value, ok := value.(model.MasterProfileFilterPointValue)
				if !ok {
					continue
				}
				center := gormtype.NewPoint(value.Latitude, value.Longitude)
				tx.Or("st_dwithin(point, ?, ?)", center, value.Radius)
			}
		}

		return tx
	}
}

func sortScope(filter map[model.MasterProfileFilter]interface{}, sorts []*model.Sort) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		for _, s := range sorts {
			if s == nil {
				continue
			}

			desc := s.Order == model.SortOrderDesc

			switch model.MasterProfileFilter(s.Field) {
			case model.MasterProfileFilterDisplayName:
				tx.Scopes(displayNameSortScope(desc))
			case model.MasterProfileFilterJob:
				tx.Scopes(jobSortScope(desc))
			case model.MasterProfileFilterAddress:
				tx.Scopes(addressSortScope(desc))
			case model.MasterProfileFilterPoint:
				value, ok := filter[model.MasterProfileFilter(s.Field)]
				if !ok {
					continue
				}
				tx.Scopes(pointSortScope(value, desc))
			}
		}

		return tx
	}
}

func displayNameSortScope(desc bool) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.
			Order(clause.OrderByColumn{
				Column: clause.Column{Name: "master_profiles.display_name"},
				Desc:   desc,
			})
	}
}

func jobSortScope(desc bool) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.
			Order(clause.OrderByColumn{
				Column: clause.Column{Name: "master_profiles.job"},
				Desc:   desc,
			})
	}
}

func addressSortScope(desc bool) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.
			Order(clause.OrderByColumn{
				Column: clause.Column{Name: "master_profiles.address"},
				Desc:   desc,
			})
	}
}

func pointSortScope(value interface{}, desc bool) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		point, ok := value.(model.MasterProfileFilterPointValue)
		if !ok {
			return tx
		}
		center := gormtype.NewPoint(point.Latitude, point.Longitude)

		orderStr := ""
		if desc {
			orderStr = " DESC"
		}

		return tx.
			Group("user_id").
			Order(clause.OrderBy{
				Expression: clause.Expr{
					SQL:                "st_distance(point, ?)" + orderStr,
					Vars:               []interface{}{center},
					WithoutParentheses: true,
				},
			})
	}
}
