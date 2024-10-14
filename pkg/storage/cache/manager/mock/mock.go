package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type CacheManagerMock struct {
	mock.Mock
}

func (m *CacheManagerMock) Get(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *CacheManagerMock) Set(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *CacheManagerMock) SetWithExpiration(
	ctx context.Context, key string, value interface{}, expiration time.Duration,
) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *CacheManagerMock) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *CacheManagerMock) Invalidate(ctx context.Context, keyRegex string) error {
	args := m.Called(ctx, keyRegex)
	return args.Error(0)
}
