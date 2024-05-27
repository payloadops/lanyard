package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// BranchesAPIService is a service that implements the logic for the BranchesAPIServicer
// This service should implement the business logic for every endpoint for the BranchesAPI API.
type BranchesAPIService struct {
	projectClient dal.ProjectManager
	promptClient  dal.PromptManager
	branchClient  dal.BranchManager
}

// NewBranchesAPIService creates a default app service
func NewBranchesAPIService(projectClient dal.ProjectManager, promptClient dal.PromptManager, branchClient dal.BranchManager) openapi.BranchesAPIServicer {
	return &BranchesAPIService{
		projectClient: projectClient,
		promptClient:  promptClient,
		branchClient:  branchClient,
	}
}

// CreatePromptBranch - Create a new branch for a prompt
func (s *BranchesAPIService) CreatePromptBranch(ctx context.Context, projectId string, promptId string, branchInput openapi.BranchInput) (openapi.ImplResponse, error) {
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
	prompt, err := s.promptClient.GetPrompt(ctx, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	branch := &dal.Branch{
		PromptID: promptId,
	}

	err = s.branchClient.CreateBranch(ctx, branch)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, branch), nil
}

// DeleteBranch - Delete a specific branch
func (s *BranchesAPIService) DeleteBranch(ctx context.Context, projectId string, promptId string, branchId string) (openapi.ImplResponse, error) {
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
	prompt, err := s.promptClient.GetPrompt(ctx, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	// Check if the branch exists and belongs to the specified prompt
	branch, err := s.branchClient.GetBranch(ctx, promptId, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil || branch.PromptID != promptId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	err = s.branchClient.DeleteBranch(ctx, promptId, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetBranch - Retrieve a specific branch
func (s *BranchesAPIService) GetBranch(ctx context.Context, projectId string, promptId string, branchId string) (openapi.ImplResponse, error) {
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
	prompt, err := s.promptClient.GetPrompt(ctx, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	branch, err := s.branchClient.GetBranch(ctx, promptId, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil || branch.PromptID != promptId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	return openapi.Response(http.StatusOK, branch), nil
}

// ListPromptBranches - List all branches of a specific prompt
func (s *BranchesAPIService) ListPromptBranches(ctx context.Context, projectId string, promptId string) (openapi.ImplResponse, error) {
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
	prompt, err := s.promptClient.GetPrompt(ctx, projectId, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	branches, err := s.branchClient.ListBranchesByPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, branches), nil
}
