package project

import (
	"context"
	projectservicemodel "plato/app_1/go/model/project"
)

type ProjectService interface {
	CreateProject(ctx context.Context, createProjectRequest projectservicemodel.CreateProjectRequest) (*projectservicemodel.CreateProjectResponse, error)
	GetProject(ctx context.Context, projectId string) (*projectservicemodel.GetProjectResponse, error)
	ListProjectsByOrg(ctx context.Context) (*projectservicemodel.ListProjectsResponse, error)
	ListProjectsByTeam(ctx context.Context, teamId string) (*projectservicemodel.ListProjectsResponse, error)
	UpdateProject(ctx context.Context, projectId string, updateProjectRequest projectservicemodel.UpdateProjectRequest) (*projectservicemodel.UpdateProjectResponse, error)
	DeleteProject(ctx context.Context, projectId string) (*projectservicemodel.DeleteProjectResponse, error)
}
