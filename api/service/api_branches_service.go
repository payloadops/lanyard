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
	branchClient dal.BranchManager
	promptClient dal.PromptManager
}

// NewBranchesAPIService creates a default api service
func NewBranchesAPIService() openapi.BranchesAPIServicer {
	branchClient, err := dal.NewBranchDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create BranchDBClient: %v", err))
	}
	promptClient, err := dal.NewPromptDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create PromptDBClient: %v", err))
	}
	return &BranchesAPIService{branchClient: branchClient, promptClient: promptClient}
}

// CreatePromptBranch - Create a new branch for a prompt
func (s *BranchesAPIService) CreatePromptBranch(ctx context.Context, promptId string, branchInput openapi.BranchInput) (openapi.ImplResponse, error) {
	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	branch := dal.Branch{
		ID:       branchInput.Id,
		PromptID: promptId,
	}

	err = s.branchClient.CreateBranch(ctx, branch)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, branch), nil
}

// DeleteBranch - Delete a specific branch
func (s *BranchesAPIService) DeleteBranch(ctx context.Context, promptId string, branchId string) (openapi.ImplResponse, error) {
	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	// Check if the branch exists and belongs to the specified prompt
	branch, err := s.branchClient.GetBranch(ctx, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil || branch.PromptID != promptId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	err = s.branchClient.DeleteBranch(ctx, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetBranch - Retrieve a specific branch
func (s *BranchesAPIService) GetBranch(ctx context.Context, promptId string, branchId string) (openapi.ImplResponse, error) {
	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("prompt not found")
	}

	branch, err := s.branchClient.GetBranch(ctx, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil || branch.PromptID != promptId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	return openapi.Response(http.StatusOK, branch), nil
}

// ListPromptBranches - List all branches of a specific prompt
func (s *BranchesAPIService) ListPromptBranches(ctx context.Context, promptId string) (openapi.ImplResponse, error) {
	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, promptId)
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
