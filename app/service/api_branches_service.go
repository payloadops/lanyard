package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
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
func (s *BranchesAPIService) CreatePromptBranch(ctx context.Context, projectID string, promptID string, branchInput openapi.BranchInput) (openapi.ImplResponse, error) {
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

	branch := &dal.Branch{
		Name: branchInput.Name,
	}

	err = s.branchClient.CreateBranch(ctx, orgID, promptID, branch)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, branch), nil
}

// DeleteBranch - Delete a specific branch
func (s *BranchesAPIService) DeleteBranch(ctx context.Context, projectID string, promptID string, branchName string) (openapi.ImplResponse, error) {
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

	// Check if the branch exists and belongs to the specified prompt
	branch, err := s.branchClient.GetBranch(ctx, orgID, promptID, branchName)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	err = s.branchClient.DeleteBranch(ctx, orgID, promptID, branchName)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetBranch - Retrieve a specific branch
func (s *BranchesAPIService) GetBranch(ctx context.Context, projectID string, promptID string, branchID string) (openapi.ImplResponse, error) {
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

	branch, err := s.branchClient.GetBranch(ctx, orgID, promptID, branchID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	return openapi.Response(http.StatusOK, branch), nil
}

// ListPromptBranches - List all branches of a specific prompt
func (s *BranchesAPIService) ListPromptBranches(ctx context.Context, projectID string, promptID string) (openapi.ImplResponse, error) {
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

	branches, err := s.branchClient.ListBranchesByPrompt(ctx, orgID, promptID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, branches), nil
}
