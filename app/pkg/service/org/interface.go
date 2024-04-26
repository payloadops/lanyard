package org

import "context"

type OrgService interface {
	CreateOrg(ctx context.Context, createOrgRequest model.CreateOrgRequest) (model.CreateOrgResponse, error)
	GetOrg(ctx context.Context, orgId string) (model.GetOrgResponse, error)
	UpdateOrg(ctx context.Context, projectId string) error
	DeleteOrg(ctx context.Context, projectId string) error
}
