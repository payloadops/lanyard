package dbdal

import (
	"context"
	"fmt"
	"time"

	awsclient "plato/app/go/client/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

var PROJECT_TABLE_NAME = aws.String("ProjectsPrompts")

type Project struct {
	Id          string `dynamodbav:"projectId" json:"project_id"`
	OrgId       string `dynamodbav:"orgId" json:"org_id"`
	TeamId      string `dynamodbav:"teamId" json:"team_id"`
	Name        string `dynamodbav:"name" json:"name"`
	Description string `dynamodbav:"description" json:"description"`
	Deleted     bool   `dynamodbav:"deleted" json:"deleted"`
}

// Fetches projects for a given project Id
func ListProjectsByOrgId(ctx context.Context, orgId string) ([]Project, error) {
	pk := fmt.Sprintf("ORG#%s", orgId)
	skPrefix := pk

	params := &dynamodb.QueryInput{
		TableName:              PROJECT_TABLE_NAME,
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":    &types.AttributeValueMemberS{Value: pk},
			":sk":    &types.AttributeValueMemberS{Value: skPrefix},
			":false": &types.AttributeValueMemberBOOL{Value: false},
		},
		FilterExpression: aws.String("deleted = :false"),
	}

	resp, err := awsclient.GetDynamoClient().Query(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error querying projects: %w", err)
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling projects: %w", err)
	}

	return projects, nil
}

func ListProjectsByTeamId(ctx context.Context, teamId string) ([]Project, error) {
	pk := fmt.Sprintf("TEAM#%s", teamId)
	skPrefix := pk

	params := &dynamodb.QueryInput{
		TableName:              PROJECT_TABLE_NAME,
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":    &types.AttributeValueMemberS{Value: pk},
			":sk":    &types.AttributeValueMemberS{Value: skPrefix},
			":false": &types.AttributeValueMemberBOOL{Value: false},
		},
		FilterExpression: aws.String("deleted = :false"),
	}

	resp, err := awsclient.GetDynamoClient().Query(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error querying projects: %w", err)
	}

	var projects []Project
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &projects)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling projects: %w", err)
	}

	return projects, nil
}

// GetProjectById retrieves a project by its Id
func GetProject(ctx context.Context, projectId string) (*Project, error) {
	project := &Project{Id: projectId}
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := pk

	resp, err := awsclient.GetDynamoClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: PROJECT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting project: %w", err)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &project)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling project: %w", err)
	}

	return project, nil
}

// AddProject adds a new project to the database
func AddProject(ctx context.Context, orgId string, projectId string, name string, description string) (*Project, error) {
	project := &Project{
		Id:          projectId,
		OrgId:       orgId,
		Name:        name,
		Description: description,
		Deleted:     false,
	}
	pk := fmt.Sprintf("PROJECT#%s", projectId)
	sk := pk

	item, err := attributevalue.MarshalMap(project)
	if err != nil {
		return nil, fmt.Errorf("error marshaling project: %w", err)
	}

	item["PK"] = &types.AttributeValueMemberS{Value: pk}
	item["SK"] = &types.AttributeValueMemberS{Value: sk}

	_, err = awsclient.GetDynamoClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: PROJECT_TABLE_NAME,
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("error adding project: %w", err)
	}

	return project, nil
}

// UpdateProjectDeletedStatus updates the 'deleted' status of a project
func SoftDeleteProject(ctx context.Context, projectId string) (string, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := pk

	modifiedAt := time.Now().UTC().Format(time.RFC3339)

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROJECT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("set deleted = :deleted, modifiedAt = :modifiedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":deleted":    &types.AttributeValueMemberBOOL{Value: true},
			":modifiedAt": &types.AttributeValueMemberS{Value: modifiedAt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error updating project deleted status: %w", err)
	}

	return modifiedAt, nil
}

func UpdateProjectActiveVersion(ctx context.Context, projectId string, stub string, version string) (string, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := pk

	modifiedAt := time.Now().UTC().Format(time.RFC3339)

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROJECT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET stub = :stub, version = :version, modifiedAt = :modifiedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":stub":       &types.AttributeValueMemberS{Value: stub},
			":version":    &types.AttributeValueMemberS{Value: version},
			":modifiedAt": &types.AttributeValueMemberS{Value: modifiedAt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error updating project deleted status: %w", err)
	}

	return modifiedAt, nil
}

func UpdateProject(ctx context.Context, name string, projectId string, stub string, version string) (string, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := pk

	modifiedAt := time.Now().UTC().Format(time.RFC3339)

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROJECT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET stub = :stub, version = :version, modifiedAt = :modifiedAt, name = :name"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":       &types.AttributeValueMemberS{Value: name},
			":stub":       &types.AttributeValueMemberS{Value: stub},
			":version":    &types.AttributeValueMemberS{Value: version},
			":modifiedAt": &types.AttributeValueMemberS{Value: modifiedAt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error updating project deleted status: %w", err)
	}

	return modifiedAt, nil
}
