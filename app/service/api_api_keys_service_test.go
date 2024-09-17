package service_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/lanyard/app/dal"
	"github.com/payloadops/lanyard/app/dal/mocks"
	"github.com/payloadops/lanyard/app/openapi"
	"github.com/payloadops/lanyard/app/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAPIKeysAPIService_DeleteApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockServiceClient := mocks.NewMockServiceManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockServiceClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	serviceID := "serv1"
	keyID := "key1"

	mockServiceClient.EXPECT().GetService(ctx, "org1", serviceID).Return(&dal.Service{}, nil)
	mockAPIKeyClient.EXPECT().GetAPIKey(ctx, keyID).Return(&dal.APIKey{ServiceID: serviceID}, nil)
	mockAPIKeyClient.EXPECT().DeleteAPIKey(ctx, "org1", serviceID, keyID).Return(nil)

	response, err := service.DeleteApiKey(ctx, serviceID, keyID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestAPIKeysAPIService_GenerateApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockServiceClient := mocks.NewMockServiceManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockServiceClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	serviceID := "serv1"
	apiKeyInput := openapi.ApiKeyInput{
		Scopes: []string{"scope1", "scope2"},
	}

	mockServiceClient.EXPECT().GetService(ctx, "org1", serviceID).Return(&dal.Service{}, nil)
	mockAPIKeyClient.EXPECT().CreateAPIKey(ctx, gomock.Any()).Return(nil)

	response, err := service.GenerateApiKey(ctx, serviceID, apiKeyInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	apiKey, ok := response.Body.(openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, serviceID, apiKey.ServiceId)
	assert.Equal(t, apiKeyInput.Scopes, apiKey.Scopes)
}

func TestAPIKeysAPIService_GetApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockServiceClient := mocks.NewMockServiceManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockServiceClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	serviceID := "serv1"
	keyID := "key1"

	mockServiceClient.EXPECT().GetService(ctx, "org1", serviceID).Return(&dal.Service{}, nil)
	mockAPIKeyClient.EXPECT().GetAPIKey(ctx, keyID).Return(&dal.APIKey{ServiceID: serviceID}, nil)

	response, err := service.GetApiKey(ctx, serviceID, keyID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	apiKey, ok := response.Body.(openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, serviceID, apiKey.ActorId)
}

func TestAPIKeysAPIService_ListApiKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockServiceClient := mocks.NewMockServiceManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockServiceClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	serviceID := "serv1"

	apiKeys := []dal.APIKey{
		{APIKeyID: "key1", ServiceID: serviceID},
		{APIKeyID: "key2", ServiceID: serviceID},
	}

	mockServiceClient.EXPECT().GetService(ctx, "org1", serviceID).Return(&dal.Service{}, nil)
	mockAPIKeyClient.EXPECT().ListAPIKeysByService(ctx, "org1", serviceID).Return(apiKeys, nil)

	response, err := service.ListApiKeys(ctx, serviceID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	keys, ok := response.Body.([]openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, "key1", keys[0].Id)
	assert.Equal(t, "key2", keys[1].Id)
}

func TestAPIKeysAPIService_UpdateApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockServiceClient := mocks.NewMockServiceManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockServiceClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	serviceID := "serv1"
	keyID := "key1"
	apiKeyInput := openapi.ApiKeyInput{
		Scopes: []string{"new-scope1", "new-scope2"},
	}

	apiKey := &dal.APIKey{
		APIKeyID:  keyID,
		ServiceID: serviceID,
		Scopes:    []string{"old-scope1", "old-scope2"},
	}

	mockServiceClient.EXPECT().GetService(ctx, "org1", serviceID).Return(&dal.Service{}, nil)
	mockAPIKeyClient.EXPECT().GetAPIKey(ctx, keyID).Return(apiKey, nil)
	mockAPIKeyClient.EXPECT().UpdateAPIKey(ctx, gomock.Any()).Return(nil)

	response, err := service.UpdateApiKey(ctx, serviceID, keyID, apiKeyInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	updatedKey, ok := response.Body.(openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, apiKeyInput.Scopes, updatedKey.Scopes)
	assert.Equal(t, serviceID, updatedKey.ActorId)
}
