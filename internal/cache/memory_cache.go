// Package cache
package cache

import (
	"sync"
	"time"

	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/utils"
)

func isOutDate[T any](data *CachedItem[T]) bool {
	return data.ExpiredAt.Before(time.Now())
}

type MemoryCache[T any] struct {
	cacheMap map[string]*CachedItem[T]
	cleaner  *utils.IntervalActuator
	lock     sync.RWMutex
}

func NewMemoryCache[T any](cleanInterval time.Duration) *MemoryCache[T] {
	if cleanInterval <= 0 {
		cleanInterval = 30 * time.Minute
	}
	cached := &MemoryCache[T]{
		cacheMap: make(map[string]*CachedItem[T]),
		lock:     sync.RWMutex{},
	}
	cached.cleaner = utils.NewIntervalActuator(cleanInterval, cached.CleanExpiredData)
	return cached
}

func (cache *MemoryCache[T]) CleanExpiredData() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for key, value := range cache.cacheMap {
		if isOutDate(value) {
			delete(cache.cacheMap, key)
		}
	}
}

func (cache *MemoryCache[T]) Set(key string, value T, expiredAt time.Time) {
	if expiredAt.Before(time.Now()) {
		return
	}
	if key == "" {
		return
	}
	cache.lock.Lock()
	cache.cacheMap[key] = &CachedItem[T]{CachedData: value, ExpiredAt: expiredAt}
	cache.lock.Unlock()
}

func (cache *MemoryCache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	expiredAt := time.Now().Add(ttl)
	cache.Set(key, value, expiredAt)
}

func (cache *MemoryCache[T]) Get(key string) (T, bool) {
	if key == "" {
		var zero T
		return zero, false
	}
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	val, ok := cache.cacheMap[key]
	if ok && isOutDate(val) {
		var zero T
		return zero, false
	}
	if val == nil {
		var zero T
		return zero, false
	}
	return val.CachedData, ok
}

func (cache *MemoryCache[T]) Del(key string) {
	if key == "" {
		return
	}
	cache.lock.Lock()
	delete(cache.cacheMap, key)
	cache.lock.Unlock()
}

func (cache *MemoryCache[T]) Close() {
	cache.cleaner.Stop()
}
