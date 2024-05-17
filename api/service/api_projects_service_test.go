package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectManager is a mock implementation of the ProjectManager interface
type MockProjectManager struct {
	mock.Mock
}

func (m *MockProjectManager) CreateProject(ctx context.Context, project dal.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectManager) GetProject(ctx context.Context, id string) (*dal.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.Project), args.Error(1)
}

func (m *MockProjectManager) UpdateProject(ctx context.Context, project dal.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectManager) DeleteProject(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectManager) ListProjects(ctx context.Context) ([]dal.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.Project), args.Error(1)
}

func (m *MockProjectManager) ListProjectsByOrganization(ctx context.Context, organizationID string) ([]dal.Project, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).([]dal.Project), args.Error(1)
}

func (m *MockProjectManager) ListProjectsByTeam(ctx context.Context, teamID string) ([]dal.Project, error) {
	args := m.Called(ctx, teamID)
	return args.Get(0).([]dal.Project), args.Error(1)
}

func TestCreateProject(t *testing.T) {
	mockProjectClient := new(MockProjectManager)
	mockOrgClient := new(MockOrganizationManager)
	service := ProjectsAPIService{client: mockProjectClient, orgClient: mockOrgClient}

	projectInput := openapi.ProjectInput{
		OrgId:       "org1",
		TeamId:      "team1",
		Name:        "Test Project",
		Description: "Test Description",
	}

	expectedProject := dal.Project{
		ID:          "foo",
		OrgID:       projectInput.OrgId,
		TeamID:      projectInput.TeamId,
		Name:        projectInput.Name,
		Description: projectInput.Description,
	}

	mockOrgClient.On("GetOrganization", mock.Anything, projectInput.OrgId).Return(&dal.Organization{}, nil)
	mockProjectClient.On("CreateProject", mock.Anything, expectedProject).Return(nil)

	resp, err := service.CreateProject(context.Background(), projectInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockProjectClient.AssertExpectations(t)
}

func TestDeleteProject(t *testing.T) {
	mockProjectClient := new(MockProjectManager)
	mockOrgClient := new(MockOrganizationManager)
	service := ProjectsAPIService{client: mockProjectClient, orgClient: mockOrgClient}

	projectId := "1"

	// Test case where project exists
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{ID: projectId, OrgID: "org1"}, nil)
	mockProjectClient.On("DeleteProject", mock.Anything, projectId).Return(nil)

	resp, err := service.DeleteProject(context.Background(), projectId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.DeleteProject(context.Background(), projectId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
}

func TestGetProject(t *testing.T) {
	mockProjectClient := new(MockProjectManager)
	mockOrgClient := new(MockOrganizationManager)
	service := ProjectsAPIService{client: mockProjectClient, orgClient: mockOrgClient}

	projectId := "1"
	project := &dal.Project{
		ID:          projectId,
		OrgID:       "org1",
		TeamID:      "team1",
		Name:        "Test Project",
		Description: "Test Description",
	}

	// Test case where project exists
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(project, nil)

	resp, err := service.GetProject(context.Background(), projectId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.GetProject(context.Background(), projectId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
}

func TestListProjects(t *testing.T) {
	mockProjectClient := new(MockProjectManager)
	mockOrgClient := new(MockOrganizationManager)
	service := ProjectsAPIService{client: mockProjectClient, orgClient: mockOrgClient}

	projects := []dal.Project{
		{ID: "1", OrgID: "org1", TeamID: "team1", Name: "Project1", Description: "Description1"},
		{ID: "2", OrgID: "org1", TeamID: "team1", Name: "Project2", Description: "Description2"},
	}

	mockProjectClient.On("ListProjects", mock.Anything).Return(projects, nil)

	resp, err := service.ListProjects(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
}

func TestUpdateProject(t *testing.T) {
	mockProjectClient := new(MockProjectManager)
	mockOrgClient := new(MockOrganizationManager)
	service := ProjectsAPIService{client: mockProjectClient, orgClient: mockOrgClient}

	projectId := "1"
	projectInput := openapi.ProjectInput{
		Name:        "Updated Project",
		Description: "Updated Description",
		OrgId:       "org1",
		TeamId:      "team1",
	}
	project := &dal.Project{
		ID:          projectId,
		OrgID:       projectInput.OrgId,
		TeamID:      projectInput.TeamId,
		Name:        projectInput.Name,
		Description: projectInput.Description,
	}

	// Test case where organization and project exist
	mockOrgClient.On("GetOrganization", mock.Anything, projectInput.OrgId).Return(&dal.Organization{}, nil)
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{ID: projectId, OrgID: projectInput.OrgId}, nil)
	mockProjectClient.On("UpdateProject", mock.Anything, *project).Return(nil)

	resp, err := service.UpdateProject(context.Background(), projectId, projectInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockProjectClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, projectInput.OrgId).Return((*dal.Organization)(nil), nil)

	resp, err = service.UpdateProject(context.Background(), projectId, projectInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)

	// Test case where project does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, projectInput.OrgId).Return(&dal.Organization{}, nil)
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.UpdateProject(context.Background(), projectId, projectInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockProjectClient.AssertExpectations(t)
}
