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
		branch string,
	) (*model.GetPromptResponse, error)
	CreatePrompt(
		ctx context.Context,
		projectId string,
		prompt string,
		branch string,
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
}
