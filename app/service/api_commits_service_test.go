package service

import (
	"context"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/payloadops/plato/app/openapi"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestCommitsAPIService_CreateBranchCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	mockCommitClient := mocks.NewMockCommitManager(ctrl)
	service := NewCommitsAPIService(mockProjectClient, mockPromptClient, mockBranchClient, mockCommitClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	ctx = context.WithValue(ctx, "userID", "user1")
	projectID := "proj1"
	promptID := "prompt1"
	branchName := "branch1"
	commitInput := openapi.CommitInput{
		Message: "Initial commit",
		Content: "This is the first commit",
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().GetBranch(ctx, "org1", promptID, branchName).Return(&dal.Branch{}, nil)
	mockCommitClient.EXPECT().CreateCommit(ctx, "org1", projectID, promptID, branchName, gomock.Any()).Return(nil)

	response, err := service.CreateBranchCommit(ctx, projectID, promptID, branchName, commitInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	commit, ok := response.Body.(openapi.Commit)
	assert.True(t, ok)
	assert.Equal(t, commitInput.Message, commit.Message)
	assert.Equal(t, commitInput.Content, commit.Content)
}

func TestCommitsAPIService_GetBranchCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	mockCommitClient := mocks.NewMockCommitManager(ctrl)
	service := NewCommitsAPIService(mockProjectClient, mockPromptClient, mockBranchClient, mockCommitClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	branchName := "branch1"
	commitID := "commit1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().GetBranch(ctx, "org1", promptID, branchName).Return(&dal.Branch{}, nil)
	mockCommitClient.EXPECT().GetCommit(ctx, "org1", projectID, promptID, branchName, commitID).Return(&dal.Commit{
		CommitID:  commitID,
		Message:   "Initial commit",
		Content:   "This is the first commit",
		UserID:    "user1",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil)

	response, err := service.GetBranchCommit(ctx, projectID, promptID, branchName, commitID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	commit, ok := response.Body.(openapi.Commit)
	assert.True(t, ok)
	assert.Equal(t, commitID, commit.Id)
	assert.Equal(t, "Initial commit", commit.Message)
}

func TestCommitsAPIService_GetTemplateCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	mockCommitClient := mocks.NewMockCommitManager(ctrl)
	service := NewCommitsAPIService(mockProjectClient, mockPromptClient, mockBranchClient, mockCommitClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	ctx = context.WithValue(ctx, "projectID", "proj1")

	projectID := "proj1"
	promptID := "prompt1"
	branchName := "branch1"
	commitID := "commit1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().GetBranch(ctx, "org1", promptID, branchName).Return(&dal.Branch{}, nil)
	mockCommitClient.EXPECT().GetCommit(ctx, "org1", projectID, promptID, branchName, commitID).Return(&dal.Commit{
		CommitID:  commitID,
		Message:   "Initial commit",
		Content:   "This is the first commit",
		UserID:    "user1",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil)

	response, err := service.GetTemplateCommit(ctx, promptID, branchName, commitID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	commit, ok := response.Body.(openapi.Commit)
	assert.True(t, ok)
	assert.Equal(t, commitID, commit.Id)
	assert.Equal(t, "Initial commit", commit.Message)
}

func TestCommitsAPIService_ListBranchCommits(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	mockCommitClient := mocks.NewMockCommitManager(ctrl)
	service := NewCommitsAPIService(mockProjectClient, mockPromptClient, mockBranchClient, mockCommitClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	branchName := "branch1"

	commits := []dal.Commit{
		{CommitID: "commit1", Message: "Initial commit", Content: "This is the first commit", UserID: "user1", CreatedAt: time.Now().UTC().Format(time.RFC3339)},
		{CommitID: "commit2", Message: "Second commit", Content: "This is the second commit", UserID: "user2", CreatedAt: time.Now().UTC().Format(time.RFC3339)},
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().GetBranch(ctx, "org1", promptID, branchName).Return(&dal.Branch{}, nil)
	mockCommitClient.EXPECT().ListCommitsByBranch(ctx, "org1", projectID, promptID, branchName).Return(commits, nil)

	response, err := service.ListBranchCommits(ctx, projectID, promptID, branchName)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	listedCommits, ok := response.Body.([]openapi.Commit)
	assert.True(t, ok)
	assert.Equal(t, 2, len(listedCommits))
	assert.Equal(t, "commit1", listedCommits[0].Id)
	assert.Equal(t, "commit2", listedCommits[1].Id)
}
