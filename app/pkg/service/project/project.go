package project

import (
	"context"
	"fmt"
	"plato/app/pkg/auth"
	awsclient "plato/app/pkg/client/aws"
	dbdal "plato/app/pkg/dal/postgres"
	"plato/app/pkg/model"
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

func (s *Service) CreateProject(ctx context.Context, createProjectRequest model.CreateProjectRequest) (model.CreateProjectResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return fmt.Errorf("failed to org id from context")
	}

	projectId := util.GenIDString()

	record, err := dbdal.AddProject(
		orgId,
		projectId,
		&createProjectRequest.Name,
		&createProjectRequest.desc,
	)

	if err != nil {
		return fmt.Errorf("failed to create project record with err: %w", err)
	}

	_, err = s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(fmt.Sprintf("%s/%s", orgId, projectId)),
	})

	if err != nil {
		return fmt.Errorf("failed to create project bucket with err: %w", err)
	}

	_, err = s.s3Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to enable bucket versioning with err: %w", err)
	}

	return err
}

func (s *Service) GetProject(ctx context.Context, projectId string) (model.GetProjectResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return fmt.Errorf("failed to org id from context")
	}

	record, err := dbdal.GetProject(
		orgId,
		projectId,
	)

	if err != nil {
		return fmt.Errorf("failed to fetch project record with err: %w", err)
	}

	return err
}

func ListProjectsByOrg(ctx context.Context) (model.ListProjectsResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return fmt.Errorf("failed to org id from context")
	}

	records, err := dbdal.ListProjectsByOrgId(
		orgId,
	)

	if err != nil {
		return fmt.Errorf("failed to list project records with err: %w", err)
	}
	return listProjectsResponse, err
}

func ListProjectsByTeam(ctx context.Context, teamId string) (model.ListProjectsResponse, error) {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return fmt.Errorf("failed to org id from context")
	}

	records, err := dbdal.ListProjectsByTeamId(
		orgId,
		teamId,
	)

	if err != nil {
		return fmt.Errorf("failed to list project records with err: %w", err)
	}
	return listProjectsResponse, err
}

func UpdateProject(ctx context.Context, projectId string, updateProjectRequest model.UpdateProjectRequest) error {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return fmt.Errorf("failed to org id from context")
	}

	record, err := dbdal.UpdateProject(
		orgId,
		projectId,
		updateProjectRequest.Name,
		updateProjectRequest.Desc,
		updateProjectRequest.Team,
	)

	if err != nil {
		return fmt.Errorf("failed to update project record with err: %w", err)
	}
	return nil
}

func DeleteProject(ctx context.Context, projectId string) error {
	orgId, ok := ctx.Value(auth.OrgContext{}).(string)

	if !ok {
		return fmt.Errorf("failed to org id from context")
	}

	record, err := dbdal.SoftDeleteProject(
		orgId,
		projectId,
	)

	if err != nil {
		return fmt.Errorf("failed to soft delete project record with err: %w", err)
	}
	return nil
}
