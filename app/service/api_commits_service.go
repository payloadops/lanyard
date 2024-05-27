package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
)

// CommitsAPIService is a service that implements the logic for the CommitsAPIServicer
// This service should implement the business logic for every endpoint for the CommitsAPI API.
type CommitsAPIService struct {
	projectClient dal.ProjectManager
	promptClient  dal.PromptManager
	branchClient  dal.BranchManager
	commitClient  dal.CommitManager
}

// NewCommitsAPIService creates a default app service
func NewCommitsAPIService(
	projectClient dal.ProjectManager,
	promptClient dal.PromptManager,
	branchClient dal.BranchManager,
	commitClient dal.CommitManager,
) openapi.CommitsAPIServicer {
	return &CommitsAPIService{
		projectClient: projectClient,
		promptClient:  promptClient,
		branchClient:  branchClient,
		commitClient:  commitClient,
	}
}

// CreateBranchCommit - Create a new commit for a branch
func (s *CommitsAPIService) CreateBranchCommit(ctx context.Context, projectId string, promptId string, branchId string, commitInput openapi.CommitInput) (openapi.ImplResponse, error) {
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

	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, promptId, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	commit := dal.Commit{
		// ID:       commitInput.Id,
		BranchID: branchId,
		// UserID:   commitInput.UserId,
		Message: commitInput.Message,
		Content: commitInput.Content,
	}

	err = s.commitClient.CreateCommit(ctx, commit)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusCreated, commit), nil
}

// GetBranchCommit - Retrieve a specific commit or the latest commit of a branch
func (s *CommitsAPIService) GetBranchCommit(ctx context.Context, projectId string, promptId string, branchId string, commitId string) (openapi.ImplResponse, error) {
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

	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, promptId, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	commit, err := s.commitClient.GetCommit(ctx, commitId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if commit == nil || commit.BranchID != branchId {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("commit not found")
	}

	return openapi.Response(http.StatusOK, commit), nil
}

// ListBranchCommits - List all commits of a specific branch
func (s *CommitsAPIService) ListBranchCommits(ctx context.Context, projectId string, promptId string, branchId string) (openapi.ImplResponse, error) {
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

	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, promptId, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	commits, err := s.commitClient.ListCommitsByBranch(ctx, branchId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, commits), nil
}
