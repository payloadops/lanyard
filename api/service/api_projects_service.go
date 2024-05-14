package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/api/utils"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// ProjectsAPIService is a service that implements the logic for the ProjectsAPIServicer
// This service should implement the business logic for every endpoint for the ProjectsAPI API.
type ProjectsAPIService struct {
	projectClient dal.ProjectManager
	orgClient     dal.OrganizationManager
}

// NewProjectsAPIService creates a default api service
func NewProjectsAPIService() openapi.ProjectsAPIServicer {
	projectClient, err := dal.NewProjectDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create ProjectDBClient: %v", err))
	}
	orgClient, err := dal.NewOrgDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create OrgDBClient: %v", err))
	}
	return &ProjectsAPIService{projectClient: projectClient, orgClient: orgClient}
}

// CreateProject - Create a new project
func (s *ProjectsAPIService) CreateProject(ctx context.Context, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, projectInput.OrgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	project := dal.Project{
		ID:          ksuid,
		OrgID:       projectInput.OrgId,
		TeamID:      projectInput.TeamId,
		Name:        projectInput.Name,
		Description: projectInput.Description,
	}

	err = s.projectClient.CreateProject(ctx, project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, project), nil
}

// DeleteProject - Delete a project
func (s *ProjectsAPIService) DeleteProject(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	err = s.projectClient.DeleteProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetProject - Retrieve a project by ID
func (s *ProjectsAPIService) GetProject(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	return openapi.Response(http.StatusOK, project), nil
}

// ListProjects - List all projects
func (s *ProjectsAPIService) ListProjects(ctx context.Context) (openapi.ImplResponse, error) {
	projects, err := s.projectClient.ListProjects(ctx)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, projects), nil
}

// UpdateProject - Update a project
func (s *ProjectsAPIService) UpdateProject(ctx context.Context, projectId string, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	// Check if the organization exists
	org, err := s.orgClient.GetOrganization(ctx, projectInput.OrgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if org == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("organization not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Update the project with the new values
	project.Name = projectInput.Name
	project.Description = projectInput.Description
	project.OrgID = projectInput.OrgId
	project.TeamID = projectInput.TeamId

	err = s.projectClient.UpdateProject(ctx, *project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, project), nil
}
