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

// MockTeamManager is a mock implementation of the TeamManager interface
type MockTeamManager struct {
	mock.Mock
}

func (m *MockTeamManager) CreateTeam(ctx context.Context, team dal.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamManager) GetTeam(ctx context.Context, id string) (*dal.Team, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.Team), args.Error(1)
}

func (m *MockTeamManager) UpdateTeam(ctx context.Context, team dal.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamManager) DeleteTeam(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamManager) ListTeams(ctx context.Context) ([]dal.Team, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.Team), args.Error(1)
}

func (m *MockTeamManager) ListTeamsByOrganization(ctx context.Context, orgID string) ([]dal.Team, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).([]dal.Team), args.Error(1)
}

func TestCreateTeam(t *testing.T) {
	mockTeamClient := new(MockTeamManager)
	mockOrgClient := new(MockOrganizationManager)
	service := TeamsAPIService{teamClient: mockTeamClient, orgClient: mockOrgClient}

	orgID := "1"
	teamInput := openapi.TeamInput{
		Name:        "Test Team",
		Description: "Test Description",
	}

	expectedTeam := dal.Team{
		ID:          "foo",
		OrgID:       orgID,
		Name:        teamInput.Name,
		Description: teamInput.Description,
	}

	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("CreateTeam", mock.Anything, expectedTeam).Return(nil)

	resp, err := service.CreateTeam(context.Background(), orgID, teamInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)
}

func TestDeleteTeam(t *testing.T) {
	mockTeamClient := new(MockTeamManager)
	mockOrgClient := new(MockOrganizationManager)
	service := TeamsAPIService{teamClient: mockTeamClient, orgClient: mockOrgClient}

	orgID := "1"
	teamID := "1"

	// Test case where organization and team exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("GetTeam", mock.Anything, teamID).Return(&dal.Team{ID: teamID, OrgID: orgID}, nil)
	mockTeamClient.On("DeleteTeam", mock.Anything, teamID).Return(nil)

	resp, err := service.DeleteTeam(context.Background(), orgID, teamID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.DeleteTeam(context.Background(), orgID, teamID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)

	// Test case where team does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("GetTeam", mock.Anything, teamID).Return((*dal.Team)(nil), nil)

	resp, err = service.DeleteTeam(context.Background(), orgID, teamID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)
}

func TestGetTeam(t *testing.T) {
	mockTeamClient := new(MockTeamManager)
	mockOrgClient := new(MockOrganizationManager)
	service := TeamsAPIService{teamClient: mockTeamClient, orgClient: mockOrgClient}

	orgID := "1"
	teamID := "1"
	team := &dal.Team{
		ID:          teamID,
		OrgID:       orgID,
		Name:        "Test Team",
		Description: "Test Description",
	}

	// Test case where organization and team exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("GetTeam", mock.Anything, teamID).Return(team, nil)

	resp, err := service.GetTeam(context.Background(), orgID, teamID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.GetTeam(context.Background(), orgID, teamID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)

	// Test case where team does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("GetTeam", mock.Anything, teamID).Return((*dal.Team)(nil), nil)

	resp, err = service.GetTeam(context.Background(), orgID, teamID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)
}

func TestListTeams(t *testing.T) {
	mockTeamClient := new(MockTeamManager)
	mockOrgClient := new(MockOrganizationManager)
	service := TeamsAPIService{teamClient: mockTeamClient, orgClient: mockOrgClient}

	orgID := "1"
	teams := []dal.Team{
		{ID: "1", OrgID: orgID, Name: "Team1", Description: "Description1"},
		{ID: "2", OrgID: orgID, Name: "Team2", Description: "Description2"},
	}

	// Test case where organization exists
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("ListTeamsByOrganization", mock.Anything, orgID).Return(teams, nil)

	resp, err := service.ListTeams(context.Background(), orgID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.ListTeams(context.Background(), orgID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)
}

func TestUpdateTeam(t *testing.T) {
	mockTeamClient := new(MockTeamManager)
	mockOrgClient := new(MockOrganizationManager)
	service := TeamsAPIService{teamClient: mockTeamClient, orgClient: mockOrgClient}

	orgID := "1"
	teamID := "1"
	teamInput := openapi.TeamInput{
		Name:        "Updated Team",
		Description: "Updated Description",
	}
	team := &dal.Team{
		ID:          teamID,
		OrgID:       orgID,
		Name:        teamInput.Name,
		Description: teamInput.Description,
	}

	// Test case where organization and team exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("GetTeam", mock.Anything, teamID).Return(&dal.Team{ID: teamID, OrgID: orgID}, nil)
	mockTeamClient.On("UpdateTeam", mock.Anything, *team).Return(nil)

	resp, err := service.UpdateTeam(context.Background(), orgID, teamID, teamInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.UpdateTeam(context.Background(), orgID, teamID, teamInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)

	// Test case where team does not exist
	mockOrgClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockTeamClient.On("GetTeam", mock.Anything, teamID).Return((*dal.Team)(nil), nil)

	resp, err = service.UpdateTeam(context.Background(), orgID, teamID, teamInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockOrgClient.AssertExpectations(t)
	mockTeamClient.AssertExpectations(t)
}
