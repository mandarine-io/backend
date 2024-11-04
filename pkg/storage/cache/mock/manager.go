package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) Get(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *ManagerMock) Set(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *ManagerMock) SetWithExpiration(
	ctx context.Context, key string, value interface{}, expiration time.Duration,
) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *ManagerMock) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *ManagerMock) Invalidate(ctx context.Context, keyRegex string) error {
	args := m.Called(ctx, keyRegex)
	return args.Error(0)
}
