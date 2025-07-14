package cache

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type cacheContent struct {
	expireAt time.Time
	content  []byte
}

type memoryCache struct {
	cache *lru.Cache[string, cacheContent]
}

var _ Cache = &memoryCache{}

// Get implements Cache.
func (m *memoryCache) Get(ctx context.Context, key string, value CacheValue) error {
	val, ok := m.cache.Get(key)
	if !ok {
		return ErrKeyNotFound
	}
	if val.expireAt.Before(time.Now()) {
		m.cache.Remove(key)
		return ErrKeyNotFound
	}
	return value.UnmarshalBinary(val.content)

}

// Set implements Cache.
func (m *memoryCache) Set(ctx context.Context, key string, value CacheValue, ttl time.Duration) error {
	val, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	// Cache the result.
	m.cache.Add(key, cacheContent{
		expireAt: time.Now().Add(ttl),
		content:  val,
	})
	return nil
}

func NewMemoryCache(size int) (*memoryCache, error) {
	cache, err := lru.New[string, cacheContent](size)
	if err != nil {
		return nil, err
	}
	return &memoryCache{cache: cache}, nil
}
