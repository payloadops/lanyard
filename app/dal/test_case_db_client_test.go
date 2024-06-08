package dal_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTestCaseDBClient(mockSvc)

	testCase := &dal.TestCase{
		Name: "Prompt1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateTestCase(context.Background(), "org1", "prompt1", testCase)
	assert.NoError(t, err)
}

func TestGetTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTestCaseDBClient(mockSvc)

	testCase := dal.TestCase{
		TestCaseID: "testCase1",
		Name:       "TestCase1",
		Deleted:    false,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(testCase)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetTestCase(context.Background(), "org1", "prompt1", "testCase1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TestCase1", result.Name)
}

func TestUpdateTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTestCaseDBClient(mockSvc)

	testCase := &dal.TestCase{
		TestCaseID: "testCase1",
		Name:       "TestCase1",
		Deleted:    false,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, "Org#org1Prompt#prompt1", input.Key["pk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "TestCase#testCase1", input.Key["sk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "TestCase1", input.ExpressionAttributeValues[":name"].(*types.AttributeValueMemberS).Value)
			assert.NotEmpty(t, input.ExpressionAttributeValues[":updatedAt"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "SET #name = :name, #description = :description, #updatedAt = :updatedAt", *input.UpdateExpression)
			assert.Equal(t, "Name", input.ExpressionAttributeNames["#name"])
			assert.Equal(t, "Description", input.ExpressionAttributeNames["#description"])
			assert.Equal(t, "UpdatedAt", input.ExpressionAttributeNames["#updatedAt"])
			return &dynamodb.UpdateItemOutput{}, nil
		})

	err := client.UpdateTestCase(context.Background(), "org1", "prompt1", testCase)
	assert.NoError(t, err)
}

func TestDeleteTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTestCaseDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteTestCase(context.Background(), "org1", "prompt1", "testCase1")
	assert.NoError(t, err)
}

func TestListTestCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTestCaseDBClient(mockSvc)

	testCase := dal.TestCase{
		TestCaseID: "testCase1",
		Name:       "TestCase1",
		Deleted:    false,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(testCase)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListTestCases(context.Background(), "prompt1", "testCase1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "TestCase1", result[0].Name)
}
