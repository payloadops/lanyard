package dal

/*
import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestCreateUser(t *testing.T) {
	ctx := context.TODO()
	user := User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockSvc := new(mockedDynamoDB)
	client := UserDBClient{service: mockSvc}

	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(nil)

	err := client.CreateUser(ctx, user)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	ctx := context.TODO()
	user := User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockSvc := new(mockedDynamoDB)
	client := UserDBClient{service: mockSvc}

	mockSvc.On("GetItemWithContext", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Return(&dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"id":    {S: aws.String("1")},
			"name":  {S: aws.String("John Doe")},
			"email": {S: aws.String("john@example.com")},
		},
	}, nil)

	result, err := client.GetUser(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, &user, result)
	mockSvc.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.TODO()
	user := User{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockSvc := new(mockedDynamoDB)
	client := UserDBClient{service: mockSvc}

	mockSvc.On("PutItemWithContext", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Return(nil)

	err := client.UpdateUser(ctx, user)
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	ctx := context.TODO()

	mockSvc := new(mockedDynamoDB)
	client := UserDBClient{service: mockSvc}

	mockSvc.On("DeleteItemWithContext", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).Return(nil)

	err := client.DeleteUser(ctx, "1")
	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	ctx := context.TODO()
	users := []User{
		{
			ID:    "1",
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			ID:    "2",
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	}

	mockSvc := new(mockedDynamoDB)
	client := UserDBClient{service: mockSvc}

	mockSvc.On("ScanWithContext", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":    {S: aws.String("1")},
				"name":  {S: aws.String("John Doe")},
				"email": {S: aws.String("john@example.com")},
			},
			{
				"id":    {S: aws.String("2")},
				"name":  {S: aws.String("Jane Doe")},
				"email": {S: aws.String("jane@example.com")},
			},
		},
	}, nil)

	result, err := client.ListUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, users, result)
	mockSvc.AssertExpectations(t)
}
*/
