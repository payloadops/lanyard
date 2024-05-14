package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_project_db_client.go "github.com/payloadops/plato/api/dal" ProjectManager

// ProjectManager defines the operations available for managing projects.
type ProjectManager interface {
	CreateProject(ctx context.Context, project Project) error
	GetProject(ctx context.Context, id string) (*Project, error)
	UpdateProject(ctx context.Context, project Project) error
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context) ([]Project, error)
	ListProjectsByOrganization(ctx context.Context, organizationID string) ([]Project, error)
	ListProjectsByTeam(ctx context.Context, teamID string) ([]Project, error)
}

// Ensure ProjectDBClient implements the ProjectManager interface
var _ ProjectManager = &ProjectDBClient{}

// Project represents a project in the system.
type Project struct {
	ID          string `json:"id"`
	OrgID       string `json:"orgId,omitempty"`
	TeamID      string `json:"teamId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ProjectDBClient is a client for interacting with DynamoDB for project-related operations.
type ProjectDBClient struct {
	service *dynamodb.Client
}

// NewProjectDBClient creates a new ProjectDBClient with AWS configuration.
func NewProjectDBClient() (*ProjectDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	return &ProjectDBClient{
		service: svc,
	}, nil
}

// CreateProject creates a new project in the DynamoDB table.
func (d *ProjectDBClient) CreateProject(ctx context.Context, project Project) error {
	now := time.Now().UTC().Format(time.RFC3339)
	project.CreatedAt = now
	project.UpdatedAt = now

	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetProject retrieves a project by ID from the DynamoDB table.
func (d *ProjectDBClient) GetProject(ctx context.Context, id string) (*Project, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := d.service.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var project Project
	err = attributevalue.UnmarshalMap(result.Item, &project)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	return &project, nil
}

// UpdateProject updates an existing project in the DynamoDB table.
func (d *ProjectDBClient) UpdateProject(ctx context.Context, project Project) error {
	project.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteProject deletes a project by ID from the DynamoDB table.
func (d *ProjectDBClient) DeleteProject(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := d.service.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %v", err)
	}

	return nil
}

// ListProjects retrieves all projects from the DynamoDB table.
func (d *ProjectDBClient) ListProjects(ctx context.Context) ([]Project, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Projects"),
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return projects, nil
}

// ListProjectsByOrganization retrieves all projects for a specific organization from the DynamoDB table.
func (d *ProjectDBClient) ListProjectsByOrganization(ctx context.Context, organizationID string) ([]Project, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Projects"),
		FilterExpression: aws.String("orgId = :orgId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":orgId": &types.AttributeValueMemberS{Value: organizationID},
		},
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return projects, nil
}

// ListProjectsByTeam retrieves all projects for a specific team from the DynamoDB table.
func (d *ProjectDBClient) ListProjectsByTeam(ctx context.Context, teamID string) ([]Project, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Projects"),
		FilterExpression: aws.String("teamId = :teamId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":teamId": &types.AttributeValueMemberS{Value: teamID},
		},
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return projects, nil
}
