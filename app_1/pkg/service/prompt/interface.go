package promptservice

import (
	"context"
	promptservicemodel "plato/app_1/pkg/model/prompt/service"
)

type PromptService interface {
	GetPrompt(
		ctx context.Context,
		projectId string,
		promptId string,
		branch string,
		version string,
	) (*promptservicemodel.GetPromptResponse, error)
	UpdatePrompt(
		ctx context.Context,
		projectId string,
		promptId string,
		updatePromptRequest promptservicemodel.UpdatePromptRequest,
	) (*promptservicemodel.GetPromptResponse, error)
	CreatePrompt(
		ctx context.Context,
		projectId string,
		createPromptRequest promptservicemodel.CreatePromptRequest,
	) (*promptservicemodel.GetPromptResponse, error)
	DeletePrompt(
		ctx context.Context,
		projectId string,
		promptId string,
	) (*promptservicemodel.DeletePromptResponse, error)
	ListPrompts(
		ctx context.Context,
		projectId string,
	) (*promptservicemodel.ListPromptsResponse, error)
	CreateBranch(
		ctx context.Context,
		projectId string,
		promptId string,
		createBranchRequest promptservicemodel.CreateBranchRequest,
	) (*promptservicemodel.CreateBranchResponse, error)
	ListBranches(
		ctx context.Context,
		projectId string,
		promptId string,
	) (*promptservicemodel.ListBranchesResponse, error)
	DeleteBranch(
		ctx context.Context,
		projectId string,
		promptId string,
		branch string,
	) (*promptservicemodel.DeleteBranchResponse, error)
	ListVersions(
		ctx context.Context,
		projectId string,
		promptId string,
	) (*promptservicemodel.ListVersionsResponse, error)
	UpdateActiveVersion(
		ctx context.Context,
		projectId string,
		promptId string,
		updateActiveVersionRequest *promptservicemodel.UpdateActiveVersionRequest,
	) (*promptservicemodel.GetPromptResponse, error)
}
