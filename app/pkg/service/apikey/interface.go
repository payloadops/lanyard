package apikey

import (
	"context"
	dbdal "plato/app/pkg/dal/postgres"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_apikey.go "plato/app/pkg/service/apikey" ApiKeyService

// ApiKeyService defines the operations available for managing Api keys.
type ApiKeyService interface {
	ListApiKeys(ctx context.Context, projectId string) (*[]dbdal.ApiKeyItem, error)
	Mint(ctx context.Context, projectId string, desc string, scopes []string) (*dbdal.ApiKeyItem, error)
	GetApiKey(ctx context.Context, keyId string) (*dbdal.ApiKeyItem, error)
	UpdateApiKey(ctx context.Context, projectId string, keyId string, newDesc string, newScopes []string) error
	DeleteApiKey(ctx context.Context, projectId string, keyId string) error
}
