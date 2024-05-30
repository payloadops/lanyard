package service

import (
	"context"
	"fmt"
	"github.com/payloadops/plato/app/utils"
	"net/http"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
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
func (s *CommitsAPIService) CreateBranchCommit(ctx context.Context, projectID string, promptID string, branchName string, commitInput openapi.CommitInput) (openapi.ImplResponse, error) {
	orgID, ok := ctx.Value("orgID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("org not found")
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("user not found")
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

	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, orgID, promptID, branchName)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	commit := &dal.Commit{
		UserID:  userID,
		Message: commitInput.Message,
		Content: commitInput.Content,
	}

	err = s.commitClient.CreateCommit(ctx, orgID, promptID, branchName, commit)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	createdAt, err := utils.ParseTimestamp(commit.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	response := openapi.Commit{
		Id:        commit.CommitID,
		Content:   commit.Content,
		Message:   commit.Message,
		UserId:    commit.UserID,
		CreatedAt: createdAt,
	}

	return openapi.Response(http.StatusCreated, response), nil
}

// GetBranchCommit - Retrieve a specific commit or the latest commit of a branch
func (s *CommitsAPIService) GetBranchCommit(ctx context.Context, projectID string, promptID string, branchName string, commitID string) (openapi.ImplResponse, error) {
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

	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, orgID, promptID, branchName)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	commit, err := s.commitClient.GetCommit(ctx, orgID, promptID, branchName, commitID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if commit == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("commit not found")
	}

	createdAt, err := utils.ParseTimestamp(commit.CreatedAt)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	response := openapi.Commit{
		Id:        commit.CommitID,
		Content:   commit.Content,
		Message:   commit.Message,
		UserId:    commit.UserID,
		CreatedAt: createdAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}

// ListBranchCommits - List all commits of a specific branch
func (s *CommitsAPIService) ListBranchCommits(ctx context.Context, projectID string, promptID string, branchName string) (openapi.ImplResponse, error) {
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

	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, orgID, promptID, branchName)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}
	if branch == nil {
		return openapi.Response(http.StatusNotFound, nil), fmt.Errorf("branch not found")
	}

	commits, err := s.commitClient.ListCommitsByBranch(ctx, orgID, promptID, branchName)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	responses := make([]openapi.Commit, len(commits))
	for i, commit := range commits {
		createdAt, err := utils.ParseTimestamp(commit.CreatedAt)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, nil), err
		}

		responses[i] = openapi.Commit{
			Id:        commit.CommitID,
			Content:   commit.Content,
			Message:   commit.Message,
			UserId:    commit.UserID,
			CreatedAt: createdAt,
		}
	}

	return openapi.Response(http.StatusOK, responses), nil
}
