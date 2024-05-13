package org

import (
	"context"
	orgservicemodel "plato/app_1/pkg/model/org"
)

type OrgService interface {
	CreateOrg(ctx context.Context, createOrgRequest orgservicemodel.CreateOrgRequest) (orgservicemodel.CreateOrgResponse, error)
	GetOrg(ctx context.Context, orgId string) (orgservicemodel.GetOrgResponse, error)
	UpdateOrg(ctx context.Context, updateOrgRequest orgservicemodel.UpdateOrgRequest) (orgservicemodel.UpdateOrgResponse, error)
	DeleteOrg(ctx context.Context) (orgservicemodel.DeleteOrgResponse, error)
}
