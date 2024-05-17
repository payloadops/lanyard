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

func TestCreatePrompt(t *testing.T) {
	ctx := context.TODO()
	prompt := Prompt{
		ID:        "1",
		ProjectID: "proj1",
		Content:   "Prompt Content",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewPromptDBClient(mockSvc)

	now := time.Now().UTC().Format(time.RFC3339)
	prompt.CreatedAt = now
	prompt.UpdatedAt = now

	av, _ := dynamodbattribute.MarshalMap(prompt)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreatePrompt(ctx, prompt)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestGetPrompt(t *testing.T) {
	ctx := context.TODO()
	prompt := Prompt{
		ID:        "1",
		ProjectID: "proj1",
		Content:   "Prompt Content",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewPromptDBClient(mockSvc)

	av, _ := dynamodbattribute.MarshalMap(prompt)
	mockSvc.On("GetItemWithContext", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Return(&dynamodb.GetItemOutput{
		Item: av,
	}, nil)

	result, err := client.GetPrompt(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, &prompt, result)
	mockSvc.AssertExpectations(t)
}

func TestUpdatePrompt(t *testing.T) {
	ctx := context.TODO()
	prompt := Prompt{
		ID:        "1",
		ProjectID: "proj1",
		Content:   "Updated Prompt Content",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewPromptDBClient(mockSvc)

	now := time.Now().UTC().Format(time.RFC3339)
	prompt.UpdatedAt = now

	av, _ := dynamodbattribute.MarshalMap(prompt)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := client.UpdatePrompt(ctx, prompt)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestDeletePrompt(t *testing.T) {
	ctx := context.TODO()

	mockSvc := new(mockedDynamoDB)
	client := NewPromptDBClient(mockSvc)

	mockSvc.On("DeleteItemWithContext", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := client.DeletePrompt(ctx, "1")
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestListPrompts(t *testing.T) {
	ctx := context.TODO()
	prompts := []Prompt{
		{
			ID:        "1",
			ProjectID: "proj1",
			Content:   "Prompt Content A",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:        "2",
			ProjectID: "proj1",
			Content:   "Prompt Content B",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(prompts[0])
	av2, _ := dynamodbattribute.MarshalMap(prompts[1])
	mockSvc := new(mockedDynamoDB)
	client := NewPromptDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListPrompts(ctx)
	assert.NoError(t, err)
	assert.Equal(t, prompts, result)
	mockSvc.AssertExpectations(t)
}

func TestListPromptsByProject(t *testing.T) {
	ctx := context.TODO()
	projectID := "proj1"
	prompts := []Prompt{
		{
			ID:        "1",
			ProjectID: projectID,
			Content:   "Prompt Content A",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:        "2",
			ProjectID: projectID,
			Content:   "Prompt Content B",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(prompts[0])
	av2, _ := dynamodbattribute.MarshalMap(prompts[1])
	mockSvc := new(mockedDynamoDB)
	client := NewPromptDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListPromptsByProject(ctx, projectID)
	assert.NoError(t, err)
	assert.Equal(t, prompts, result)
	mockSvc.AssertExpectations(t)
}
*/
