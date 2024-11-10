package gorm

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/persistence/model"
	repo2 "github.com/mandarine-io/Backend/internal/persistence/repo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repo2.UserRepository {
	return &userRepository{db}
}

func (u *userRepository) CreateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error) {
	log.Debug().Msg("create user")
	tx := u.db.WithContext(ctx).Create(user)
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return user, repo2.ErrDuplicateUser
	}
	return user, tx.Error
}

func (u *userRepository) UpdateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error) {
	log.Debug().Msg("update user")
	tx := u.db.WithContext(ctx).Save(user)
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return user, repo2.ErrDuplicateUser
	}
	return user, tx.Error
}

func (u *userRepository) FindUserById(ctx context.Context, id uuid.UUID, roleJoin bool) (*model.UserEntity, error) {
	log.Debug().Msg("find user by id")
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).InnerJoins("Role")
	} else {
		db = u.db.WithContext(ctx)
	}

	userEntity := &model.UserEntity{}
	tx := db.Scopes(notDeletedUsers).
		Where("users.id = ?", id).
		First(userEntity)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return userEntity, tx.Error
}

func (u *userRepository) FindUserByUsername(ctx context.Context, username string, roleJoin bool) (*model.UserEntity, error) {
	log.Debug().Msg("find user by username")
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).InnerJoins("Role")
	} else {
		db = u.db.WithContext(ctx)
	}

	userEntity := &model.UserEntity{}
	tx := db.Scopes(notDeletedUsers).
		Where("username = ?", username).
		First(userEntity)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return userEntity, tx.Error
}

func (u *userRepository) FindUserByEmail(ctx context.Context, email string, roleJoin bool) (*model.UserEntity, error) {
	log.Debug().Msg("find user by email")
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).InnerJoins("Role")
	} else {
		db = u.db.WithContext(ctx)
	}

	userEntity := &model.UserEntity{}
	tx := db.Scopes(notDeletedUsers).
		Where("email = ?", email).
		First(userEntity)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return userEntity, tx.Error
}

func (u *userRepository) FindUserByUsernameOrEmail(ctx context.Context, login string, roleJoin bool) (*model.UserEntity, error) {
	log.Debug().Msg("find user by username or email")
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).InnerJoins("Role")
	} else {
		db = u.db.WithContext(ctx)
	}

	userEntity := &model.UserEntity{}
	tx := db.Scopes(notDeletedUsers).
		Where("username = ?", login).
		Or("email = ?", login).
		First(userEntity)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return userEntity, tx.Error
}

func (u *userRepository) ExistsUserById(ctx context.Context, id uuid.UUID) (bool, error) {
	log.Debug().Msg("exists user by id")
	var exists bool
	tx := u.db.WithContext(ctx).
		Model(&model.UserEntity{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists)
	return exists, tx.Error
}

func (u *userRepository) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	log.Debug().Msg("exists user by username")
	var exists bool
	tx := u.db.WithContext(ctx).
		Model(&model.UserEntity{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists)
	return exists, tx.Error
}

func (u *userRepository) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	log.Debug().Msg("exists user by email")
	var exists bool
	tx := u.db.WithContext(ctx).
		Model(&model.UserEntity{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists)
	return exists, tx.Error
}

func (u *userRepository) ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error) {
	log.Debug().Msg("exists user by username or email")
	var exists bool
	tx := u.db.WithContext(ctx).
		Model(&model.UserEntity{}).
		Scopes(notDeletedUsers).
		Select("count(*) > 0").
		Where("username = ?", username).
		Or("email = ?", email).
		Find(&exists)
	return exists, tx.Error
}

func (u *userRepository) DeleteExpiredUser(ctx context.Context) (*model.UserEntity, error) {
	log.Debug().Msg("delete expired user")
	userEntity := &model.UserEntity{}
	tx := u.db.WithContext(ctx).
		Scopes(deletedUsers).
		Delete(userEntity)
	return userEntity, tx.Error
}

func notDeletedUsers(db *gorm.DB) *gorm.DB {
	return db.Where("(deleted_at is NULL OR now() - deleted_at <= INTERVAL '365 days')")
}

func deletedUsers(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at is NOT NULL").Where("now() - deleted_at > INTERVAL '365 days'")
}
