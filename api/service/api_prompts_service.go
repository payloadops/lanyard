package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/api/utils"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// PromptsAPIService is a service that implements the logic for the PromptsAPIServicer
// This service should implement the business logic for every endpoint for the PromptsAPI API.
type PromptsAPIService struct {
	promptClient  dal.PromptManager
	projectClient dal.ProjectManager
}

// NewPromptsAPIService creates a default api service
func NewPromptsAPIService() openapi.PromptsAPIServicer {
	promptClient, err := dal.NewPromptDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create PromptDBClient: %v", err))
	}
	projectClient, err := dal.NewProjectDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create ProjectDBClient: %v", err))
	}
	return &PromptsAPIService{promptClient: promptClient, projectClient: projectClient}
}

// CreatePrompt - Create a new prompt in a project
func (s *PromptsAPIService) CreatePrompt(ctx context.Context, projectId string, promptInput openapi.PromptInput) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	prompt := dal.Prompt{
		ID:          ksuid,
		ProjectID:   projectId,
		Name:        promptInput.Name,
		Description: promptInput.Description,
	}

	err = s.promptClient.CreatePrompt(ctx, prompt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, prompt), nil
}

// DeletePrompt - Delete a specific prompt from a project
func (s *PromptsAPIService) DeletePrompt(ctx context.Context, projectId string, promptId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil || prompt.ProjectID != projectId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	err = s.promptClient.DeletePrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetPrompt - Retrieve a specific prompt within a project
func (s *PromptsAPIService) GetPrompt(ctx context.Context, projectId string, promptId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil || prompt.ProjectID != projectId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	return openapi.Response(http.StatusOK, prompt), nil
}

// ListPrompts - List all prompts in a project
func (s *PromptsAPIService) ListPrompts(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompts, err := s.promptClient.ListPromptsByProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, prompts), nil
}

// UpdatePrompt - Update a specific prompt in a project
func (s *PromptsAPIService) UpdatePrompt(ctx context.Context, projectId string, promptId string, promptInput openapi.PromptInput) (openapi.ImplResponse, error) {
	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil || prompt.ProjectID != projectId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	// Update the prompt with the new values
	prompt.Name = promptInput.Name
	prompt.Description = promptInput.Description
	prompt.ProjectID = projectId

	err = s.promptClient.UpdatePrompt(ctx, *prompt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, prompt), nil
}
