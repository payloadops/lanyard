package dal_test

import (
	"context"
	"testing"
	"time"

	"github.com/payloadops/lanyard/app/dal"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/lanyard/app/dal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewServiceDBClient(mockSvc)

	service := &dal.Service{
		Name:        "Service1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateService(context.Background(), "org1", service)
	assert.NoError(t, err)
}

func TestGetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewServiceDBClient(mockSvc)

	service := dal.Service{
		ServiceID:   "proj1",
		Name:        "Service1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(service)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetService(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Service1", result.Name)
}

func TestUpdateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewServiceDBClient(mockSvc)

	service := &dal.Service{
		ServiceID:   "proj1",
		Name:        "Service1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, "Org#org1", input.Key["pk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Service#proj1", input.Key["sk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Service1", input.ExpressionAttributeValues[":name"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Description1", input.ExpressionAttributeValues[":description"].(*types.AttributeValueMemberS).Value)
			assert.NotEmpty(t, input.ExpressionAttributeValues[":updatedAt"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "SET #name = :name, #description = :description, #updatedAt = :updatedAt", *input.UpdateExpression)
			assert.Equal(t, "Name", input.ExpressionAttributeNames["#name"])
			assert.Equal(t, "Description", input.ExpressionAttributeNames["#description"])
			assert.Equal(t, "UpdatedAt", input.ExpressionAttributeNames["#updatedAt"])
			return &dynamodb.UpdateItemOutput{}, nil
		})

	err := client.UpdateService(context.Background(), "org1", service)
	assert.NoError(t, err)
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewServiceDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteService(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
}

func TestListServicesByOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewServiceDBClient(mockSvc)

	service := dal.Service{
		ServiceID:   "proj1",
		Name:        "Service1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(service)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListServicesByOrganization(context.Background(), "org1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Service1", result[0].Name)
}
