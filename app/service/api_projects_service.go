package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
	"net/http"
)

// ProjectsAPIService is a service that implements the logic for the ProjectsAPIServicer
// This service should implement the business logic for every endpoint for the ProjectsAPI API.
type ProjectsAPIService struct {
	client dal.ProjectManager
}

// NewProjectsAPIService creates a default api service
func NewProjectsAPIService(client dal.ProjectManager) openapi.ProjectsAPIServicer {
	return &ProjectsAPIService{client: client}
}

// CreateProject - Create a new project
func (s *ProjectsAPIService) CreateProject(ctx context.Context, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	project := &dal.Project{
		OrgID:       orgId,
		Name:        projectInput.Name,
		Description: projectInput.Description,
	}

	err := s.client.CreateProject(ctx, project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, project), nil
}

// DeleteProject - Delete a project
func (s *ProjectsAPIService) DeleteProject(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.client.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	err = s.client.DeleteProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetProject - Retrieve a project by ID
func (s *ProjectsAPIService) GetProject(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	project, err := s.client.GetProject(ctx, orgId, projectId)
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
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	projects, err := s.client.ListProjectsByOrganization(ctx, orgId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, projects), nil
}

// UpdateProject - Update a project
func (s *ProjectsAPIService) UpdateProject(ctx context.Context, projectId string, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.client.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Update the project with the new values
	project.Name = projectInput.Name
	project.Description = projectInput.Description

	err = s.client.UpdateProject(ctx, project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, project), nil
}
