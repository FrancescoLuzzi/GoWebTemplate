package cache

import (
	"context"
	"encoding"
	"errors"
	"time"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
)

var (
	ErrKeyNotFound = errors.New("key not found in cache")
)

type CacheValue interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type Cache interface {
	Set(ctx context.Context, key string, value CacheValue, ttl time.Duration) error
	Get(ctx context.Context, key string, value CacheValue) error
}

func NewCache(cfg config.CacheConfig) (Cache, error) {
	if cfg.Host == "" || cfg.Host == "memory" {
		return NewMemoryCache(400)
	}
	return NewRedisCache(cfg)
}
