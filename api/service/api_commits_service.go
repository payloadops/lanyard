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
	commitClient dal.CommitManager
	branchClient dal.BranchManager
}

// NewCommitsAPIService creates a default api service
func NewCommitsAPIService() openapi.CommitsAPIServicer {
	/*commitClient, err := dal.NewCommitDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create CommitDBClient: %v", err))
	}
	branchClient, err := dal.NewBranchDBClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create BranchDBClient: %v", err))
	}
	return &CommitsAPIService{commitClient: commitClient, branchClient: branchClient}
	*/
	return nil
}

// CreateBranchCommit - Create a new commit for a branch
func (s *CommitsAPIService) CreateBranchCommit(ctx context.Context, promptId string, branchId string, commitInput openapi.CommitInput) (openapi.ImplResponse, error) {
	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, branchId)
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
func (s *CommitsAPIService) GetBranchCommit(ctx context.Context, promptId string, branchId string, commitId string) (openapi.ImplResponse, error) {
	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, branchId)
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
func (s *CommitsAPIService) ListBranchCommits(ctx context.Context, promptId string, branchId string) (openapi.ImplResponse, error) {
	// Check if the branch exists
	branch, err := s.branchClient.GetBranch(ctx, branchId)
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
