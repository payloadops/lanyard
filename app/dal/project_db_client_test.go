package dal

/*
import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *mockedDynamoDB) PutItemWithContext(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *mockedDynamoDB) GetItemWithContext(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockedDynamoDB) DeleteItemWithContext(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func (m *mockedDynamoDB) ScanWithContext(ctx aws.Context, input *dynamodb.ScanInput, opts ...request.Option) (*dynamodb.ScanOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func TestCreateProject(t *testing.T) {
	ctx := context.TODO()
	project := Project{
		ID:          "1",
		Name:        "Project A",
		Description: "Description A",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	now := time.Now().UTC().Format(time.RFC3339)
	project.CreatedAt = now
	project.UpdatedAt = now

	av, _ := dynamodbattribute.MarshalMap(project)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateProject(ctx, project)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestGetProject(t *testing.T) {
	ctx := context.TODO()
	project := Project{
		ID:          "1",
		Name:        "Project A",
		Description: "Description A",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	av, _ := dynamodbattribute.MarshalMap(project)
	mockSvc.On("GetItemWithContext", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Return(&dynamodb.GetItemOutput{
		Item: av,
	}, nil)

	result, err := client.GetProject(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, &project, result)
	mockSvc.AssertExpectations(t)
}

func TestUpdateProject(t *testing.T) {
	ctx := context.TODO()
	project := Project{
		ID:          "1",
		Name:        "Project A",
		Description: "Description A",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	now := time.Now().UTC().Format(time.RFC3339)
	project.UpdatedAt = now

	av, _ := dynamodbattribute.MarshalMap(project)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := client.UpdateProject(ctx, project)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestDeleteProject(t *testing.T) {
	ctx := context.TODO()

	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	mockSvc.On("DeleteItemWithContext", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := client.DeleteProject(ctx, "1")
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestListProjects(t *testing.T) {
	ctx := context.TODO()
	projects := []Project{
		{
			ID:          "1",
			Name:        "Project A",
			Description: "Description A",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:          "2",
			Name:        "Project B",
			Description: "Description B",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(projects[0])
	av2, _ := dynamodbattribute.MarshalMap(projects[1])
	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListProjects(ctx)
	assert.NoError(t, err)
	assert.Equal(t, projects, result)
	mockSvc.AssertExpectations(t)
}

func TestListProjectsByOrganization(t *testing.T) {
	ctx := context.TODO()
	organizationID := "org1"
	projects := []Project{
		{
			ID:          "1",
			OrgID:       organizationID,
			Name:        "Project A",
			Description: "Description A",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:          "2",
			OrgID:       organizationID,
			Name:        "Project B",
			Description: "Description B",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(projects[0])
	av2, _ := dynamodbattribute.MarshalMap(projects[1])
	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListProjectsByOrganization(ctx, organizationID)
	assert.NoError(t, err)
	assert.Equal(t, projects, result)
	mockSvc.AssertExpectations(t)
}

func TestListProjectsByTeam(t *testing.T) {
	ctx := context.TODO()
	teamID := "team1"
	projects := []Project{
		{
			ID:          "1",
			TeamID:      teamID,
			Name:        "Project A",
			Description: "Description B",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:          "2",
			TeamID:      teamID,
			Name:        "Project B",
			Description: "Description B",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(projects[0])
	av2, _ := dynamodbattribute.MarshalMap(projects[1])
	mockSvc := new(mockedDynamoDB)
	client := NewProjectDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListProjectsByTeam(ctx, teamID)
	assert.NoError(t, err)
	assert.Equal(t, projects, result)
	mockSvc.AssertExpectations(t)
}
*/
