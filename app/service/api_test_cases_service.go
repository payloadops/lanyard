package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/payloadops/plato/app/utils"
	"go.uber.org/zap"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
)

// TestCasesAPIService is a service that implements the logic for the TestCasesAPIServicer
// This service should implement the business logic for every endpoint for the TestCaesAPI API.
type TestCasesAPIService struct {
	projectClient  dal.ProjectManager
	promptClient   dal.PromptManager
	testCaseClient dal.TestCaseManager
	logger         *zap.Logger
}

// NewTestCasesAPIService creates a default app service
func NewTestCasesAPIService(projectClient dal.ProjectManager, promptClient dal.PromptManager, logger *zap.Logger) openapi.TestCasesAPIServicer {
	return &TestCasesAPIService{projectClient: projectClient, promptClient: promptClient, logger: logger}
}

// CreatePrompt - Create a new prompt in a project
func (s *TestCasesAPIService) CreateTestCase(ctx context.Context, projectID string, testCaseInput openapi.TestCaseInput) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		s.logger.Error("failed to get project",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	prompt := &dal.Prompt{
		Name: testCaseInput.Name,
	}

	err = s.promptClient.CreatePrompt(ctx, orgID, projectID, prompt)
	if err != nil {
		s.logger.Error("failed to create prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
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

// DeleteTestCase - Delete a specific test case for a prompt
func (s *TestCasesAPIService) DeleteTestCase(ctx context.Context, projectID string, promptID string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		s.logger.Error("failed to get project",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	err = s.promptClient.DeletePrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		s.logger.Error("failed to delete prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetPrompt - Retrieve a specific prompt within a project
func (s *TestCasesAPIService) GetTestCase(ctx context.Context, projectID string, promptID string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		s.logger.Error("failed to get project",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
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

// ListTestCases - List all prompts in a project
func (s *TestCasesAPIService) ListTestCases(ctx context.Context, projectID string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		s.logger.Error("failed to get project",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	prompts, err := s.promptClient.ListPromptsByProject(ctx, orgID, projectID)
	if err != nil {
		s.logger.Error("failed to list prompts by project",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	responses := make([]openapi.Prompt, len(prompts))
	for i, prompt := range prompts {
		createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
		if err != nil {
			s.logger.Error("failed to parse timestamp",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
		if err != nil {
			s.logger.Error("failed to parse timestamp",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
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
func (s *TestCasesAPIService) UpdateTestCase(ctx context.Context, projectID string, promptID string, testCaseInput openapi.TestCaseInput) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the project exists
	project, err := s.projectClient.GetProject(ctx, orgID, projectID)
	if err != nil {
		s.logger.Error("failed to get project",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if project == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("project not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	err = s.promptClient.UpdatePrompt(ctx, orgID, projectID, prompt)
	if err != nil {
		s.logger.Error("failed to update prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(prompt.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(prompt.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
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
