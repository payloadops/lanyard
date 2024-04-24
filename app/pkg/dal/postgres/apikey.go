package dbdal

import (
	"context"
	"fmt"

	dbClient "plato/app/pkg/client/db"
	"plato/app/pkg/util"

	"github.com/uptrace/bun"
)

type ApiKeyItem struct {
	bun.BaseModel `bun:"table:api_keys,alias:ak"`
	ApiKey        string   `bun:"api_key,pk" json:"api_key"`
	ProjectId     string   `bun:"project_id" json:"project_id"`
	OrgId         string   `bun:"org_id" json:"org_id"`
	TeamId        string   `bun:"team_id" json:"team_id"`
	RateLimit     int      `bun:"rate_limit" json:"rate_limit"`
	Active        bool     `bun:"active" json:"active"`
	Scopes        []string `bun:"scopes,array" json:"scopes"`
}

// Lists active Api keys by project id
func ListApiKeysByProjectId(ctx context.Context, projectId string) (*[]ApiKeyItem, error) {
	apiKeys := &[]ApiKeyItem{}
	err := dbClient.GetClient().NewSelect().Model(apiKeys).Where("project_id = ? AND active = true", projectId).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying Api keys: %w", err)
	}
	return apiKeys, nil
}

// GetApiKey retrieves an Api key by its string value from the database.
func GetApiKey(ctx context.Context, apiKeyString string) (*ApiKeyItem, error) {
	apiKey := &ApiKeyItem{}
	err := dbClient.GetClient().NewSelect().Model(apiKey).Where("api_key = ?", apiKeyString).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying Api key: %w", err)
	}
	return apiKey, nil
}

// CreateApiKey creates a new Api key in the database.
func CreateApiKey(ctx context.Context, projectId string, desc string, scopes []string) (*ApiKeyItem, error) {
	apiKey := &ApiKeyItem{
		ApiKey:    util.GenUUIDString(),
		ProjectId: projectId,
		RateLimit: 1000,
		Active:    true,
		Scopes:    scopes,
	}
	_, err := dbClient.GetClient().NewInsert().Model(apiKey).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating Api key: %w", err)
	}
	return apiKey, nil
}

// UpdateApiKey updates an existing Api key's description and scopes.
func UpdateApiKey(ctx context.Context, apiKeyId, newDesc string, newScopes []string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&ApiKeyItem{}).Set("description = ?", newDesc).Set("scopes = ?", newScopes).Where("api_key = ?", apiKeyId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating Api key: %w", err)
	}
	return nil
}

// DeactivateApiKey deactivates a specific Api key.
func DeactivateApiKey(ctx context.Context, apiKeyId string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&ApiKeyItem{}).Set("active = false").Where("api_key = ?", apiKeyId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deactivating Api key: %w", err)
	}
	return nil
}
