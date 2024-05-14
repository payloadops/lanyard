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

// MockBranchManager is a mock implementation of the BranchManager interface
type MockBranchManager struct {
	mock.Mock
}

func (m *MockBranchManager) CreateBranch(ctx context.Context, branch dal.Branch) error {
	args := m.Called(ctx, branch)
	return args.Error(0)
}

func (m *MockBranchManager) GetBranch(ctx context.Context, id string) (*dal.Branch, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.Branch), args.Error(1)
}

func (m *MockBranchManager) DeleteBranch(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBranchManager) ListBranches(ctx context.Context) ([]dal.Branch, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.Branch), args.Error(1)
}

func (m *MockBranchManager) ListBranchesByPrompt(ctx context.Context, promptID string) ([]dal.Branch, error) {
	args := m.Called(ctx, promptID)
	return args.Get(0).([]dal.Branch), args.Error(1)
}

func TestCreatePromptBranch(t *testing.T) {
	mockBranchClient := new(MockBranchManager)
	mockPromptClient := new(MockPromptManager)
	service := BranchesAPIService{branchClient: mockBranchClient, promptClient: mockPromptClient}

	promptId := "prompt1"
	branchInput := openapi.BranchInput{
		Id: "branch1",
	}
	expectedBranch := dal.Branch{
		ID:       branchInput.Id,
		PromptID: promptId,
	}

	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{}, nil)
	mockBranchClient.On("CreateBranch", mock.Anything, expectedBranch).Return(nil)

	resp, err := service.CreatePromptBranch(context.Background(), promptId, branchInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockPromptClient.AssertExpectations(t)
	mockBranchClient.AssertExpectations(t)
}

func TestDeleteBranch(t *testing.T) {
	mockBranchClient := new(MockBranchManager)
	mockPromptClient := new(MockPromptManager)
	service := BranchesAPIService{branchClient: mockBranchClient, promptClient: mockPromptClient}

	promptId := "prompt1"
	branchId := "branch1"

	// Test case where prompt and branch exist and belong to the prompt
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{}, nil)
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return(&dal.Branch{ID: branchId, PromptID: promptId}, nil)
	mockBranchClient.On("DeleteBranch", mock.Anything, branchId).Return(nil)

	resp, err := service.DeleteBranch(context.Background(), promptId, branchId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.Code)
	mockPromptClient.AssertExpectations(t)
	mockBranchClient.AssertExpectations(t)

	// Test case where prompt does not exist
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return((*dal.Prompt)(nil), nil)

	resp, err = service.DeleteBranch(context.Background(), promptId, branchId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockPromptClient.AssertExpectations(t)

	// Test case where branch does not exist
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{}, nil)
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return((*dal.Branch)(nil), nil)

	resp, err = service.DeleteBranch(context.Background(), promptId, branchId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockPromptClient.AssertExpectations(t)
	mockBranchClient.AssertExpectations(t)
}

func TestGetBranch(t *testing.T) {
	mockBranchClient := new(MockBranchManager)
	mockPromptClient := new(MockPromptManager)
	service := BranchesAPIService{branchClient: mockBranchClient, promptClient: mockPromptClient}

	promptId := "prompt1"
	branchId := "branch1"
	branch := &dal.Branch{
		ID:       branchId,
		PromptID: promptId,
	}

	// Test case where prompt and branch exist and belong to the prompt
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{}, nil)
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return(branch, nil)

	resp, err := service.GetBranch(context.Background(), promptId, branchId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockPromptClient.AssertExpectations(t)
	mockBranchClient.AssertExpectations(t)

	// Test case where prompt does not exist
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return((*dal.Prompt)(nil), nil)

	resp, err = service.GetBranch(context.Background(), promptId, branchId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockPromptClient.AssertExpectations(t)

	// Test case where branch does not exist
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{}, nil)
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return((*dal.Branch)(nil), nil)

	resp, err = service.GetBranch(context.Background(), promptId, branchId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockPromptClient.AssertExpectations(t)
	mockBranchClient.AssertExpectations(t)
}

func TestListPromptBranches(t *testing.T) {
	mockBranchClient := new(MockBranchManager)
	mockPromptClient := new(MockPromptManager)
	service := BranchesAPIService{branchClient: mockBranchClient, promptClient: mockPromptClient}

	promptId := "prompt1"
	branches := []dal.Branch{
		{ID: "1", PromptID: promptId},
		{ID: "2", PromptID: promptId},
	}

	// Test case where prompt exists
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return(&dal.Prompt{}, nil)
	mockBranchClient.On("ListBranchesByPrompt", mock.Anything, promptId).Return(branches, nil)

	resp, err := service.ListPromptBranches(context.Background(), promptId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockPromptClient.AssertExpectations(t)
	mockBranchClient.AssertExpectations(t)

	// Test case where prompt does not exist
	mockPromptClient.On("GetPrompt", mock.Anything, promptId).Return((*dal.Prompt)(nil), nil)

	resp, err = service.ListPromptBranches(context.Background(), promptId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockPromptClient.AssertExpectations(t)
}
