package service

import (
	"context"
	"errors"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
	"github.com/payloadops/plato/app/utils"
	"net/http"
)

// ProjectsAPIService is a service that implements the logic for the ProjectsAPIServicer
// This service should implement the business logic for every endpoint for the ProjectsAPI API.
type ProjectsAPIService struct {
	client dal.ProjectManager
}

// NewProjectsAPIService creates a default app service
func NewProjectsAPIService(client dal.ProjectManager) openapi.ProjectsAPIServicer {
	return &ProjectsAPIService{client: client}
}

// CreateProject - Create a new project
func (s *ProjectsAPIService) CreateProject(ctx context.Context, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	project := &dal.Project{
		Name:        projectInput.Name,
		Description: projectInput.Description,
	}

	err := s.client.CreateProject(ctx, orgID, project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(project.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(project.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.Project{
		Id:          project.ProjectID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return openapi.Response(http.StatusCreated, response), nil
}

// DeleteProject - Delete a project
func (s *ProjectsAPIService) DeleteProject(ctx context.Context, projectID string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.client.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	err = s.client.DeleteProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetProject - Retrieve a project by ID
func (s *ProjectsAPIService) GetProject(ctx context.Context, projectID string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	project, err := s.client.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	createdAt, err := utils.ParseTimestamp(project.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(project.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.Project{
		Id:          project.ProjectID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}

// ListProjects - List all projects
func (s *ProjectsAPIService) ListProjects(ctx context.Context) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	projects, err := s.client.ListProjectsByOrganization(ctx, orgID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	responses := make([]openapi.Project, len(projects))
	for i, project := range projects {
		createdAt, err := utils.ParseTimestamp(project.CreatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		updatedAt, err := utils.ParseTimestamp(project.UpdatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		responses[i] = openapi.Project{
			Id:          project.ProjectID,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
	}

	return openapi.Response(http.StatusOK, responses), nil
}

// UpdateProject - Update a project
func (s *ProjectsAPIService) UpdateProject(ctx context.Context, projectID string, projectInput openapi.ProjectInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.client.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	// Update the project with the new values
	project.Name = projectInput.Name
	project.Description = projectInput.Description

	err = s.client.UpdateProject(ctx, orgID, project)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(project.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(project.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	response := openapi.Project{
		Id:          project.ProjectID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}
