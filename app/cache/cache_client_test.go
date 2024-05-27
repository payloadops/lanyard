package cache

/*
import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	// Setup the Redis client for testing (make sure Redis is running locally)
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Default Redis address
	})
}

func TestRedisCache_SetAndGet(t *testing.T) {
	client := setupTestRedis()
	cache := NewRedisCache(client)

	type TestData struct {
		Value string `json:"value"`
	}

	ctx := context.Background()
	key := "testKey"
	value := TestData{Value: "testValue"}

	// Set the value in the cache
	err := cache.Set(ctx, key, value, 1*time.Hour)
	assert.NoError(t, err)

	// Get the value from the cache
	var result TestData
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestRedisCache_GetNonExistentKey(t *testing.T) {
	client := setupTestRedis()
	cache := NewRedisCache(client)

	ctx := context.Background()
	key := "nonExistentKey"

	var result interface{}
	err := cache.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisCache_SetWithExpiration(t *testing.T) {
	client := setupTestRedis()
	cache := NewRedisCache(client)

	type TestData struct {
		Value string `json:"value"`
	}

	ctx := context.Background()
	key := "expiringKey"
	value := TestData{Value: "expiringValue"}

	// Set the value in the cache with a short expiration
	err := cache.Set(ctx, key, value, 1*time.Second)
	assert.NoError(t, err)

	// Get the value from the cache immediately
	var result TestData
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Wait for the expiration time to pass
	time.Sleep(2 * time.Second)

	// Try to get the value from the cache again
	err = cache.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisCache_SetInvalidValue(t *testing.T) {
	client := setupTestRedis()
	cache := NewRedisCache(client)

	ctx := context.Background()
	key := "invalidValueKey"
	value := make(chan int) // Invalid value type for JSON marshalling

	err := cache.Set(ctx, key, value, 1*time.Hour)
	assert.Error(t, err)
}

func TestRedisCache_GetInvalidValue(t *testing.T) {
	client := setupTestRedis()
	cache := NewRedisCache(client)

	type TestData struct {
		Value string `json:"value"`
	}

	ctx := context.Background()
	key := "invalidGetKey"
	value := TestData{Value: "invalidGetValue"}

	// Set a valid value in the cache
	err := cache.Set(ctx, key, value, 1*time.Hour)
	assert.NoError(t, err)

	// Try to get the value into an invalid type
	var result chan int
	err = cache.Get(ctx, key, &result)
	assert.Error(t, err)
}
*/
