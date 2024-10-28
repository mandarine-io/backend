package memory

import (
	"context"
	"errors"
	"github.com/mandarine-io/Backend/pkg/storage/cache/manager"
	"github.com/rs/zerolog/log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type cacheEntry struct {
	value      interface{}
	expiration int64
}

type cacheManager struct {
	lock    sync.RWMutex
	storage map[string]cacheEntry
	ttl     time.Duration
}

func NewCacheManager(ttl time.Duration) manager.CacheManager {
	return &cacheManager{
		storage: make(map[string]cacheEntry),
		ttl:     ttl,
	}
}

func (m *cacheManager) Get(_ context.Context, key string, value interface{}) error {
	m.cleanExpiredEntry()
	m.lock.RLock()
	defer m.lock.RUnlock()

	log.Debug().Msgf("get from cache: %s", key)

	entry, ok := m.storage[key]
	if !ok {
		return manager.ErrCacheEntryNotFound
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("value must be a non-nil pointer")
	}

	val.Elem().Set(reflect.ValueOf(entry.value))
	return nil
}

func (m *cacheManager) Set(_ context.Context, key string, value interface{}) error {
	return m.SetWithExpiration(context.Background(), key, value, m.ttl)
}

func (m *cacheManager) SetWithExpiration(
	_ context.Context, key string, value interface{}, expiration time.Duration,
) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	log.Debug().Msgf("set to cache: %s", key)

	m.storage[key] = cacheEntry{
		value:      value,
		expiration: time.Now().Add(expiration).Unix(),
	}
	return nil
}

func (m *cacheManager) Delete(_ context.Context, keys ...string) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	log.Debug().Msgf("delete from cache: %s", strings.Join(keys, ","))

	for _, key := range keys {
		delete(m.storage, key)
	}
	return nil
}

func (m *cacheManager) Invalidate(ctx context.Context, keyRegex string) error {
	m.cleanExpiredEntry()
	m.lock.Lock()
	defer m.lock.Unlock()

	log.Debug().Msgf("invalidate cache by regex %s", keyRegex)

	for key := range m.storage {
		matched, err := regexp.MatchString(keyRegex, key)
		if err == nil && matched {
			delete(m.storage, key)
		}
	}

	return nil
}

func (m *cacheManager) cleanExpiredEntry() {
	m.lock.Lock()
	defer m.lock.Unlock()

	log.Debug().Msg("clean expired entry")

	now := time.Now().Unix()
	for key, entry := range m.storage {
		if entry.expiration <= now {
			delete(m.storage, key)
		}
	}
}
