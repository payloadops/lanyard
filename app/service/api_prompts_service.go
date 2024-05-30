package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
)

// PromptsAPIService is a service that implements the logic for the PromptsAPIServicer
// This service should implement the business logic for every endpoint for the PromptsAPI API.
type PromptsAPIService struct {
	projectClient dal.ProjectManager
	promptClient  dal.PromptManager
}

// NewPromptsAPIService creates a default app service
func NewPromptsAPIService(projectClient dal.ProjectManager, promptClient dal.PromptManager) openapi.PromptsAPIServicer {
	return &PromptsAPIService{projectClient: projectClient, promptClient: promptClient}
}

// CreatePrompt - Create a new prompt in a project
func (s *PromptsAPIService) CreatePrompt(ctx context.Context, projectId string, promptInput openapi.PromptInput) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompt := &dal.Prompt{
		Name:        promptInput.Name,
		Description: promptInput.Description,
	}

	err = s.promptClient.CreatePrompt(ctx, orgId, projectId, prompt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, prompt), nil
}

// DeletePrompt - Delete a specific prompt from a project
func (s *PromptsAPIService) DeletePrompt(ctx context.Context, projectId string, promptId string) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgId, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	err = s.promptClient.DeletePrompt(ctx, orgId, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetPrompt - Retrieve a specific prompt within a project
func (s *PromptsAPIService) GetPrompt(ctx context.Context, projectId string, promptId string) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompt, err := s.promptClient.GetPrompt(ctx, orgId, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	return openapi.Response(http.StatusOK, prompt), nil
}

// ListPrompts - List all prompts in a project
func (s *PromptsAPIService) ListPrompts(ctx context.Context, projectId string) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompts, err := s.promptClient.ListPromptsByProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, prompts), nil
}

// UpdatePrompt - Update a specific prompt in a project
func (s *PromptsAPIService) UpdatePrompt(ctx context.Context, projectId string, promptId string, promptInput openapi.PromptInput) (openapi.ImplResponse, error) {
	orgId, ok := ctx.Value("orgId").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgId, projectId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgId, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	// Update the prompt with the new values
	prompt.Name = promptInput.Name
	prompt.Description = promptInput.Description

	err = s.promptClient.UpdatePrompt(ctx, orgId, projectId, prompt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, prompt), nil
}
