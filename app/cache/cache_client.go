package cache

/*
import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_cache_client.go "github.com/payloadops/plato/api/cache" Cache

// Cache is an interface defining methods for a caching layer.
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
}

// RedisCache implements the Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Set stores a value in the cache.
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value from the cache.
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}
*/
