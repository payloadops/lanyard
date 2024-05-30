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

func TestCreateBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewBranchDBClient(mockSvc)

	branch := &dal.Branch{
		Name: "prompt1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateBranch(context.Background(), "org1", "prompt1", branch)
	assert.NoError(t, err)
}

func TestGetBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewBranchDBClient(mockSvc)

	branch := dal.Branch{
		Name:      "branch1",
		Deleted:   false,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(branch)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetBranch(context.Background(), "org1", "prompt1", "branch1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "branch1", result.Name)
}

func TestDeleteBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewBranchDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteBranch(context.Background(), "org1", "prompt1", "branch1")
	assert.NoError(t, err)
}

func TestListBranchesByPrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewBranchDBClient(mockSvc)

	branch := dal.Branch{
		Name:      "branch1",
		Deleted:   false,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(branch)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListBranchesByPrompt(context.Background(), "org1", "prompt1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "branch1", result[0].Name)
}
