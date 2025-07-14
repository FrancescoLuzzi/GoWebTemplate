package cache

import (
	"context"
	"errors"
	"time"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/config"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisNotConnected = errors.New("redis not connected")
)

type redisClient struct {
	client *redis.Client
}

var _ Cache = &redisClient{}

func (r *redisClient) Get(ctx context.Context, key string, value CacheValue) error {
	val, err := r.client.Get(ctx, key).Bytes()
	switch err {
	case nil:
		return value.UnmarshalBinary(val)
	default:
		return ErrKeyNotFound
	}
}

func (r *redisClient) Set(ctx context.Context, key string, value CacheValue, ttl time.Duration) error {
	val, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = r.client.Set(ctx, key, val, ttl).Result()
	return err
}

func redisOptionsFromConfig(cfg config.CacheConfig) *redis.Options {
	return &redis.Options{
		Addr:        cfg.Host + ":" + cfg.Port,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DB:          cfg.DB,
		PoolSize:    cfg.PoolSize,
		PoolTimeout: 30 * time.Second,
	}
}

func NewRedisCache(cfg config.CacheConfig) (*redisClient, error) {
	rdb := redis.NewClient(redisOptionsFromConfig(cfg))
	status, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	if status != "PONG" {
		return nil, ErrRedisNotConnected
	}
	return &redisClient{client: rdb}, nil
}
