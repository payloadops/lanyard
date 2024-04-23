package apikey

import (
	"context"
	dbdal "plato/app/pkg/dal/postgres"
)

// APIKeyService defines the operations available for managing API keys.
type APIKeyService interface {
	Mint(ctx context.Context, projectId, desc string, scopes []string) (*dbdal.APIKeyItem, error)
	GetAPIKey(ctx context.Context, keyId string) (*dbdal.APIKeyItem, error)
	UpdateAPIKey(ctx context.Context, keyId, newDesc string, newScopes []string) error
	DeleteAPIKey(ctx context.Context, keyId string) error
}
