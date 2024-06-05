package service_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/payloadops/plato/app/openapi"
	"github.com/payloadops/plato/app/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAPIKeysAPIService_DeleteApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockProjectClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	keyID := "key1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockAPIKeyClient.EXPECT().GetAPIKey(ctx, keyID).Return(&dal.APIKey{ProjectID: projectID}, nil)
	mockAPIKeyClient.EXPECT().DeleteAPIKey(ctx, "org1", projectID, keyID).Return(nil)

	response, err := service.DeleteApiKey(ctx, projectID, keyID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestAPIKeysAPIService_GenerateApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockProjectClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	apiKeyInput := openapi.ApiKeyInput{
		Scopes: []string{"scope1", "scope2"},
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockAPIKeyClient.EXPECT().CreateAPIKey(ctx, gomock.Any()).Return(nil)

	response, err := service.GenerateApiKey(ctx, projectID, apiKeyInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	apiKey, ok := response.Body.(openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, projectID, apiKey.ProjectId)
	assert.Equal(t, apiKeyInput.Scopes, apiKey.Scopes)
}

func TestAPIKeysAPIService_GetApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockProjectClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	keyID := "key1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockAPIKeyClient.EXPECT().GetAPIKey(ctx, keyID).Return(&dal.APIKey{ProjectID: projectID}, nil)

	response, err := service.GetApiKey(ctx, projectID, keyID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	apiKey, ok := response.Body.(openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, projectID, apiKey.ProjectId)
}

func TestAPIKeysAPIService_ListApiKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAPIKeyClient := mocks.NewMockAPIKeyManager(ctrl)
	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockProjectClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"

	apiKeys := []dal.APIKey{
		{APIKeyID: "key1", ProjectID: projectID},
		{APIKeyID: "key2", ProjectID: projectID},
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockAPIKeyClient.EXPECT().ListAPIKeysByProject(ctx, "org1", projectID).Return(apiKeys, nil)

	response, err := service.ListApiKeys(ctx, projectID)
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
	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	service := service.NewAPIKeysAPIService(mockAPIKeyClient, mockProjectClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	keyID := "key1"
	apiKeyInput := openapi.ApiKeyInput{
		Scopes: []string{"new-scope1", "new-scope2"},
	}

	apiKey := &dal.APIKey{
		APIKeyID:  keyID,
		ProjectID: projectID,
		Scopes:    []string{"old-scope1", "old-scope2"},
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockAPIKeyClient.EXPECT().GetAPIKey(ctx, keyID).Return(apiKey, nil)
	mockAPIKeyClient.EXPECT().UpdateAPIKey(ctx, gomock.Any()).Return(nil)

	response, err := service.UpdateApiKey(ctx, projectID, keyID, apiKeyInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	updatedKey, ok := response.Body.(openapi.ApiKey)
	assert.True(t, ok)
	assert.Equal(t, apiKeyInput.Scopes, updatedKey.Scopes)
	assert.Equal(t, projectID, updatedKey.ProjectId)
}
