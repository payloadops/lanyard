package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/api/utils"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

const (
	// ApiKeyLength represents the length of an API key to generate
	ApiKeyLength = 32
)

// APIKeysAPIService is a service that implements the logic for the APIKeysAPIServicer
// This service should implement the business logic for every endpoint for the APIKeysAPI API.
type APIKeysAPIService struct {
	apiKeyClient  dal.APIKeyManager
	projectClient dal.ProjectManager
}

// NewAPIKeysAPIService creates a default api service
func NewAPIKeysAPIService() openapi.APIKeysAPIServicer {
	apiKeyClient, err := dal.NewAPIKeyDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create APIKeyDBClient: %v", err))
	}
	projectClient, err := dal.NewProjectDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create ProjectDBClient: %v", err))
	}
	return &APIKeysAPIService{apiKeyClient: apiKeyClient, projectClient: projectClient}
}

// DeleteApiKey - Delete a specific API key
func (s *APIKeysAPIService) DeleteApiKey(ctx context.Context, projectId string, keyId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the API key exists
	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if apiKey == nil || apiKey.ProjectID != projectId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("API key not found")
	}

	err = s.apiKeyClient.DeleteAPIKey(ctx, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GenerateApiKey - Generate a new API key with specific scopes for a project
func (s *APIKeysAPIService) GenerateApiKey(ctx context.Context, projectId string, apiKeyInput openapi.ApiKeyInput) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	keySecret, err := utils.GenerateSecret(ApiKeyLength)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	apiKey := dal.APIKey{
		ID:        ksuid,
		ProjectID: projectId,
		Key:       keySecret,
		Scopes:    apiKeyInput.Scopes,
	}

	err = s.apiKeyClient.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, apiKey), nil
}

// GetApiKey - Retrieve a specific API key
func (s *APIKeysAPIService) GetApiKey(ctx context.Context, projectId string, keyId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if apiKey == nil || apiKey.ProjectID != projectId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("API key not found")
	}

	return openapi.Response(http.StatusOK, apiKey), nil
}

// ListApiKeys - List all API keys for a project
func (s *APIKeysAPIService) ListApiKeys(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	apiKeys, err := s.apiKeyClient.ListAPIKeys(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, apiKeys), nil
}

// UpdateApiKey - Update an API key's scopes
func (s *APIKeysAPIService) UpdateApiKey(ctx context.Context, projectId string, keyId string, apiKeyInput openapi.ApiKeyInput) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the API key exists
	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if apiKey == nil || apiKey.ProjectID != projectId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("API key not found")
	}

	// Update the API key with the new values
	apiKey.Scopes = apiKeyInput.Scopes
	apiKey.ProjectID = projectId

	err = s.apiKeyClient.UpdateAPIKey(ctx, *apiKey)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, apiKey), nil
}
