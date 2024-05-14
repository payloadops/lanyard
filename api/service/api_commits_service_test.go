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

// MockCommitManager is a mock implementation of the CommitManager interface
type MockCommitManager struct {
	mock.Mock
}

func (m *MockCommitManager) CreateCommit(ctx context.Context, commit dal.Commit) error {
	args := m.Called(ctx, commit)
	return args.Error(0)
}

func (m *MockCommitManager) GetCommit(ctx context.Context, id string) (*dal.Commit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dal.Commit), args.Error(1)
}

func (m *MockCommitManager) ListCommits(ctx context.Context) ([]dal.Commit, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dal.Commit), args.Error(1)
}

func (m *MockCommitManager) ListCommitsByBranch(ctx context.Context, branchID string) ([]dal.Commit, error) {
	args := m.Called(ctx, branchID)
	return args.Get(0).([]dal.Commit), args.Error(1)
}

func TestCreateBranchCommit(t *testing.T) {
	mockCommitClient := new(MockCommitManager)
	mockBranchClient := new(MockBranchManager)
	service := CommitsAPIService{commitClient: mockCommitClient, branchClient: mockBranchClient}

	promptId := "prompt1"
	branchId := "branch1"
	commitInput := openapi.CommitInput{
		Message: "Initial commit",
		Content: "Commit content",
	}
	expectedCommit := dal.Commit{
		ID:       "id",
		BranchID: branchId,
		UserID:   "user_id",
		Message:  commitInput.Message,
		Content:  commitInput.Content,
	}

	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return(&dal.Branch{}, nil)
	mockCommitClient.On("CreateCommit", mock.Anything, expectedCommit).Return(nil)

	resp, err := service.CreateBranchCommit(context.Background(), promptId, branchId, commitInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	mockBranchClient.AssertExpectations(t)
	mockCommitClient.AssertExpectations(t)
}

func TestGetBranchCommit(t *testing.T) {
	mockCommitClient := new(MockCommitManager)
	mockBranchClient := new(MockBranchManager)
	service := CommitsAPIService{commitClient: mockCommitClient, branchClient: mockBranchClient}

	promptId := "prompt1"
	branchId := "branch1"
	commitId := "commit1"
	commit := &dal.Commit{
		ID:       commitId,
		BranchID: branchId,
	}

	// Test case where branch and commit exist and belong to the branch
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return(&dal.Branch{}, nil)
	mockCommitClient.On("GetCommit", mock.Anything, commitId).Return(commit, nil)

	resp, err := service.GetBranchCommit(context.Background(), promptId, branchId, commitId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockBranchClient.AssertExpectations(t)
	mockCommitClient.AssertExpectations(t)

	// Test case where branch does not exist
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return((*dal.Branch)(nil), nil)

	resp, err = service.GetBranchCommit(context.Background(), promptId, branchId, commitId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockBranchClient.AssertExpectations(t)

	// Test case where commit does not exist
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return(&dal.Branch{}, nil)
	mockCommitClient.On("GetCommit", mock.Anything, commitId).Return((*dal.Commit)(nil), nil)

	resp, err = service.GetBranchCommit(context.Background(), promptId, branchId, commitId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockBranchClient.AssertExpectations(t)
	mockCommitClient.AssertExpectations(t)
}

func TestListBranchCommits(t *testing.T) {
	mockCommitClient := new(MockCommitManager)
	mockBranchClient := new(MockBranchManager)
	service := CommitsAPIService{commitClient: mockCommitClient, branchClient: mockBranchClient}

	promptId := "prompt1"
	branchId := "branch1"
	commits := []dal.Commit{
		{ID: "1", BranchID: branchId},
		{ID: "2", BranchID: branchId},
	}

	// Test case where branch exists
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return(&dal.Branch{}, nil)
	mockCommitClient.On("ListCommitsByBranch", mock.Anything, branchId).Return(commits, nil)

	resp, err := service.ListBranchCommits(context.Background(), promptId, branchId)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockBranchClient.AssertExpectations(t)
	mockCommitClient.AssertExpectations(t)

	// Test case where branch does not exist
	mockBranchClient.On("GetBranch", mock.Anything, branchId).Return((*dal.Branch)(nil), nil)

	resp, err = service.ListBranchCommits(context.Background(), promptId, branchId)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockBranchClient.AssertExpectations(t)
}
