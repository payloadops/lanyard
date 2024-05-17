package dal

/*
import (
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type mockedDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
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

func TestCreateOrganization(t *testing.T) {
	mockSvc := &MockDynamoDBAPI{
		putItemOutput: &dynamodb.PutItemOutput{},
	}

	client := NewOrgDBClient(mockSvc)

	org := Organization{
		ID:          "1",
		Name:        "Test Org",
		Description: "Test Description",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	err := client.CreateOrganization(context.Background(), org)
	assert.NoError(t, err)
}

func TestGetOrganization(t *testing.T) {
	org := Organization{
		ID:          "1",
		Name:        "Test Org",
		Description: "Test Description",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	av, _ := dynamodbattribute.MarshalMap(org)
	mockSvc := &MockDynamoDBAPI{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: av,
		},
	}

	client := NewOrgDBClient(mockSvc)

	result, err := client.GetOrganization(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, &org, result)
}

func TestUpdateOrganization(t *testing.T) {
	mockSvc := &MockDynamoDBAPI{
		putItemOutput: &dynamodb.PutItemOutput{},
	}

	client := NewOrgDBClient(mockSvc)

	org := Organization{
		ID:          "1",
		Name:        "Updated Org",
		Description: "Test Description",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	err := client.UpdateOrganization(context.Background(), org)
	assert.NoError(t, err)
}

func TestDeleteOrganization(t *testing.T) {
	mockSvc := &MockDynamoDBAPI{
		deleteItemOutput: &dynamodb.DeleteItemOutput{},
	}

	client := NewOrgDBClient(mockSvc)

	err := client.DeleteOrganization(context.Background(), "1")
	assert.NoError(t, err)
}

func TestListOrganizations(t *testing.T) {
	org := Organization{
		ID:          "1",
		Name:        "Test Org",
		Description: "Test Description",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	av, _ := dynamodbattribute.MarshalMap(org)
	mockSvc := &MockDynamoDBAPI{
		scanOutput: &dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{av},
		},
	}

	client := NewOrgDBClient(mockSvc)

	result, err := client.ListOrganizations(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []Organization{org}, result)
}
*/
