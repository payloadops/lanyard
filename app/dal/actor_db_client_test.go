package dal_test

import (
	"context"
	"testing"
	"time"

	"github.com/payloadops/plato/app/dal"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := &dal.Actor{
		Name:        "Actor1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateActor(context.Background(), "org1", actor)
	assert.NoError(t, err)
}

func TestGetActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := dal.Actor{
		ActorID:     "proj1",
		Name:        "Actor1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(actor)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetActor(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Actor1", result.Name)
}

func TestUpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := &dal.Actor{
		ActorID:     "proj1",
		Name:        "Actor1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, "Org#org1", input.Key["pk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Actor#proj1", input.Key["sk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Actor1", input.ExpressionAttributeValues[":name"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Description1", input.ExpressionAttributeValues[":description"].(*types.AttributeValueMemberS).Value)
			assert.NotEmpty(t, input.ExpressionAttributeValues[":updatedAt"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "SET #name = :name, #description = :description, #updatedAt = :updatedAt", *input.UpdateExpression)
			assert.Equal(t, "Name", input.ExpressionAttributeNames["#name"])
			assert.Equal(t, "Description", input.ExpressionAttributeNames["#description"])
			assert.Equal(t, "UpdatedAt", input.ExpressionAttributeNames["#updatedAt"])
			return &dynamodb.UpdateItemOutput{}, nil
		})

	err := client.UpdateActor(context.Background(), "org1", actor)
	assert.NoError(t, err)
}

func TestDeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteActor(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
}

func TestListActorsByOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := dal.Actor{
		ActorID:     "proj1",
		Name:        "Actor1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(actor)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListActors(context.Background(), "org1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Actor1", result[0].Name)
}
