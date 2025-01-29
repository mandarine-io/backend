package gorm

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mandarine-io/backend/internal/persistence/entity"
	"github.com/mandarine-io/backend/internal/persistence/repo"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type userRepo struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type UserRepoOption func(*userRepo)

func WithUserRepoLogger(logger zerolog.Logger) UserRepoOption {
	return func(r *userRepo) {
		r.logger = logger
	}
}

func NewUserRepository(db *gorm.DB, opts ...UserRepoOption) repo.UserRepository {
	r := &userRepo{
		db:     db,
		logger: zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *userRepo) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	r.logger.Debug().Msg("create user")

	tx := r.db.WithContext(ctx).Create(user)
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return user, repo.ErrDuplicateUser
	}

	return user, tx.Error
}

func (r *userRepo) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	r.logger.Debug().Msg("update user")

	tx := r.db.WithContext(ctx).Save(user)
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return user, repo.ErrDuplicateUser
	}

	return user, tx.Error
}

func (r *userRepo) FindUserByID(ctx context.Context, id uuid.UUID, scopes ...repo.Scope) (*entity.User, error) {
	r.logger.Debug().Msg("find user by id")

	tx := r.db.WithContext(ctx)

	for _, option := range scopes {
		tx = tx.Scopes(option)
	}

	user := &entity.User{}
	tx = tx.Scopes(notDeletedUsers).
		Where("users.id = ?", id).
		First(user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, tx.Error
}

func (r *userRepo) FindUserByUsername(ctx context.Context, username string, scopes ...repo.Scope) (
	*entity.User,
	error,
) {
	r.logger.Debug().Msg("find user by username")

	tx := r.db.WithContext(ctx)

	for _, option := range scopes {
		tx = tx.Scopes(option)
	}

	user := &entity.User{}
	tx = tx.Scopes(notDeletedUsers).
		Where("username = ?", username).
		First(user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, tx.Error
}

func (r *userRepo) FindUserByEmail(ctx context.Context, email string, scopes ...repo.Scope) (*entity.User, error) {
	r.logger.Debug().Msg("find user by email")

	tx := r.db.WithContext(ctx)

	for _, option := range scopes {
		tx = tx.Scopes(option)
	}

	user := &entity.User{}
	tx = tx.Scopes(notDeletedUsers).
		Where("email = ?", email).
		First(user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, tx.Error
}

func (r *userRepo) FindUserByUsernameOrEmail(ctx context.Context, login string, scopes ...repo.Scope) (
	*entity.User,
	error,
) {
	r.logger.Debug().Msg("find user by username or email")

	tx := r.db.WithContext(ctx)

	for _, option := range scopes {
		tx = tx.Scopes(option)
	}

	user := &entity.User{}
	tx = tx.Scopes(notDeletedUsers).
		Where("username = ?", login).
		Or("email = ?", login).
		First(user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, tx.Error
}

func (r *userRepo) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {
	r.logger.Debug().Msg("exists user by id")

	var exists bool
	tx := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists)
	return exists, tx.Error
}

func (r *userRepo) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	r.logger.Debug().Msg("exists user by username")

	var exists bool
	tx := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists)
	return exists, tx.Error
}

func (r *userRepo) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	r.logger.Debug().Msg("exists user by email")

	var exists bool
	tx := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists)
	return exists, tx.Error
}

func (r *userRepo) ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error) {
	r.logger.Debug().Msg("exists user by username or email")

	var exists bool
	tx := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("username = ?", username).
		Or("email = ?", email).
		Find(&exists)
	return exists, tx.Error
}

func (r *userRepo) DeleteExpiredUser(ctx context.Context) (*entity.User, error) {
	r.logger.Debug().Msg("delete expired user")

	user := &entity.User{}
	tx := r.db.WithContext(ctx).
		Scopes(deletedUsers).
		Delete(user)
	return user, tx.Error
}

func (r *userRepo) WithRolePreload() repo.Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.InnerJoins("Role")
	}
}

func notDeletedUsers(db *gorm.DB) *gorm.DB {
	return db.Where("(deleted_at is NULL OR now() - deleted_at <= INTERVAL '365 days')")
}

func deletedUsers(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at is NOT NULL").Where("now() - deleted_at > INTERVAL '365 days'")
}
