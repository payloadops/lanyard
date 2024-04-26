package org

import "context"

type OrgService interface {
	CreateOrg(ctx context.Context, createOrgRequest orgservicemodel.CreateOrgRequest) (model.CreateOrgResponse, error)
	GetOrg(ctx context.Context, orgId string) (orgservicemodel.GetOrgResponse, error)
	UpdateOrg(ctx context.Context, updateOrgRequest orgservicemodel.UpdateOrgRequest) (model.UpdateOrgResponse, error)
	DeleteOrg(ctx context.Context) (orgservicemodel.DeleteOrgResponse, error)
}
