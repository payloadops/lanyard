package service

import (
	"context"
	"errors"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
	"github.com/payloadops/plato/app/utils"
	"go.uber.org/zap"
	"net/http"
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
	logger        zap.Logger
}

// NewAPIKeysAPIService creates a default app service
func NewAPIKeysAPIService(apiKeyClient dal.APIKeyManager, projectClient dal.ProjectManager) openapi.APIKeysAPIServicer {
	return &APIKeysAPIService{apiKeyClient: apiKeyClient, projectClient: projectClient}
}

// DeleteApiKey - Delete a specific API key
func (s *APIKeysAPIService) DeleteApiKey(ctx context.Context, projectId string, keyId string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	// Check if the API key exists
	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, orgID, projectId, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if apiKey == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("API key not found")
	}

	err = s.apiKeyClient.DeleteAPIKey(ctx, orgID, projectId, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GenerateApiKey - Generate a new API key with specific scopes for a project
func (s *APIKeysAPIService) GenerateApiKey(ctx context.Context, projectId string, apiKeyInput openapi.ApiKeyInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	keySecret, err := utils.GenerateSecret(ApiKeyLength)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	apiKey := dal.APIKey{
		ProjectID: projectId,
		Secret:    keySecret,
		Scopes:    apiKeyInput.Scopes,
	}

	err = s.apiKeyClient.CreateAPIKey(ctx, orgID, &apiKey)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.ApiKey{
		Id:        apiKey.APIKeyID,
		Secret:    apiKey.Secret,
		Scopes:    apiKey.Scopes,
		ProjectId: apiKey.ProjectID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return openapi.Response(http.StatusCreated, response), nil
}

// GetApiKey - Retrieve a specific API key
func (s *APIKeysAPIService) GetApiKey(ctx context.Context, projectId string, keyId string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, orgID, projectId, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if apiKey == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("API key not found")
	}

	createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.ApiKey{
		Id:        apiKey.APIKeyID,
		Secret:    apiKey.Secret,
		Scopes:    apiKey.Scopes,
		ProjectId: apiKey.ProjectID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}

// ListApiKeys - List all API keys for a project
func (s *APIKeysAPIService) ListApiKeys(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	apiKeys, err := s.apiKeyClient.ListAPIKeysByProject(ctx, orgID, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	responses := make([]openapi.ApiKey, len(apiKeys))
	for i, apiKey := range apiKeys {
		createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		responses[i] = openapi.ApiKey{
			Id:        apiKey.APIKeyID,
			Secret:    apiKey.Secret,
			Scopes:    apiKey.Scopes,
			ProjectId: apiKey.ProjectID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return openapi.Response(http.StatusOK, responses), nil
}

// UpdateApiKey - Update an API key's scopes
func (s *APIKeysAPIService) UpdateApiKey(ctx context.Context, projectId string, keyId string, apiKeyInput openapi.ApiKeyInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	// Check if the API key exists
	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, orgID, projectId, keyId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if apiKey == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("API key not found")
	}

	// Update the API key with the new values
	apiKey.Scopes = apiKeyInput.Scopes
	err = s.apiKeyClient.UpdateAPIKey(ctx, orgID, apiKey)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.ApiKey{
		Id:        apiKey.APIKeyID,
		Secret:    apiKey.Secret,
		Scopes:    apiKey.Scopes,
		ProjectId: apiKey.ProjectID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}
