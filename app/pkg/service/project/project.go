package project

import (
	"context"
	"fmt"
	"plato/app/pkg/auth"
	awsclient "plato/app/pkg/client/aws"
	dbdal "plato/app/pkg/dal/postgres"
	projectservicemodel "plato/app/pkg/model/project"
	"plato/app/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

type Service struct {
	s3Client *s3.Client
}

func NewService() (ProjectService, error) {
	return &Service{s3Client: awsclient.GetS3Client()}, nil
}

func (s *Service) CreateProject(ctx context.Context, createProjectRequest projectservicemodel.CreateProjectRequest) (*projectservicemodel.CreateProjectResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return nil, fmt.Errorf("failed to org id from context")
	}

	projectId := util.GenIDString()

	var createProjectResponse projectservicemodel.CreateProjectResponse
	_, err := dbdal.AddProject(
		ctx,
		orgId,
		createProjectRequest.TeamId,
		projectId,
		createProjectRequest.Name,
		// createProjectRequest.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project record with err: %w", err)
	}

	_, err = s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(fmt.Sprintf("%s/%s", orgId, projectId)),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create project bucket with err: %w", err)
	}

	_, err = s.s3Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to enable bucket versioning with err: %w", err)
	}

	return &createProjectResponse, err
}

func (s *Service) GetProject(ctx context.Context, projectId string) (*projectservicemodel.GetProjectResponse, error) {
	// orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	// if !ok {
	// 	return fmt.Errorf("failed to org id from context")
	// }

	var getProjectResponse projectservicemodel.GetProjectResponse
	_, err := dbdal.GetProject(
		ctx,
		projectId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch project record with err: %w", err)
	}

	return &getProjectResponse, err
}

func (s *Service) ListProjectsByOrg(ctx context.Context) (*projectservicemodel.ListProjectsResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return nil, fmt.Errorf("failed to org id from context")
	}

	var listProjectsResponse projectservicemodel.ListProjectsResponse
	_, err := dbdal.ListProjectsByOrgId(
		ctx,
		orgId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to list project records with err: %w", err)
	}
	return &listProjectsResponse, err
}

func (s *Service) ListProjectsByTeam(ctx context.Context, teamId string) (*projectservicemodel.ListProjectsResponse, error) {
	// orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	// if !ok {
	// 	return fmt.Errorf("failed to org id from context")
	// }

	var listProjectsResponse projectservicemodel.ListProjectsResponse
	_, err := dbdal.ListProjectsByTeamId(
		ctx,
		teamId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to list project records with err: %w", err)
	}
	return &listProjectsResponse, err
}

func (s *Service) UpdateProject(ctx context.Context, projectId string, updateProjectRequest projectservicemodel.UpdateProjectRequest) (*projectservicemodel.UpdateProjectResponse, error) {
	// orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	// if !ok {
	// 	return fmt.Errorf("failed to org id from context")
	// }

	var updateProjectResponse projectservicemodel.UpdateProjectResponse
	_, err := dbdal.UpdateProject(
		ctx,
		projectId,
		updateProjectRequest.Name,
		updateProjectRequest.Description,
		updateProjectRequest.TeamId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update project record with err: %w", err)
	}

	return &updateProjectResponse, err
}

func (s *Service) DeleteProject(ctx context.Context, projectId string) (*projectservicemodel.DeleteProjectResponse, error) {
	// orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	// if !ok {
	// 	return nil, fmt.Errorf("failed to org id from context")
	// }

	var deleteProjectResponse projectservicemodel.DeleteProjectResponse
	_, err := dbdal.SoftDeleteProject(
		ctx,
		projectId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to soft delete project record with err: %w", err)
	}
	return &deleteProjectResponse, err
}
