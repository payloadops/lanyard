package dbdal

import (
	"context"
	"fmt"

	dbClient "plato/app/pkg/client/db"
	"plato/app/pkg/util"

	"github.com/uptrace/bun"
)

type APIKeyItem struct {
	bun.BaseModel `bun:"table:api_keys,alias:ak"`
	ApiKey        string   `bun:"api_key,pk" json:"api_key"`
	ProjectId     string   `bun:"project_id" json:"project_id"`
	OrgId         string   `bun:"org_id" json:"org_id"`
	TeamId        string   `bun:"team_id" json:"team_id"`
	RateLimit     int      `bun:"rate_limit" json:"rate_limit"`
	Active        bool     `bun:"active" json:"active"`
	Scopes        []string `bun:"scopes,array" json:"scopes"`
}

// GetApiKey retrieves an API key by its string value from the database.
func GetApiKey(ctx context.Context, apiKeyString string) (*APIKeyItem, error) {
	apiKey := &APIKeyItem{}
	err := dbClient.GetClient().NewSelect().Model(apiKey).Where("api_key = ?", apiKeyString).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying API key: %w", err)
	}
	return apiKey, nil
}

// CreateApiKey creates a new API key in the database.
func CreateApiKey(ctx context.Context, projectId string, desc string, scopes []string) (*APIKeyItem, error) {
	apiKey := &APIKeyItem{
		ApiKey:    util.GenUUIDString(),
		ProjectId: projectId,
		RateLimit: 1000,
		Active:    true,
		Scopes:    scopes,
	}
	_, err := dbClient.GetClient().NewInsert().Model(apiKey).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating API key: %w", err)
	}
	return apiKey, nil
}

// UpdateApiKey updates an existing API key's description and scopes.
func UpdateApiKey(ctx context.Context, apiKeyId, newDesc string, newScopes []string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&APIKeyItem{}).Set("description = ?", newDesc).Set("scopes = ?", newScopes).Where("api_key = ?", apiKeyId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating API key: %w", err)
	}
	return nil
}

// DeactivateApiKey deactivates a specific API key.
func DeactivateApiKey(ctx context.Context, apiKeyId string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&APIKeyItem{}).Set("active = false").Where("api_key = ?", apiKeyId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deactivating API key: %w", err)
	}
	return nil
}
