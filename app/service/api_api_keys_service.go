package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
	"github.com/payloadops/plato/app/utils"
	"go.uber.org/zap"
)

const (
	// ApiKeyLength represents the length of an API key to generate
	ApiKeyLength = 32
)

// APIKeysAPIService is a service that implements the logic for the APIKeysAPIServicer
// This service should implement the business logic for every endpoint for the APIKeysAPI API.
type APIKeysAPIService struct {
	apiKeyClient  dal.APIKeyManager
	serviceClient dal.ServiceManager
	logger        *zap.Logger
}

// AuthApiKey implements openapi.APIKeysAPIServicer.
func (s *APIKeysAPIService) AuthApiKey(context.Context, string, string, openapi.AuthApiKeyRequest) (openapi.ImplResponse, error) {
	panic("unimplemented")
}

// NewAPIKeysAPIService creates a default app service
func NewAPIKeysAPIService(apiKeyClient dal.APIKeyManager, serviceClient dal.ServiceManager, logger *zap.Logger) openapi.APIKeysAPIServicer {
	return &APIKeysAPIService{apiKeyClient: apiKeyClient, serviceClient: serviceClient, logger: logger}
}

// DeleteApiKey - Delete a specific API key
func (s *APIKeysAPIService) DeleteApiKey(ctx context.Context, serviceId string, keyId string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the service exists
	service, err := s.serviceClient.GetService(ctx, orgID, serviceId)
	if err != nil {
		s.logger.Error("failed to get service",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if service == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("service not found")
	}

	// Check if the API key exists
	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, keyId)
	if err != nil {
		s.logger.Error("failed to get API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if apiKey == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("API key not found")
	}

	err = s.apiKeyClient.DeleteAPIKey(ctx, orgID, serviceId, keyId)
	if err != nil {
		s.logger.Error("failed to delete API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GenerateApiKey - Generate a new API key with specific scopes for a service
func (s *APIKeysAPIService) GenerateApiKey(ctx context.Context, serviceId string, apiKeyInput openapi.ApiKeyInput) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the service exists
	service, err := s.serviceClient.GetService(ctx, orgID, serviceId)
	if err != nil {
		s.logger.Error("failed to get service",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if service == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("service not found")
	}

	keySecret, err := utils.GenerateSecret(ApiKeyLength)
	if err != nil {
		s.logger.Error("failed to generate API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	apiKey := dal.APIKey{
		ServiceID: serviceId,
		OrgID:     orgID,
		Secret:    keySecret,
		Scopes:    apiKeyInput.Scopes,
	}

	err = s.apiKeyClient.CreateAPIKey(ctx, &apiKey)
	if err != nil {
		s.logger.Error("failed to create API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.ApiKey{
		Id:        apiKey.APIKeyID,
		Secret:    apiKey.Secret,
		Scopes:    apiKey.Scopes,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return openapi.Response(http.StatusCreated, response), nil
}

// GetApiKey - Retrieve a specific API key
func (s *APIKeysAPIService) GetApiKey(ctx context.Context, serviceId string, keyId string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the service exists
	service, err := s.serviceClient.GetService(ctx, orgID, serviceId)
	if err != nil {
		s.logger.Error("failed to get service",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if service == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("service not found")
	}

	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, keyId)
	if err != nil {
		s.logger.Error("failed to get API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if apiKey == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("API key not found")
	}

	createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.ApiKey{
		Id:        apiKey.APIKeyID,
		Secret:    apiKey.Secret,
		Scopes:    apiKey.Scopes,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}

// ListApiKeys - List all API keys for a service
func (s *APIKeysAPIService) ListApiKeys(ctx context.Context, serviceId string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the service exists
	service, err := s.serviceClient.GetService(ctx, orgID, serviceId)
	if err != nil {
		s.logger.Error("failed to get service",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if service == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("service not found")
	}

	apiKeys, err := s.apiKeyClient.ListAPIKeysByService(ctx, orgID, serviceId)
	if err != nil {
		s.logger.Error("failed to list API keys",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	responses := make([]openapi.ApiKey, len(apiKeys))
	for i, apiKey := range apiKeys {
		createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
		if err != nil {
			s.logger.Error("failed to parse timestamp",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
		if err != nil {
			s.logger.Error("failed to parse timestamp",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		responses[i] = openapi.ApiKey{
			Id:        apiKey.APIKeyID,
			Secret:    apiKey.Secret,
			Scopes:    apiKey.Scopes,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return openapi.Response(http.StatusOK, responses), nil
}

// UpdateApiKey - Update an API key's scopes
func (s *APIKeysAPIService) UpdateApiKey(ctx context.Context, serviceId string, keyId string, apiKeyInput openapi.ApiKeyInput) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the service exists
	service, err := s.serviceClient.GetService(ctx, orgID, serviceId)
	if err != nil {
		s.logger.Error("failed to get service",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if service == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("service not found")
	}

	// Check if the API key exists
	apiKey, err := s.apiKeyClient.GetAPIKey(ctx, keyId)
	if err != nil {
		s.logger.Error("failed to get API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if apiKey == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("API key not found")
	}

	// Update the API key with the new values
	apiKey.Scopes = apiKeyInput.Scopes
	err = s.apiKeyClient.UpdateAPIKey(ctx, apiKey)
	if err != nil {
		s.logger.Error("failed to update API key",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(apiKey.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(apiKey.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.ApiKey{
		Id:        apiKey.APIKeyID,
		Secret:    apiKey.Secret,
		Scopes:    apiKey.Scopes,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}
