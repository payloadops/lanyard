package project

import (
	"context"
	"fmt"
	awsclient "plato/app/pkg/client/aws"
	dbdal "plato/app/pkg/dal/postgres"
	projectservicemodel "plato/app/pkg/model/project"
	"plato/app/pkg/auth"
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

	_, err := dbdal.AddProject(
		ctx,
		orgId,
		projectId,
		createProjectRequest.Name,
		createProjectRequest.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project record with err: %w", err)
	}

	bucketName := aws.String(util.GetBucketString(orgId, projectId))
	_, err = s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: bucketName,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create project bucket with err: %w", err)
	}

	_, err = s.s3Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: bucketName,
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to enable bucket versioning with err: %w", err)
	}

	createProjectResponse := &projectservicemodel.CreateProjectResponse{
		ProjectId:   projectId,
		OrgId:       orgId,
		TeamId:      createProjectRequest.TeamId,
		Name:        createProjectRequest.Name,
		Description: createProjectRequest.Description,
	}

	return createProjectResponse, err
}

func (s *Service) GetProject(ctx context.Context, projectId string) (*projectservicemodel.GetProjectResponse, error) {
	record, err := dbdal.GetProject(
		ctx,
		projectId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch project record with err: %w", err)
	}

	getProjectResponse := &projectservicemodel.GetProjectResponse{
		OrgId:       record.OrgId,
		ProjectId:   record.Id,
		TeamId:      record.TeamId,
		Name:        record.Name,
		Description: record.Description,
	}

	return getProjectResponse, err
}

func (s *Service) ListProjectsByOrg(ctx context.Context) (*projectservicemodel.ListProjectsResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return nil, fmt.Errorf("failed to org id from context")
	}

	records, err := dbdal.ListProjectsByOrgId(
		ctx,
		orgId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to list project records with err: %w", err)
	}

	listProjectsResponse := &projectservicemodel.ListProjectsResponse{
		OrgId:    orgId,
		Projects: &records,
	}
	return listProjectsResponse, err
}

func (s *Service) ListProjectsByTeam(ctx context.Context, teamId string) (*projectservicemodel.ListProjectsResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return nil, fmt.Errorf("failed to org id from context")
	}

	records, err := dbdal.ListProjectsByTeamId(
		ctx,
		teamId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to list project records with err: %w", err)
	}

	listProjectsResponse := &projectservicemodel.ListProjectsResponse{
		OrgId:    orgId,
		TeamId:   teamId,
		Projects: &records,
	}
	return listProjectsResponse, err
}

func (s *Service) UpdateProject(ctx context.Context, projectId string, updateProjectRequest projectservicemodel.UpdateProjectRequest) (*projectservicemodel.UpdateProjectResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return nil, fmt.Errorf("failed to org id from context")
	}

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

	updateProjectResponse := &projectservicemodel.UpdateProjectResponse{
		OrgId:       orgId,
		ProjectId:   projectId,
		TeamId:      updateProjectRequest.TeamId,
		Name:        updateProjectRequest.Name,
		Description: updateProjectRequest.Description,
	}
	return updateProjectResponse, err
}

func (s *Service) DeleteProject(ctx context.Context, projectId string) (*projectservicemodel.DeleteProjectResponse, error) {
	_, err := dbdal.SoftDeleteProject(
		ctx,
		projectId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to soft delete project record with err: %w", err)
	}

	deleteProjectResponse := &projectservicemodel.DeleteProjectResponse{}

	return deleteProjectResponse, err
}
