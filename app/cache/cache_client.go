package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_redis_client.go github.com/go-redis/redis/v8 Cmdable

//go:generate mockgen -package=mocks -destination=mocks/mock_cache_client.go "github.com/payloadops/plato/app/cache" Cache

// Cache is an interface defining methods for a caching layer.
type Cache interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

// RedisCache implements the Cache interface using Redis.
type RedisCache struct {
	client redis.Cmdable
}

// NewRedisCache creates a new RedisCache.
func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{client: client}
}

// Set stores a value in the cache.
func (r *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from the cache.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
