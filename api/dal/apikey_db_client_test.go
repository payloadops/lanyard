package dal

/*
import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockedDynamoDB is a mock of the DynamoDBAPI
type mockedDynamoDB struct {
	mock.Mock
	dynamodbiface.DynamoDBAPI
}

func (m *mockedDynamoDB) PutItemWithContext(ctx context.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, input)
	return nil, args.Error(1)
}

func (m *mockedDynamoDB) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockedDynamoDB) DeleteItemWithContext(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(ctx, input)
	return nil, args.Error(1)
}

func (m *mockedDynamoDB) ScanWithContext(ctx context.Context, input *dynamodb.ScanInput, opts ...request.Option) (*dynamodb.ScanOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

// TestCreateAPIKey tests the CreateAPIKey method
func TestCreateAPIKey(t *testing.T) {
	mockSvc := new(mockedDynamoDB)
	client := APIKeyDBClient{service: mockSvc}
	ctx := context.TODO()
	apiKey := APIKey{
		ID:        "abc123",
		ProjectID: "proj1",
		Key:       "key123",
		Scopes:    []string{"scope1", "scope2", "scope3"},
	}

	now := time.Now().UTC().Format(time.RFC3339)
	apiKey.CreatedAt = now
	apiKey.UpdatedAt = now

	av, _ := dynamodbattribute.MarshalMap(apiKey)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(nil)

	err := client.CreateAPIKey(ctx, apiKey)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestGetAPIKey tests the GetAPIKey method
func TestGetAPIKey(t *testing.T) {
	mockSvc := new(mockedDynamoDB)
	client := APIKeyDBClient{service: mockSvc}
	ctx := context.TODO()
	apiKey := APIKey{
		ID:        "abc123",
		ProjectID: "proj1",
		Key:       "key123",
		Scopes:    []string{"scope1", "scope2", "scope3"},
	}

	av, _ := dynamodbattribute.MarshalMap(apiKey)
	mockSvc.On("GetItemWithContext", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Return(&dynamodb.GetItemOutput{
		Item: av,
	}, nil)

	result, err := client.GetAPIKey(ctx, "abc123")
	assert.NoError(t, err)
	assert.Equal(t, &apiKey, result)
	mockSvc.AssertExpectations(t)
}

// TestUpdateAPIKey tests the UpdateAPIKey method
func TestUpdateAPIKey(t *testing.T) {
	mockSvc := new(mockedDynamoDB)
	client := APIKeyDBClient{service: mockSvc}
	ctx := context.TODO()
	apiKey := APIKey{
		ID:        "abc123",
		ProjectID: "proj1",
		Scopes:    []string{"scope1", "scope2", "scope3"},
		Key:       "key123Updated",
	}

	now := time.Now().UTC().Format(time.RFC3339)
	apiKey.UpdatedAt = now

	av, _ := dynamodbattribute.MarshalMap(apiKey)
	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(nil)

	err := client.UpdateAPIKey(ctx, apiKey)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestDeleteAPIKey tests the DeleteAPIKey method
func TestDeleteAPIKey(t *testing.T) {
	mockSvc := new(mockedDynamoDB)
	client := APIKeyDBClient{service: mockSvc}
	ctx := context.TODO()

	mockSvc.On("DeleteItemWithContext", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(nil)

	err := client.DeleteAPIKey(ctx, "abc123")
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestListAPIKeys tests the ListAPIKeys method
func TestListAPIKeys(t *testing.T) {
	mockSvc := new(mockedDynamoDB)
	client := APIKeyDBClient{service: mockSvc}
	ctx := context.TODO()
	apiKeys := []APIKey{
		{
			ID:        "abc123",
			ProjectID: "proj1",
			Scopes:    []string{"scope1", "scope2", "scope3"},
			Key:       "key123",
		},
	}

	av, _ := dynamodbattribute.MarshalList(apiKeys)
	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: av,
	}, nil)

	result, err := client.ListAPIKeys(ctx)
	assert.NoError(t, err)
	assert.Equal(t, apiKeys, result)
	mockSvc.AssertExpectations(t)
}
*/
