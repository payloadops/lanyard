package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/app/utils"
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
func (s *PromptsAPIService) CreatePrompt(ctx context.Context, projectID string, promptInput openapi.PromptInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
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

	err = s.promptClient.CreatePrompt(ctx, orgID, projectID, prompt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	response := openapi.Prompt{
		Id:          prompt.PromptID,
		Name:        prompt.Name,
		Description: prompt.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return openapi.Response(http.StatusCreated, response), nil
}

// DeletePrompt - Delete a specific prompt from a project
func (s *PromptsAPIService) DeletePrompt(ctx context.Context, projectID string, promptID string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	err = s.promptClient.DeletePrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetPrompt - Retrieve a specific prompt within a project
func (s *PromptsAPIService) GetPrompt(ctx context.Context, projectID string, promptID string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	response := openapi.Prompt{
		Id:          prompt.PromptID,
		Name:        prompt.Name,
		Description: prompt.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}

// ListPrompts - List all prompts in a project
func (s *PromptsAPIService) ListPrompts(ctx context.Context, projectID string) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	prompts, err := s.promptClient.ListPromptsByProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	responses := make([]openapi.Prompt, len(prompts))
	for i, prompt := range prompts {
		createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), err
		}

		updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), err
		}

		responses[i] = openapi.Prompt{
			Id:          prompt.PromptID,
			Name:        prompt.Name,
			Description: prompt.Description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
	}

	return openapi.Response(http.StatusOK, responses), nil
}

// UpdatePrompt - Update a specific prompt in a project
func (s *PromptsAPIService) UpdatePrompt(ctx context.Context, projectID string, promptID string, promptInput openapi.PromptInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	// Update the prompt with the new values
	prompt.Name = promptInput.Name
	prompt.Description = promptInput.Description

	err = s.promptClient.UpdatePrompt(ctx, orgID, projectID, prompt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	response := openapi.Prompt{
		Id:          prompt.PromptID,
		Name:        prompt.Name,
		Description: prompt.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}
