package dbdal

import (
	"context"
	"fmt"
	dbClient "plato/app/pkg/client/db"

	"github.com/uptrace/bun"
)

// Prompt represents the structure of a prompt record in the database.
type Prompt struct {
	bun.BaseModel `bun:"table:prompts,alias:p"`
	ID            string `bun:"prompt_id,pk" json:"prompt_id"`
	ProjectID     string `bun:"project_id" json:"project_id"`
	PromptS3Path  string `bun:"prompt_s3_path" json:"prompt_s3_path"`
	Deleted       bool   `bun:"deleted" json:"deleted"`
	Version       string `bun:"version" json:"version"`
	Stub          string `bun:"stub" json:"stub"`
}

func ListPromptsByProjectId(ctx context.Context, projectId string) (*[]Prompt, error) {
	prompts := &[]Prompt{}
	err := dbClient.GetClient().NewSelect().Model(prompts).Where("project_id = ?", projectId).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting prompt: %w", err)
	}
	return prompts, nil
}

func GetPromptById(ctx context.Context, id string) (*Prompt, error) {
	prompt := &Prompt{}
	err := dbClient.GetClient().NewSelect().Model(prompt).Where("prompt_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting prompt: %w", err)
	}
	return prompt, nil
}

func AddPrompt(ctx context.Context, stub string, projectID string, promptId string, promptS3Path string, version string) (*Prompt, error) {
	prompt := &Prompt{
		ProjectID:    projectID,
		ID:           promptId,
		PromptS3Path: promptS3Path,
		Version:      version,
		Deleted:      false,
		Stub:         stub,
	}
	_, err := dbClient.GetClient().NewInsert().Model(prompt).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error adding prompt: %w", err)
	}
	return prompt, nil
}

func UpdatePrompt(ctx context.Context, id string, stub string, version string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&Prompt{}).Set("stub = ? AND version = ?", stub, version).Where("prompt_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating prompt deleted status: %w", err)
	}
	return nil
}

func UpdatePromptDeletedStatus(ctx context.Context, id string, deleted bool) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&Prompt{}).Set("deleted = ?", deleted).Where("prompt_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating prompt deleted status: %w", err)
	}
	return nil
}
