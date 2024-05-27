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

func TestCreateCommit(t *testing.T) {
	ctx := context.TODO()
	commit := Commit{
		ID:       "1",
		BranchID: "branch1",
		Content:  "Initial commit",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewCommitDBClient(mockSvc)

	now := time.Now().UTC().Format(time.RFC3339)
	commit.CreatedAt = now

	av, _ := dynamodbattribute.MarshalMap(commit)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateCommit(ctx, commit)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestGetCommit(t *testing.T) {
	ctx := context.TODO()
	commit := Commit{
		ID:       "1",
		BranchID: "branch1",
		Content:  "Initial commit",
	}

	mockSvc := new(mockedDynamoDB)
	client := NewCommitDBClient(mockSvc)

	av, _ := dynamodbattribute.MarshalMap(commit)
	mockSvc.On("GetItemWithContext", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Return(&dynamodb.GetItemOutput{
		Item: av,
	}, nil)

	result, err := client.GetCommit(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, &commit, result)
	mockSvc.AssertExpectations(t)
}

func TestListCommits(t *testing.T) {
	ctx := context.TODO()
	commits := []Commit{
		{
			ID:        "1",
			BranchID:  "branch1",
			Content:   "Initial commit",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:        "2",
			BranchID:  "branch1",
			Content:   "Second commit",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(commits[0])
	av2, _ := dynamodbattribute.MarshalMap(commits[1])
	mockSvc := new(mockedDynamoDB)
	client := NewCommitDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListCommits(ctx)
	assert.NoError(t, err)
	assert.Equal(t, commits, result)
	mockSvc.AssertExpectations(t)
}

func TestListCommitsByBranch(t *testing.T) {
	ctx := context.TODO()
	branchID := "branch1"
	commits := []Commit{
		{
			ID:        "1",
			BranchID:  branchID,
			Content:   "Initial commit",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
		{
			ID:        "2",
			BranchID:  branchID,
			Content:   "Second commit",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	av1, _ := dynamodbattribute.MarshalMap(commits[0])
	av2, _ := dynamodbattribute.MarshalMap(commits[1])
	mockSvc := new(mockedDynamoDB)
	client := NewCommitDBClient(mockSvc)

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{av1, av2},
	}, nil)

	result, err := client.ListCommitsByBranch(ctx, branchID)
	assert.NoError(t, err)
	assert.Equal(t, commits, result)
	mockSvc.AssertExpectations(t)
}
*/
