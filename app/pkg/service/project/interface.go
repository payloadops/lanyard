package project

import (
	"context"
	"plato/app/pkg/model"
)

type ProjectService interface {
	CreateProject(ctx context.Context, createProjectRequest model.CreateProjectRequest) (model.CreateProjectResponse, error)
	GetProject(ctx context.Context, projectId string) (model.GetProjectResponse, error)
	ListProjectsByOrg(ctx context.Context) (model.ListProjectsResponse, error)
	ListProjectsByTeam(ctx context.Context, teamId string) (model.ListProjectsResponse, error)
	UpdateProject(ctx context.Context, projectId string) error
	DeleteProject(ctx context.Context, projectId string) error
}
