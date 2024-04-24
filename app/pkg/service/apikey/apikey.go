package apikey

import (
	"context"
	"errors"

	dbdal "plato/app/pkg/dal/postgres"
)

type service struct {
}

// NewService creates a new Api key management service.
func NewService() ApiKeyService {
	return &service{}
}

func (s *service) ListApiKeys(ctx context.Context, projectId string) (*[]dbdal.ApiKeyItem, error) {
	if projectId == "" {
		return nil, errors.New("project ID cannot be empty")
	}

	keys, err := dbdal.ListApiKeysByProjectId(ctx, projectId)

	return keys, err
}

func (s *service) Mint(ctx context.Context, projectId, desc string, scopes []string) (*dbdal.ApiKeyItem, error) {
	if projectId == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	if len(scopes) == 0 {
		return nil, errors.New("scopes cannot be empty")
	}
	return dbdal.CreateApiKey(ctx, projectId, desc, scopes)
}

func (s *service) GetApiKey(ctx context.Context, keyId string) (*dbdal.ApiKeyItem, error) {
	if keyId == "" {
		return nil, errors.New("key ID cannot be empty")
	}
	return dbdal.GetApiKey(ctx, keyId)
}

func (s *service) UpdateApiKey(ctx context.Context, keyId, newDesc string, newScopes []string) error {
	if keyId == "" {
		return errors.New("key ID cannot be empty")
	}
	return dbdal.UpdateApiKey(ctx, keyId, newDesc, newScopes)
}

func (s *service) DeleteApiKey(ctx context.Context, keyId string) error {
	if keyId == "" {
		return errors.New("key ID cannot be empty")
	}
	return dbdal.DeactivateApiKey(ctx, keyId)
}
