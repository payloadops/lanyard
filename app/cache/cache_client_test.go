package cache_test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/payloadops/plato/app/cache"
	"testing"
	"time"

	"github.com/payloadops/plato/app/cache/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCmdable(ctrl)
	cache := cache.NewRedisCache(mockClient)

	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	expiration := time.Hour

	mockClient.EXPECT().Set(ctx, key, value, expiration).Return(&redis.StatusCmd{})

	err := cache.Set(ctx, key, value, expiration)
	assert.NoError(t, err)
}

func TestRedisCache_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCmdable(ctrl)
	cache := cache.NewRedisCache(mockClient)

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	mockClient.EXPECT().Get(ctx, key).Return(redis.NewStringResult(value, nil))

	result, err := cache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestRedisCache_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockCmdable(ctrl)
	cache := cache.NewRedisCache(mockClient)

	ctx := context.Background()
	key := "non-existent-key"

	mockClient.EXPECT().Get(ctx, key).Return(redis.NewStringResult("", redis.Nil))

	_, err := cache.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestNoopCache_Set(t *testing.T) {
	nc := cache.NewNoopCache()

	err := nc.Set(context.Background(), "key", "value", 10*time.Minute)
	assert.NoError(t, err, "Expected no error for NoopCache Set")
}

func TestNoopCache_Get(t *testing.T) {
	nc := cache.NewNoopCache()

	value, err := nc.Get(context.Background(), "key")
	assert.NoError(t, err, "Expected no error for NoopCache Get")
	assert.Equal(t, "", value, "Expected empty string value for NoopCache Get")
}
