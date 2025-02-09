package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_redis_client.go github.com/go-redis/redis/v8 Cmdable

//go:generate mockgen -package=mocks -destination=mocks/mock_cache_client.go github.com/payloadops/plato/app/cache Cache

// Cache is an interface defining methods for a caching layer.
type Cache interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string, expiration time.Duration) (string, error)
}

// Ensure RedisCache implements the Cache interface
var _ Cache = &RedisCache{}

// Ensure NoopCache implements the Cache interface
var _ Cache = &NoopCache{}

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

// Get retrieves a value from the cache and resets the expiration atomically.
func (r *RedisCache) Get(ctx context.Context, key string, expiration time.Duration) (string, error) {
	script := `
		local value = redis.call('GET', KEYS[1])
		if value then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return value
	`

	result, err := r.client.Eval(ctx, script, []string{key}, int(expiration.Seconds())).Result()
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", redis.Nil
	}

	return result.(string), nil
}

// NoopCache implements the Cache interface as a no-op.
type NoopCache struct{}

// NewNoopCache creates a new NoopCache.
func NewNoopCache() *NoopCache {
	return &NoopCache{}
}

// Set is a no-op for NoopCache.
func (n *NoopCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	// No operation performed
	return nil
}

// Get is a no-op for NoopCache.
func (n *NoopCache) Get(ctx context.Context, key string, expiration time.Duration) (string, error) {
	// No operation performed, return an empty string and no error
	return "", nil
}
