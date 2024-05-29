package service

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrganizationManager is a mock implementation of the OrganizationManager interface
type MockOrganizationManager struct {
	mock.Mock
}

func (m *MockOrganizationManager) CreateOrganization(ctx context.Context, org dal.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationManager) GetOrganization(ctx context.Context, id string) (*dal.Organization, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.Organization), args.Error(1)
}

func (m *MockOrganizationManager) UpdateOrganization(ctx context.Context, org dal.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrganizationManager) DeleteOrganization(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrganizationManager) ListOrganizations(ctx context.Context) ([]dal.Organization, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.Organization), args.Error(1)
}

func TestCreateOrganization(t *testing.T) {
	mockClient := new(MockOrganizationManager)
	service := OrganizationsAPIService{client: mockClient}

	orgInput := openapi.OrganizationInput{
		Name:        "Test Org",
		Description: "Test Description",
	}

	expectedOrg := dal.Organization{
		ID:          "foo",
		Name:        orgInput.Name,
		Description: orgInput.Description,
	}

	mockClient.On("CreateOrganization", mock.Anything, expectedOrg).Return(nil)

	resp, err := service.CreateOrganization(context.Background(), orgInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestDeleteOrganization(t *testing.T) {
	mockClient := new(MockOrganizationManager)
	service := OrganizationsAPIService{client: mockClient}

	orgID := "1"

	// Test case where organization exists
	mockClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockClient.On("DeleteOrganization", mock.Anything, orgID).Return(nil)

	resp, err := service.DeleteOrganization(context.Background(), orgID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.DeleteOrganization(context.Background(), orgID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestGetOrganization(t *testing.T) {
	mockClient := new(MockOrganizationManager)
	service := OrganizationsAPIService{client: mockClient}

	orgID := "1"
	org := &dal.Organization{
		ID:          orgID,
		Name:        "Test Org",
		Description: "Test Description",
	}

	// Test case where organization exists
	mockClient.On("GetOrganization", mock.Anything, orgID).Return(org, nil)

	resp, err := service.GetOrganization(context.Background(), orgID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.GetOrganization(context.Background(), orgID)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestListOrganizations(t *testing.T) {
	mockClient := new(MockOrganizationManager)
	service := OrganizationsAPIService{client: mockClient}

	orgs := []dal.Organization{
		{ID: "1", Name: "Org1", Description: "Description1"},
		{ID: "2", Name: "Org2", Description: "Description2"},
	}

	mockClient.On("ListOrganizations", mock.Anything).Return(orgs, nil)

	resp, err := service.ListOrganizations(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockClient.AssertExpectations(t)
}

func TestUpdateOrganization(t *testing.T) {
	mockClient := new(MockOrganizationManager)
	service := OrganizationsAPIService{client: mockClient}

	orgID := "1"
	orgInput := openapi.OrganizationInput{
		Name:        "Updated Org",
		Description: "Updated Description",
	}
	org := &dal.Organization{
		ID:          orgID,
		Name:        orgInput.Name,
		Description: orgInput.Description,
	}

	// Test case where organization exists
	mockClient.On("GetOrganization", mock.Anything, orgID).Return(&dal.Organization{}, nil)
	mockClient.On("UpdateOrganization", mock.Anything, *org).Return(nil)

	resp, err := service.UpdateOrganization(context.Background(), orgID, orgInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockClient.AssertExpectations(t)

	// Test case where organization does not exist
	mockClient.On("GetOrganization", mock.Anything, orgID).Return((*dal.Organization)(nil), nil)

	resp, err = service.UpdateOrganization(context.Background(), orgID, orgInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockClient.AssertExpectations(t)
}
*/
