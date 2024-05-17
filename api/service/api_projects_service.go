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
	project := dal.Project{
		OrgID:       projectInput.OrgId,
		TeamID:      projectInput.TeamId,
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
	// Check if the project exists
	project, err := s.client.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	err = s.client.DeleteProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetProject - Retrieve a project by ID
func (s *ProjectsAPIService) GetProject(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	project, err := s.client.GetProject(ctx, projectId)
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
	projects, err := s.client.ListProjects(ctx)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, projects), nil
}

// UpdateProject - Update a project
func (s *ProjectsAPIService) UpdateProject(ctx context.Context, projectId string, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.client.GetProject(ctx, projectId)
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

	err = s.client.UpdateProject(ctx, *project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, project), nil
}
