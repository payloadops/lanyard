package apikey

import (
	"context"
	"errors"
	"fmt"

	"plato/app/pkg/auth"
	dbdal "plato/app/pkg/dal/postgres"
)

var _ ApiKeyService = (*Service)(nil)

type Service struct {
}

// NewService creates a new Api key management Service.
func NewService() ApiKeyService {
	return &Service{}
}

func (s *Service) ListApiKeys(ctx context.Context, projectId string) (*[]dbdal.ApiKeyItem, error) {
	if projectId == "" {
		return nil, errors.New("project ID cannot be empty")
	}

	keys, err := dbdal.ListApiKeysByProjectId(ctx, projectId)

	return keys, err
}

func (s *Service) Mint(ctx context.Context, projectId, desc string, scopes []string) (*dbdal.ApiKeyItem, error) {
	orgId, orgIdOk := ctx.Value(auth.OrgContext{}).(string)

	// Check if all required context values are successfully retrieved
	if !orgIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}
	if projectId == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	if len(scopes) == 0 {
		return nil, errors.New("scopes cannot be empty")
	}
	return dbdal.CreateApiKey(ctx, orgId, projectId, desc, scopes)
}

func (s *Service) GetApiKey(ctx context.Context, keyId string) (*dbdal.ApiKeyItem, error) {
	if keyId == "" {
		return nil, errors.New("key ID cannot be empty")
	}
	return dbdal.GetApiKey(ctx, keyId)
}

func (s *Service) UpdateApiKey(ctx context.Context, projectId string, keyId string, newDesc string, newScopes []string) error {
	if keyId == "" {
		return errors.New("key ID cannot be empty")
	}
	return dbdal.UpdateApiKey(ctx, projectId, keyId, newDesc, newScopes)
}

func (s *Service) DeleteApiKey(ctx context.Context, projectId string, keyId string) error {
	if keyId == "" {
		return errors.New("key ID cannot be empty")
	}
	return dbdal.DeactivateApiKey(ctx, projectId, keyId)
}
