package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/api/utils"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// OrganizationsAPIService is a service that implements the logic for the OrganizationsAPIServicer
// This service should implement the business logic for every endpoint for the OrganizationsAPI API.
type OrganizationsAPIService struct {
	client dal.OrganizationManager
}

// NewOrganizationsAPIService creates a default api service
func NewOrganizationsAPIService() openapi.OrganizationsAPIServicer {
	client, err := dal.NewOrgDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create OrgDBClient: %v", err))
	}
	return &OrganizationsAPIService{client: client}
}

// CreateOrganization - Create a new organization
func (s *OrganizationsAPIService) CreateOrganization(ctx context.Context, organizationInput openapi.OrganizationInput) (openapi.ImplResponse, error) {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	org := dal.Organization{
		ID:          ksuid,
		Name:        organizationInput.Name,
		Description: organizationInput.Description,
	}

	err = s.client.CreateOrganization(ctx, org)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, org), nil
}

// DeleteOrganization - Delete a specific organization
func (s *OrganizationsAPIService) DeleteOrganization(ctx context.Context, orgId string) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.client.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	err = s.client.DeleteOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetOrganization - Get a specific organization
func (s *OrganizationsAPIService) GetOrganization(ctx context.Context, orgId string) (openapi.ImplResponse, error) {
	org, err := s.client.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	return openapi.Response(http.StatusOK, org), nil
}

// ListOrganizations - List all organizations
func (s *OrganizationsAPIService) ListOrganizations(ctx context.Context) (openapi.ImplResponse, error) {
	orgs, err := s.client.ListOrganizations(ctx)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, orgs), nil
}

// UpdateOrganization - Update a specific organization
func (s *OrganizationsAPIService) UpdateOrganization(ctx context.Context, orgId string, organizationInput openapi.OrganizationInput) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.client.GetOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	// Update the organization with the new values
	org.Name = organizationInput.Name
	org.Description = organizationInput.Description

	err = s.client.UpdateOrganization(ctx, *org)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, org), nil
}
