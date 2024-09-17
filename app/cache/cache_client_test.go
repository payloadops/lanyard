package cache_test

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/go-redis/redis/v8"
	"github.com/payloadops/lanyard/app/cache"
	"github.com/payloadops/lanyard/app/cache/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := mocks.NewMockCmdable(ctrl)
	redisCache := cache.NewRedisCache(mockRedisClient)

	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := 10 * time.Second

	mockRedisClient.EXPECT().Set(ctx, key, value, expiration).Return(redis.NewStatusResult("OK", nil))

	err := redisCache.Set(ctx, key, value, expiration)
	assert.NoError(t, err)
}

func TestRedisCache_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := mocks.NewMockCmdable(ctrl)
	redisCache := cache.NewRedisCache(mockRedisClient)

	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := 10 * time.Second

	script := `
		local value = redis.call('GET', KEYS[1])
		if value then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return value
	`

	mockRedisClient.EXPECT().
		Eval(ctx, script, []string{key}, int(expiration.Seconds())).
		Return(redis.NewCmdResult(value, nil))

	result, err := redisCache.Get(ctx, key, expiration)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestRedisCache_Get_KeyNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := mocks.NewMockCmdable(ctrl)
	redisCache := cache.NewRedisCache(mockRedisClient)

	ctx := context.Background()
	key := "test-key"
	expiration := 10 * time.Second

	script := `
		local value = redis.call('GET', KEYS[1])
		if value then
			redis.call('EXPIRE', KEYS[1], ARGV[1])
		end
		return value
	`

	mockRedisClient.EXPECT().
		Eval(ctx, script, []string{key}, int(expiration.Seconds())).
		Return(redis.NewCmdResult(nil, redis.Nil))

	result, err := redisCache.Get(ctx, key, expiration)
	assert.ErrorIs(t, err, redis.Nil)
	assert.Equal(t, "", result)
}

func TestNoopCache_Set(t *testing.T) {
	noopCache := cache.NewNoopCache()
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := 10 * time.Second

	err := noopCache.Set(ctx, key, value, expiration)
	assert.NoError(t, err)
}

func TestNoopCache_Get(t *testing.T) {
	noopCache := cache.NewNoopCache()
	ctx := context.Background()
	key := "test-key"
	expiration := 10 * time.Second

	result, err := noopCache.Get(ctx, key, expiration)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}
