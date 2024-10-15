package gorm

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mandarine/internal/api/persistence/model"
	"mandarine/internal/api/persistence/repo"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repo.UserRepository {
	return &userRepository{db}
}

func (u *userRepository) CreateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error) {
	tx := u.db.WithContext(ctx).Create(user)
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return user, repo.ErrDuplicateUser
	}
	return user, tx.Error
}

func (u *userRepository) UpdateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error) {
	tx := u.db.WithContext(ctx).Save(user)
	if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
		return user, repo.ErrDuplicateUser
	}
	return user, tx.Error
}

func (u *userRepository) FindUserById(ctx context.Context, id uuid.UUID, roleJoin bool) (*model.UserEntity, error) {
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).Joins("Role")
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
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).Joins("Role")
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
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).Joins("Role")
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
	var db *gorm.DB
	if roleJoin {
		db = u.db.WithContext(ctx).Joins("Role")
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
