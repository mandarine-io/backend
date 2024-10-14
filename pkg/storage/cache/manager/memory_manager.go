package manager

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"sync"
	"time"
)

type cacheEntry struct {
	value      interface{}
	expiration int64
}

type MemoryCacheManager struct {
	lock    sync.RWMutex
	storage map[string]cacheEntry
	ttl     time.Duration
}

func NewMemoryCacheManager(ttl time.Duration) CacheManager {
	return &MemoryCacheManager{
		storage: make(map[string]cacheEntry),
		ttl:     ttl,
	}
}

func (m *MemoryCacheManager) Get(_ context.Context, key string, value interface{}) error {
	m.cleanExpiredEntry()
	m.lock.RLock()
	defer m.lock.RUnlock()

	entry, ok := m.storage[key]
	if !ok {
		return ErrCacheEntryNotFound
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("value must be a non-nil pointer")
	}

	val.Elem().Set(reflect.ValueOf(entry.value))
	return nil
}

func (m *MemoryCacheManager) Set(_ context.Context, key string, value interface{}) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	m.storage[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(m.ttl).Unix(),
	}
	return nil
}

func (m *MemoryCacheManager) SetWithExpiration(
	_ context.Context, key string, value interface{}, expiration time.Duration,
) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	m.storage[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(expiration).Unix(),
	}
	return nil
}

func (m *MemoryCacheManager) Delete(_ context.Context, keys ...string) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, key := range keys {
		delete(m.storage, key)
	}
	return nil
}

func (m *MemoryCacheManager) Invalidate(ctx context.Context, keyRegex string) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	for key := range m.storage {
		matched, err := regexp.MatchString(keyRegex, key)
		if err == nil && matched {
			delete(m.storage, key)
		}
	}

	return nil
}

func (m *MemoryCacheManager) cleanExpiredEntry() {
	m.lock.Lock()
	defer m.lock.Unlock()

	now := time.Now().Unix()
	for key, entry := range m.storage {
		if entry.expiration <= now {
			delete(m.storage, key)
		}
	}
}
