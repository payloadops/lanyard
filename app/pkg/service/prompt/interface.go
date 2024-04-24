package promptservice

import (
	"context"
	"plato/app/pkg/model"
)

type PromptService interface {
	GetPrompt(
		ctx context.Context,
		projectId string,
		promptId string,
		branch string,
	) (*model.GetPromptResponse, error)
	UpdatePrompt(
		ctx context.Context,
		projectId string,
		promptId string,
		updatePromptRequest model.UpdatePromptRequest,
	) (*model.GetPromptResponse, error)
	CreatePrompt(
		ctx context.Context,
		projectId string,
		createPromptRequest model.CreatePromptRequest,
	) (*model.GetPromptResponse, error)
	DeletePrompt(
		ctx context.Context,
		projectId string,
		promptId string,
	) (*model.DeletePromptResponse, error)
	ListPrompts(
		ctx context.Context,
		projectId string,
	) (*model.ListPromptsResponse, error)
	// CreateBranch(
	// 	ctx context.Context,
	// 	projectId string,
	// 	promptId string,
	// 	createBranchRequest model.CreateBranchRequest,
	// ) (*model.CreateBranchResponse, error)
	// ListBranches(
	// 	ctx context.Context,
	// 	projectId string,
	// 	promptId string,
	// ) (*model.ListBranchesResponse, error)
	UpdateActiveVersion(
		ctx context.Context,
		projectId string,
		promptId string,
		updateActiveVersionRequest model.UpdateActiveVersionRequest,
	) (*model.GetPromptResponse, error)
}
