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

func TestBranchesAPIService_CreatePromptBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	service := NewBranchesAPIService(mockProjectClient, mockPromptClient, mockBranchClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	branchInput := openapi.BranchInput{
		Name: "branch1",
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().CreateBranch(ctx, "org1", promptID, gomock.Any()).Return(nil)

	response, err := service.CreatePromptBranch(ctx, projectID, promptID, branchInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	branch, ok := response.Body.(*openapi.Branch)
	assert.True(t, ok)
	assert.Equal(t, branchInput.Name, branch.Name)
}

func TestBranchesAPIService_DeleteBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	service := NewBranchesAPIService(mockProjectClient, mockPromptClient, mockBranchClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	branchName := "branch1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().GetBranch(ctx, "org1", promptID, branchName).Return(&dal.Branch{}, nil)
	mockBranchClient.EXPECT().DeleteBranch(ctx, "org1", promptID, branchName).Return(nil)

	response, err := service.DeleteBranch(ctx, projectID, promptID, branchName)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestBranchesAPIService_GetBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	service := NewBranchesAPIService(mockProjectClient, mockPromptClient, mockBranchClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	branchID := "branch1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().GetBranch(ctx, "org1", promptID, branchID).Return(&dal.Branch{
		Name:      "branch1",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil)

	response, err := service.GetBranch(ctx, projectID, promptID, branchID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	branch, ok := response.Body.(*openapi.Branch)
	assert.True(t, ok)
	assert.Equal(t, "branch1", branch.Name)
}

func TestBranchesAPIService_ListPromptBranches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockBranchClient := mocks.NewMockBranchManager(ctrl)
	service := NewBranchesAPIService(mockProjectClient, mockPromptClient, mockBranchClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"

	branches := []dal.Branch{
		{Name: "branch1", CreatedAt: time.Now().UTC().Format(time.RFC3339)},
		{Name: "branch2", CreatedAt: time.Now().UTC().Format(time.RFC3339)},
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockBranchClient.EXPECT().ListBranchesByPrompt(ctx, "org1", promptID).Return(branches, nil)

	response, err := service.ListPromptBranches(ctx, projectID, promptID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	listedBranches, ok := response.Body.([]openapi.Branch)
	assert.True(t, ok)
	assert.Equal(t, 2, len(listedBranches))
	assert.Equal(t, "branch1", listedBranches[0].Name)
	assert.Equal(t, "branch2", listedBranches[1].Name)
}
