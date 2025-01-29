package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/rs/zerolog"
	"math"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Option func(*manager) error

func WithTTL(ttl time.Duration) Option {
	return func(m *manager) error {
		m.ttl = ttl
		return nil
	}
}

func WithLogger(logger zerolog.Logger) Option {
	return func(m *manager) error {
		m.logger = logger
		return nil
	}
}

type entry struct {
	value      any
	expiration int64
}

type manager struct {
	lock    sync.RWMutex
	storage map[string]entry
	ttl     time.Duration
	logger  zerolog.Logger
}

func NewManager(opts ...Option) (cache.Manager, error) {
	m := &manager{
		storage: make(map[string]entry),
		ttl:     5 * time.Minute,
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return m, nil
}

func (m *manager) Get(_ context.Context, key string, value any) error {
	m.cleanExpiredEntry()
	m.lock.RLock()
	defer m.lock.RUnlock()

	m.logger.Debug().Msgf("get from cache: %s", key)

	e, ok := m.storage[key]
	if !ok {
		return cache.ErrCacheEntryNotFound
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("value must be a non-nil pointer")
	}

	val.Elem().Set(reflect.ValueOf(e.value))
	return nil
}

func (m *manager) Set(ctx context.Context, key string, value any) error {
	return m.SetWithExpiration(ctx, key, value, m.ttl)
}

func (m *manager) SetWithExpiration(_ context.Context, key string, value any, expiration time.Duration) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	m.logger.Debug().Msgf("set to cache: %s with expiration %s", key, expiration)

	var expiredTime int64 = math.MaxInt64
	if expiration > 0 {
		expiredTime = time.Now().Add(expiration).Unix()
	}

	m.storage[key] = entry{
		value:      value,
		expiration: expiredTime,
	}
	return nil
}

func (m *manager) Delete(_ context.Context, keys ...string) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	m.logger.Debug().Msgf("delete from cache: %s", strings.Join(keys, ","))

	for _, key := range keys {
		delete(m.storage, key)
	}
	return nil
}

func (m *manager) Invalidate(_ context.Context, keyRegex string) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	m.logger.Debug().Msgf("invalidate cache by regex %s", keyRegex)

	for key := range m.storage {
		matched, err := regexp.MatchString(keyRegex, key)
		if err == nil && matched {
			delete(m.storage, key)
		}
	}

	return nil
}

func (m *manager) cleanExpiredEntry() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.logger.Debug().Msg("clean expired entry")

	now := time.Now().Unix()
	for key, entry := range m.storage {
		if entry.expiration <= now {
			delete(m.storage, key)
		}
	}
}
