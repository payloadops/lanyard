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

func TestCreateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewAPIKeyDBClient(mockSvc)

	apiKey := &dal.APIKey{
		ProjectID: "proj1",
		Secret:    "key1",
		OrgID:     "org1",
		Scopes:    []string{"scope1", "scope2"},
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateAPIKey(context.Background(), apiKey)
	assert.NoError(t, err)
}

func TestGetAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewAPIKeyDBClient(mockSvc)

	apiKey := dal.APIKey{
		ProjectID: "proj1",
		APIKeyID:  "key1",
		Secret:    "key1",
		OrgID:     "org1",
		Scopes:    []string{"scope1", "scope2"},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(apiKey)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetAPIKey(context.Background(), "key1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "key1", result.Secret)
}

func TestUpdateAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewAPIKeyDBClient(mockSvc)

	apiKey := &dal.APIKey{
		ProjectID: "proj1",
		APIKeyID:  "key1",
		OrgID:     "org1",
		Secret:    "key1",
		Scopes:    []string{"scope1", "scope2"},
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.UpdateAPIKey(context.Background(), apiKey)
	assert.NoError(t, err)
}

func TestDeleteAPIKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewAPIKeyDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteAPIKey(context.Background(), "org1", "proj1", "key1")
	assert.NoError(t, err)
}

func TestListAPIKeysByProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewAPIKeyDBClient(mockSvc)

	apiKey := dal.APIKey{
		ProjectID: "proj1",
		APIKeyID:  "key1",
		OrgID:     "org1",
		Secret:    "key1",
		Scopes:    []string{"scope1", "scope2"},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(apiKey)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListAPIKeysByProject(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "key1", result[0].Secret)
}
