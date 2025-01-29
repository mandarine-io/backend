package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

var (
	// User errors
	ErrDuplicateUser = errors.New("duplicate user")

	// Master Profile errors
	ErrDuplicateMasterProfile       = errors.New("duplicate master profile")
	ErrUserForMasterProfileNotExist = errors.New("user for master profile does not exist")

	// Master Service errors
	ErrDuplicateMasterService       = errors.New("duplicate master service")
	ErrUserForMasterServiceNotExist = errors.New("user for master service does not exist")
)

type Scope func(db *gorm.DB) *gorm.DB

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID, scopes ...Scope) (*entity.User, error)
	FindUserByUsername(ctx context.Context, username string, scopes ...Scope) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string, scopes ...Scope) (*entity.User, error)
	FindUserByUsernameOrEmail(ctx context.Context, login string, scopes ...Scope) (*entity.User, error)
	ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error)
	DeleteExpiredUser(ctx context.Context) (*entity.User, error)

	WithRolePreload() Scope
}

type MasterProfileRepository interface {
	CreateMasterProfile(ctx context.Context, masterProfile *entity.MasterProfile) (*entity.MasterProfile, error)
	UpdateMasterProfile(ctx context.Context, masterProfile *entity.MasterProfile) (*entity.MasterProfile, error)
	FindMasterProfiles(ctx context.Context, scopes ...Scope) ([]*entity.MasterProfile, error)
	CountMasterProfiles(ctx context.Context, scopes ...Scope) (int64, error)
	FindMasterProfileByUserID(ctx context.Context, id uuid.UUID) (*entity.MasterProfile, error)
	FindMasterProfileByUsername(ctx context.Context, username string) (*entity.MasterProfile, error)
	FindEnabledMasterProfileByUsername(ctx context.Context, username string) (*entity.MasterProfile, error)
	ExistsMasterProfileByUserID(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsMasterProfileByUsername(ctx context.Context, username string) (bool, error)

	WithPagination(page, pageSize int) Scope
	WithColumnSort(field string, asc bool) Scope
	WithPointSort(latitude, longitude decimal.Decimal, asc bool) Scope
	WithDisplayNameFilter(displayName string) Scope
	WithJobFilter(job string) Scope
	WithPointFilter(latitude, longitude, radius decimal.Decimal) Scope
	WithAddressFilter(address string) Scope
}

type MasterServiceRepository interface {
	CreateMasterService(ctx context.Context, masterService *entity.MasterService) (*entity.MasterService, error)
	UpdateMasterService(ctx context.Context, masterService *entity.MasterService) (*entity.MasterService, error)
	DeleteMasterServiceByID(ctx context.Context, masterProfileID uuid.UUID, id uuid.UUID) error
	FindMasterServices(ctx context.Context, scopes ...Scope) ([]*entity.MasterService, error)
	CountMasterServices(ctx context.Context, scopes ...Scope) (int64, error)
	FindMasterServicesByMasterProfileID(
		ctx context.Context,
		masterProfileID uuid.UUID,
		scopes ...Scope,
	) ([]*entity.MasterService, error)
	CountMasterServicesByMasterProfileID(ctx context.Context, masterProfileID uuid.UUID, scopes ...Scope) (int64, error)
	FindMasterServiceByID(ctx context.Context, masterProfileID uuid.UUID, id uuid.UUID) (*entity.MasterService, error)
	ExistsMasterServiceByMasterID(ctx context.Context, masterID uuid.UUID) (bool, error)

	WithPagination(page, pageSize int) Scope
	WithSort(field string, asc bool) Scope
	WithNameFilter(name string) Scope
	WithMinPriceFilter(minPrice decimal.Decimal) Scope
	WithMaxPriceFilter(maxPrice decimal.Decimal) Scope
	WithMinIntervalFilter(minDuration time.Duration) Scope
	WithMaxIntervalFilter(maxDuration time.Duration) Scope
}
