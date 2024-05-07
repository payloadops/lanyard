// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Payload Ops API
 *
 * The Payload Ops API streamlines management of AI prompts, projects, organizations, teams, and users through conventional HTTP requests. This platform enables efficient automation and control of resources, providing robust tools for developers to manage settings, memberships, and activities seamlessly.
 *
 * API version: 1.0
 * Contact: info@payloadops.com
 */

package openapi

import (
	"context"
	"net/http"
	"errors"
)

// VersionsAPIService is a service that implements the logic for the VersionsAPIServicer
// This service should implement the business logic for every endpoint for the VersionsAPI API.
// Include any external packages or services that will be required by this service.
type VersionsAPIService struct {
}

// NewVersionsAPIService creates a default api service
func NewVersionsAPIService() VersionsAPIServicer {
	return &VersionsAPIService{}
}

// GetPromptVersion - Retrieve a specific version or the latest version of a prompt
func (s *VersionsAPIService) GetPromptVersion(ctx context.Context, promptId string, versionId string) (ImplResponse, error) {
	// TODO - update GetPromptVersion with the required logic for this service method.
	// Add api_versions_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, PromptVersion{}) or use other options such as http.Ok ...
	// return Response(200, PromptVersion{}), nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	// TODO: Uncomment the next line to return response Response(500, {}) or use other options such as http.Ok ...
	// return Response(500, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetPromptVersion method not implemented")
}

// ListPromptVersions - List all versions of a specific prompt
func (s *VersionsAPIService) ListPromptVersions(ctx context.Context, promptId string) (ImplResponse, error) {
	// TODO - update ListPromptVersions with the required logic for this service method.
	// Add api_versions_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, []PromptVersion{}) or use other options such as http.Ok ...
	// return Response(200, []PromptVersion{}), nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	// TODO: Uncomment the next line to return response Response(500, {}) or use other options such as http.Ok ...
	// return Response(500, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("ListPromptVersions method not implemented")
}
