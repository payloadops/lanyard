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

func TestPromptsAPIService_CreatePrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	service := NewPromptsAPIService(mockProjectClient, mockPromptClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptInput := openapi.PromptInput{
		Name:        "Prompt1",
		Description: "Description1",
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().CreatePrompt(ctx, "org1", projectID, gomock.Any()).Return(nil)

	response, err := service.CreatePrompt(ctx, projectID, promptInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	prompt, ok := response.Body.(openapi.Prompt)
	assert.True(t, ok)
	assert.Equal(t, promptInput.Name, prompt.Name)
	assert.Equal(t, promptInput.Description, prompt.Description)
}

func TestPromptsAPIService_DeletePrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	service := NewPromptsAPIService(mockProjectClient, mockPromptClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{}, nil)
	mockPromptClient.EXPECT().DeletePrompt(ctx, "org1", projectID, promptID).Return(nil)

	response, err := service.DeletePrompt(ctx, projectID, promptID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestPromptsAPIService_GetPrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	service := NewPromptsAPIService(mockProjectClient, mockPromptClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{
		PromptID:    promptID,
		Name:        "Prompt1",
		Description: "Description1",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil)

	response, err := service.GetPrompt(ctx, projectID, promptID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	prompt, ok := response.Body.(openapi.Prompt)
	assert.True(t, ok)
	assert.Equal(t, promptID, prompt.Id)
}

func TestPromptsAPIService_ListPrompts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	service := NewPromptsAPIService(mockProjectClient, mockPromptClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"

	prompts := []dal.Prompt{
		{PromptID: "prompt1", Name: "Prompt1", Description: "Description1", CreatedAt: time.Now().UTC().Format(time.RFC3339), UpdatedAt: time.Now().UTC().Format(time.RFC3339)},
		{PromptID: "prompt2", Name: "Prompt2", Description: "Description2", CreatedAt: time.Now().UTC().Format(time.RFC3339), UpdatedAt: time.Now().UTC().Format(time.RFC3339)},
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().ListPromptsByProject(ctx, "org1", projectID).Return(prompts, nil)

	response, err := service.ListPrompts(ctx, projectID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	listedPrompts, ok := response.Body.([]openapi.Prompt)
	assert.True(t, ok)
	assert.Equal(t, 2, len(listedPrompts))
	assert.Equal(t, "prompt1", listedPrompts[0].Id)
	assert.Equal(t, "prompt2", listedPrompts[1].Id)
}

func TestPromptsAPIService_UpdatePrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProjectClient := mocks.NewMockProjectManager(ctrl)
	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	service := NewPromptsAPIService(mockProjectClient, mockPromptClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	promptInput := openapi.PromptInput{
		Name:        "UpdatedPrompt",
		Description: "UpdatedDescription",
	}

	mockProjectClient.EXPECT().GetProject(ctx, "org1", projectID).Return(&dal.Project{}, nil)
	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(&dal.Prompt{
		PromptID:    promptID,
		Name:        "Prompt1",
		Description: "Description1",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}, nil)
	mockPromptClient.EXPECT().UpdatePrompt(ctx, "org1", projectID, gomock.Any()).Return(nil)

	response, err := service.UpdatePrompt(ctx, projectID, promptID, promptInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	updatedPrompt, ok := response.Body.(openapi.Prompt)
	assert.True(t, ok)
	assert.Equal(t, promptInput.Name, updatedPrompt.Name)
	assert.Equal(t, promptInput.Description, updatedPrompt.Description)
}
