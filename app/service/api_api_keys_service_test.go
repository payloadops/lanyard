package service

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPIKeyManager is a mock implementation of the APIKeyManager interface
type MockAPIKeyManager struct {
	mock.Mock
}

func (m *MockAPIKeyManager) CreateAPIKey(ctx context.Context, apiKey dal.APIKey) error {
	args := m.Called(ctx, apiKey)
	return args.Error(0)
}

func (m *MockAPIKeyManager) GetAPIKey(ctx context.Context, id string) (*dal.APIKey, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.APIKey), args.Error(1)
}

func (m *MockAPIKeyManager) UpdateAPIKey(ctx context.Context, apiKey dal.APIKey) error {
	args := m.Called(ctx, apiKey)
	return args.Error(0)
}

func (m *MockAPIKeyManager) DeleteAPIKey(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPIKeyManager) ListAPIKeys(ctx context.Context, projectId string) ([]dal.APIKey, error) {
	args := m.Called(ctx, projectId)
	return args.Get(0).([]dal.APIKey), args.Error(1)
}

func TestDeleteApiKey(t *testing.T) {
	mockAPIKeyClient := new(MockAPIKeyManager)
	mockProjectClient := new(MockProjectManager)
	service := APIKeysAPIService{apiKeyClient: mockAPIKeyClient, projectClient: mockProjectClient}

	projectId := "project1"
	keyId := "key1"

	// Test case where project and API key exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("GetAPIKey", mock.Anything, keyId).Return(&dal.APIKey{ID: keyId, ProjectID: projectId}, nil)
	mockAPIKeyClient.On("DeleteAPIKey", mock.Anything, keyId).Return(nil)

	resp, err := service.DeleteApiKey(context.Background(), projectId, keyId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.DeleteApiKey(context.Background(), projectId, keyId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where API key does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("GetAPIKey", mock.Anything, keyId).Return((*dal.APIKey)(nil), nil)

	resp, err = service.DeleteApiKey(context.Background(), projectId, keyId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)
}

func TestGenerateApiKey(t *testing.T) {
	mockAPIKeyClient := new(MockAPIKeyManager)
	mockProjectClient := new(MockProjectManager)
	service := APIKeysAPIService{apiKeyClient: mockAPIKeyClient, projectClient: mockProjectClient}

	projectId := "project1"
	apiKeyInput := openapi.ApiKeyInput{
		Scopes: []string{"read", "write"},
	}

	expectedAPIKey := dal.APIKey{
		ID:        "id",
		ProjectID: projectId,
		Key:       "key",
		Scopes:    apiKeyInput.Scopes,
	}

	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("CreateAPIKey", mock.Anything, expectedAPIKey).Return(nil)

	resp, err := service.GenerateApiKey(context.Background(), projectId, apiKeyInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)
}

func TestGetApiKey(t *testing.T) {
	mockAPIKeyClient := new(MockAPIKeyManager)
	mockProjectClient := new(MockProjectManager)
	service := APIKeysAPIService{apiKeyClient: mockAPIKeyClient, projectClient: mockProjectClient}

	projectId := "project1"
	keyId := "key1"
	apiKey := &dal.APIKey{
		ID:        keyId,
		ProjectID: projectId,
		Key:       "secret-key",
		Scopes:    []string{"read", "write"},
	}

	// Test case where project and API key exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("GetAPIKey", mock.Anything, keyId).Return(apiKey, nil)

	resp, err := service.GetApiKey(context.Background(), projectId, keyId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.GetApiKey(context.Background(), projectId, keyId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where API key does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("GetAPIKey", mock.Anything, keyId).Return((*dal.APIKey)(nil), nil)

	resp, err = service.GetApiKey(context.Background(), projectId, keyId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)
}

func TestListApiKeys(t *testing.T) {
	mockAPIKeyClient := new(MockAPIKeyManager)
	mockProjectClient := new(MockProjectManager)
	service := APIKeysAPIService{apiKeyClient: mockAPIKeyClient, projectClient: mockProjectClient}

	projectId := "project1"
	apiKeys := []dal.APIKey{
		{ID: "key1", ProjectID: projectId, Key: "key1-secret", Scopes: []string{"read"}},
		{ID: "key2", ProjectID: projectId, Key: "key2-secret", Scopes: []string{"write"}},
	}

	// Test case where project exists
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("ListAPIKeys", mock.Anything, projectId).Return(apiKeys, nil)

	resp, err := service.ListApiKeys(context.Background(), projectId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.ListApiKeys(context.Background(), projectId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
}

func TestUpdateApiKey(t *testing.T) {
	mockAPIKeyClient := new(MockAPIKeyManager)
	mockProjectClient := new(MockProjectManager)
	service := APIKeysAPIService{apiKeyClient: mockAPIKeyClient, projectClient: mockProjectClient}

	projectId := "project1"
	keyId := "key1"
	apiKeyInput := openapi.ApiKeyInput{
		Scopes: []string{"read", "write"},
	}
	apiKey := &dal.APIKey{
		ID:        keyId,
		ProjectID: projectId,
		Key:       "key",
		Scopes:    apiKeyInput.Scopes,
	}

	// Test case where project and API key exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("GetAPIKey", mock.Anything, keyId).Return(&dal.APIKey{ID: keyId, ProjectID: projectId}, nil)
	mockAPIKeyClient.On("UpdateAPIKey", mock.Anything, *apiKey).Return(nil)

	resp, err := service.UpdateApiKey(context.Background(), projectId, keyId, apiKeyInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.UpdateApiKey(context.Background(), projectId, keyId, apiKeyInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where API key does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockAPIKeyClient.On("GetAPIKey", mock.Anything, keyId).Return((*dal.APIKey)(nil), nil)

	resp, err = service.UpdateApiKey(context.Background(), projectId, keyId, apiKeyInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockAPIKeyClient.AssertExpectations(t)
}
*/
