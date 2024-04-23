package dbdal

import (
	"context"
	"fmt"

	dbClient "plato/app/pkg/client/db"

	"github.com/google/uuid"
)

type APIKeyItem struct {
	ApiKey    string   `json:"api_key"`
	ProjectId string   `json:"project_id"`
	RateLimit int      `json:"rate_limit"`
	Active    bool     `json:"active"`
	Scopes    []string `json:"scopes"`
}

var API_KEYS_TABLE_NAME = "api_keys"

// GetApiKey retrieves an API key by its string value from the database.
func GetApiKey(ctx context.Context, apiKeyString string) (*APIKeyItem, error) {
	query := `SELECT api_key, project_id, rate_limit, active, scopes FROM ` + API_KEYS_TABLE_NAME + ` WHERE api_key = $1`
	row := dbClient.GetPGClient().QueryRow(ctx, query, apiKeyString)

	var apiKey APIKeyItem
	err := row.Scan(&apiKey.ApiKey, &apiKey.ProjectId, &apiKey.RateLimit, &apiKey.Active, &apiKey.Scopes)
	if err != nil {
		return nil, fmt.Errorf("error querying API key: %w", err)
	}
	return &apiKey, nil
}

// CreateApiKey creates a new API key in the database.
func CreateApiKey(ctx context.Context, projectId, desc string, scopes []string) (*APIKeyItem, error) {
	newUUID := mintUUID()
	query := `INSERT INTO ` + API_KEYS_TABLE_NAME + ` (api_key, project_id, description, rate_limit, active, scopes) VALUES ($1, $2, $3, $4, $5, $6) RETURNING api_key`
	var apiKey string
	err := dbClient.GetPGClient().QueryRow(ctx, query, newUUID, projectId, desc, 1000, true, scopes).Scan(&apiKey)
	if err != nil {
		return nil, fmt.Errorf("error creating API key: %w", err)
	}
	return &APIKeyItem{
		ApiKey:    apiKey,
		ProjectId: projectId,
		RateLimit: 1000,
		Active:    true,
		Scopes:    scopes,
	}, nil
}

// UpdateApiKey updates an existing API key's description and scopes.
func UpdateApiKey(ctx context.Context, apiKeyId, newDesc string, newScopes []string) error {
	query := `UPDATE ` + API_KEYS_TABLE_NAME + ` SET description = $1, scopes = $2 WHERE api_key = $3`
	_, err := dbClient.GetPGClient().Exec(ctx, query, newDesc, newScopes, apiKeyId)
	if err != nil {
		return fmt.Errorf("error updating API key: %w", err)
	}
	return nil
}

// DeactivateApiKey deactivates a specific API key.
func DeactivateApiKey(ctx context.Context, apiKeyId string) error {
	query := `UPDATE ` + API_KEYS_TABLE_NAME + ` SET active = false WHERE api_key = $1`
	_, err := dbClient.GetPGClient().Exec(ctx, query, apiKeyId)
	if err != nil {
		return fmt.Errorf("error deactivating API key: %w", err)
	}
	return nil
}

// mintUUID generates a new UUID string.
func mintUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}
