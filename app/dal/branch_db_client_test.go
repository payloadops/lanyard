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

func TestCreateBranch(t *testing.T) {
	ctx := context.TODO()
	branch := Branch{
		ID:       "1",
		PromptID: "prompt1",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewBranchDBClient(mockSvc)

	now := time.Now().UTC().Format(time.RFC3339)
	branch.CreatedAt = now

	av, _ := dynamodbattribute.MarshalMap(branch)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateBranch(ctx, branch)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestGetBranch(t *testing.T) {
	ctx := context.TODO()
	branch := Branch{
		ID:       "1",
		PromptID: "prompt1",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewBranchDBClient(mockSvc)

	av, _ := dynamodbattribute.MarshalMap(branch)
	mockSvc.On("GetItemWithContext", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Return(&dynamodb.GetItemOutput{
		Item: av,
	}, nil)

	result, err := client.GetBranch(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, &branch, result)
	mockSvc.AssertExpectations(t)
}

func TestDeleteBranch(t *testing.T) {
	ctx := context.TODO()

	mockSvc := new(mockedDynamoDB)
	client := NewBranchDBClient(mockSvc)

	mockSvc.On("DeleteItemWithContext", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := client.DeleteBranch(ctx, "1")
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestListBranches(t *testing.T) {
	ctx := context.TODO()
	branches := []Branch{
		{
			ID:        "1",
			PromptID:  "prompt1",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:        "2",
			PromptID:  "prompt1",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(branches[0])
	av2, _ := dynamodbattribute.MarshalMap(branches[1])
	mockSvc := new(mockedDynamoDB)
	client := NewBranchDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListBranches(ctx)
	assert.NoError(t, err)
	assert.Equal(t, branches, result)
	mockSvc.AssertExpectations(t)
}

func TestListBranchesByPrompt(t *testing.T) {
	ctx := context.TODO()
	promptID := "prompt1"
	branches := []Branch{
		{
			ID:        "1",
			PromptID:  promptID,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:        "2",
			PromptID:  promptID,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(branches[0])
	av2, _ := dynamodbattribute.MarshalMap(branches[1])
	mockSvc := new(mockedDynamoDB)
	client := NewBranchDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListBranchesByPrompt(ctx, promptID)
	assert.NoError(t, err)
	assert.Equal(t, branches, result)
	mockSvc.AssertExpectations(t)
}
*/
