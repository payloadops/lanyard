package project

import (
	"context"
	"plato/app/pkg/model"
)

type ProjectService interface {
	CreateProject(ctx context.Context, createProjectRequest projectservicemodel.CreateProjectRequest) (model.CreateProjectResponse, error)
	GetProject(ctx context.Context, projectId string) (projectservicemodel.GetProjectResponse, error)
	ListProjectsByOrg(ctx context.Context) (projectservicemodel.ListProjectsResponse, error)
	ListProjectsByTeam(ctx context.Context, teamId string) (projectservicemodel.ListProjectsResponse, error)
	UpdateProject(ctx context.Context, projectId string) (projectservicemodel.UpdateProjectResponse, error)
	DeleteProject(ctx context.Context, projectId string) (projectservicemodel.DeleteProjectResponse, error)
}
