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

// MockPromptManager is a mock implementation of the PromptManager interface
type MockPromptManager struct {
	mock.Mock
}

func (m *MockPromptManager) CreatePrompt(ctx context.Context, prompt dal.Prompt) error {
	args := m.Called(ctx, prompt)
	return args.Error(0)
}

func (m *MockPromptManager) GetPrompt(ctx context.Context, id string) (*dal.Prompt, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.Prompt), args.Error(1)
}

func (m *MockPromptManager) UpdatePrompt(ctx context.Context, prompt dal.Prompt) error {
	args := m.Called(ctx, prompt)
	return args.Error(0)
}

func (m *MockPromptManager) DeletePrompt(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPromptManager) ListPrompts(ctx context.Context) ([]dal.Prompt, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.Prompt), args.Error(1)
}

func (m *MockPromptManager) ListPromptsByProject(ctx context.Context, projectID string) ([]dal.Prompt, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]dal.Prompt), args.Error(1)
}

func TestCreatePrompt(t *testing.T) {
	mockPromptClient := new(MockPromptManager)
	mockProjectClient := new(MockProjectManager)
	service := PromptsAPIService{promptClient: mockPromptClient, projectClient: mockProjectClient}

	projectId := "project1"
	promptInput := openapi.PromptInput{
		Name:        "Test Prompt",
		Description: "Test Description",
	}

	expectedPrompt := dal.Prompt{
		ID:          "foo",
		ProjectID:   projectId,
		Name:        promptInput.Name,
		Description: promptInput.Description,
	}

	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("CreatePrompt", mock.Anything, expectedPrompt).Return(nil)

	resp, err := service.CreatePrompt(context.Background(), projectId, promptInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)
}

func TestDeletePrompt(t *testing.T) {
	mockPromptClient := new(MockPromptManager)
	mockProjectClient := new(MockProjectManager)
	service := PromptsAPIService{promptClient: mockPromptClient, projectClient: mockProjectClient}

	projectId := "project1"
	promptId := "prompt1"

	// Test case where project and prompt exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{ID: promptId, ProjectID: projectId}, nil)
	mockPromptClient.On("DeletePrompt", mock.Anything, promptId).Return(nil)

	resp, err := service.DeletePrompt(context.Background(), projectId, promptId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.DeletePrompt(context.Background(), projectId, promptId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where prompt does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return((*dal.Prompt)(nil), nil)

	resp, err = service.DeletePrompt(context.Background(), projectId, promptId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)
}

func TestGetPrompt(t *testing.T) {
	mockPromptClient := new(MockPromptManager)
	mockProjectClient := new(MockProjectManager)
	service := PromptsAPIService{promptClient: mockPromptClient, projectClient: mockProjectClient}

	projectId := "project1"
	promptId := "prompt1"
	prompt := &dal.Prompt{
		ID:          promptId,
		ProjectID:   projectId,
		Name:        "Test Prompt",
		Description: "Test Description",
	}

	// Test case where project and prompt exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(prompt, nil)

	resp, err := service.GetPrompt(context.Background(), projectId, promptId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.GetPrompt(context.Background(), projectId, promptId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where prompt does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return((*dal.Prompt)(nil), nil)

	resp, err = service.GetPrompt(context.Background(), projectId, promptId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)
}

func TestListPrompts(t *testing.T) {
	mockPromptClient := new(MockPromptManager)
	mockProjectClient := new(MockProjectManager)
	service := PromptsAPIService{promptClient: mockPromptClient, projectClient: mockProjectClient}

	projectId := "project1"
	prompts := []dal.Prompt{
		{ID: "1", ProjectID: projectId, Name: "Prompt1", Description: "Description1"},
		{ID: "2", ProjectID: projectId, Name: "Prompt2", Description: "Description2"},
	}

	// Test case where project exists
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("ListPromptsByProject", mock.Anything, projectId).Return(prompts, nil)

	resp, err := service.ListPrompts(context.Background(), projectId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.ListPrompts(context.Background(), projectId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
}

func TestUpdatePrompt(t *testing.T) {
	mockPromptClient := new(MockPromptManager)
	mockProjectClient := new(MockProjectManager)
	service := PromptsAPIService{promptClient: mockPromptClient, projectClient: mockProjectClient}

	projectId := "project1"
	promptId := "prompt1"
	promptInput := openapi.PromptInput{
		Name:        "Updated Prompt",
		Description: "Updated Description",
	}
	prompt := &dal.Prompt{
		ID:          promptId,
		ProjectID:   projectId,
		Name:        promptInput.Name,
		Description: promptInput.Description,
	}

	// Test case where project and prompt exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{ID: promptId, ProjectID: projectId}, nil)
	mockPromptClient.On("UpdatePrompt", mock.Anything, *prompt).Return(nil)

	resp, err := service.UpdatePrompt(context.Background(), projectId, promptId, promptInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)

	// Test case where project does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return((*dal.Project)(nil), nil)

	resp, err = service.UpdatePrompt(context.Background(), projectId, promptId, promptInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)

	// Test case where prompt does not exist
	mockProjectClient.On("GetProject", mock.Anything, projectId).Return(&dal.Project{}, nil)
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return((*dal.Prompt)(nil), nil)

	resp, err = service.UpdatePrompt(context.Background(), projectId, promptId, promptInput)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockProjectClient.AssertExpectations(t)
	mockPromptClient.AssertExpectations(t)
}
