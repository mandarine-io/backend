package mock

import (
	"context"
	"github.com/google/uuid"
	"github.com/mandarine-io/Backend/internal/api/persistence/model"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user1 *model.UserEntity) (*model.UserEntity, error) {
	args := m.Called(ctx, user1)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *UserRepositoryMock) UpdateUser(ctx context.Context, user1 *model.UserEntity) (*model.UserEntity, error) {
	args := m.Called(ctx, user1)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *UserRepositoryMock) FindUserById(ctx context.Context, id uuid.UUID, rolePreload bool) (*model.UserEntity, error) {
	args := m.Called(ctx, id, rolePreload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *UserRepositoryMock) FindUserByEmail(ctx context.Context, email string, rolePreload bool) (*model.UserEntity, error) {
	args := m.Called(ctx, email, rolePreload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *UserRepositoryMock) FindUserByUsername(ctx context.Context, username string, rolePreload bool) (*model.UserEntity, error) {
	args := m.Called(ctx, username, rolePreload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *UserRepositoryMock) FindUserByUsernameOrEmail(ctx context.Context, login string, rolePreload bool) (*model.UserEntity, error) {
	args := m.Called(ctx, login, rolePreload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (m *UserRepositoryMock) ExistsUserById(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepositoryMock) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepositoryMock) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(bool), args.Error(1)
}

func (m *UserRepositoryMock) ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error) {
	args := m.Called(ctx, username, email)
	return args.Get(0).(bool), args.Error(1)
}

func (m *UserRepositoryMock) DeleteExpiredUser(ctx context.Context) (*model.UserEntity, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.UserEntity), args.Error(1)
}
