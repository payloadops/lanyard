package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/payloadops/plato/app/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_project_db_client.go "github.com/payloadops/plato/app/dal" ProjectManager

// ProjectManager defines the operations available for managing projects.
type ProjectManager interface {
	CreateProject(ctx context.Context, orgID string, project *Project) error
	GetProject(ctx context.Context, orgID string, projectID string) (*Project, error)
	UpdateProject(ctx context.Context, orgID string, project *Project) error
	DeleteProject(ctx context.Context, orgID string, projectID string) error
	ListProjectsByOrganization(ctx context.Context, orgID string) ([]Project, error)
}

// Ensure ProjectDBClient implements the ProjectManager interface
var _ ProjectManager = &ProjectDBClient{}

// Project represents a project in the system.
type Project struct {
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deleted     bool   `json:"deleted"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ProjectDBClient is a client for interacting with DynamoDB for project-related operations.
type ProjectDBClient struct {
	service DynamoDBAPI
}

// NewProjectDBClient creates a new ProjectDBClient.
func NewProjectDBClient(service DynamoDBAPI) *ProjectDBClient {
	return &ProjectDBClient{
		service: service,
	}
}

// createProjectCompositeKeys generates the partition key (pk) and sort key (SK) for a project.
func createProjectCompositeKeys(orgID, projectID string) (string, string) {
	return "Org#" + orgID, "Project#" + projectID
}

// CreateProject creates a new project in the DynamoDB table.
func (d *ProjectDBClient) CreateProject(ctx context.Context, orgID string, project *Project) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	project.ProjectID = ksuid
	pk, sk := createProjectCompositeKeys(orgID, project.ProjectID)

	now := time.Now().UTC().Format(time.RFC3339)
	project.CreatedAt = now
	project.UpdatedAt = now

	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Projects"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetProject retrieves a project by organization ID and project ID from the DynamoDB table.
func (d *ProjectDBClient) GetProject(ctx context.Context, orgID, projectID string) (*Project, error) {
	pk, sk := createProjectCompositeKeys(orgID, projectID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
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

	if project.Deleted {
		return nil, nil
	}

	return &project, nil
}

// UpdateProject updates an existing project in the DynamoDB table.
func (d *ProjectDBClient) UpdateProject(ctx context.Context, orgID string, project *Project) error {
	pk, sk := createProjectCompositeKeys(orgID, project.ProjectID)
	project.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(project)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String("Projects"),
		Item:                item,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(SK)"),
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteProject marks a project as deleted by organization ID and project ID in the DynamoDB table.
func (d *ProjectDBClient) DeleteProject(ctx context.Context, orgID, projectID string) error {
	pk, sk := createProjectCompositeKeys(orgID, projectID)
	now := time.Now().UTC().Format(time.RFC3339)

	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
		"UpdatedAt": {
			Value:  &types.AttributeValueMemberS{Value: now},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Projects"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(SK)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListProjectsByOrganization retrieves all projects for a specific organization from the DynamoDB table.
func (d *ProjectDBClient) ListProjectsByOrganization(ctx context.Context, orgID string) ([]Project, error) {
	pk, _ := createProjectCompositeKeys(orgID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Projects"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.service.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(result.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return projects, nil
}
