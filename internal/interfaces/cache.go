// Package interfaces
package interfaces

import "time"

type CachedItem[T any] struct {
	CachedData T
	ExpiredAt  time.Time
}

type CacheInterface[T any] interface {
	Set(key string, value T, expiredAt time.Time)
	SetWithTTL(key string, value T, ttl time.Duration)
	Get(key string) (T, bool)
	Del(key string)
	Close()
}
