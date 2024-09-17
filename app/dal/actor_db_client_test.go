package dal_test

import (
	"context"
	"testing"

	"github.com/payloadops/lanyard/app/dal"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/lanyard/app/dal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := &dal.Actor{
		ExternalID:          "12342341234",
		MonthlyRequestLimit: 1000000,
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateActor(context.Background(), "serv1", actor)
	assert.NoError(t, err)
}

func TestGetActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := dal.Actor{
		ActorID:             "actor1",
		ExternalID:          "12342341234",
		MonthlyRequestLimit: 1000000,
		Deleted:             false,
	}

	item, _ := attributevalue.MarshalMap(actor)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetActor(context.Background(), "serv1", "proj1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "actor1", result.ActorID)
}

func TestUpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := &dal.Actor{
		ActorID:             "actor1",
		ExternalID:          "12342341234",
		MonthlyRequestLimit: 1000000,
		Deleted:             false,
	}

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, "Service#serv1", input.Key["pk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Actor#actor1", input.Key["sk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "12342341234", input.ExpressionAttributeValues[":externalId"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "1000000", input.ExpressionAttributeValues[":monthlyRequestLimit"].(*types.AttributeValueMemberN).Value)
			assert.Equal(t, "SET #externalId = :externalId, #monthlyRequestLimit = :monthlyRequestLimit", *input.UpdateExpression)
			assert.Equal(t, "ExternalId", input.ExpressionAttributeNames["#externalId"])
			assert.Equal(t, "MonthlyRequestLimit", input.ExpressionAttributeNames["#monthlyRequestLimit"])
			return &dynamodb.UpdateItemOutput{}, nil
		})

	err := client.UpdateActor(context.Background(), "serv1", actor)
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

	err := client.DeleteActor(context.Background(), "serv1", "proj1")
	assert.NoError(t, err)
}

func TestListActorsByServiceanization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewActorDBClient(mockSvc)

	actor := dal.Actor{
		ActorID:             "actor1",
		Deleted:             false,
		ExternalID:          "12342341234",
		MonthlyRequestLimit: 1000000,
	}

	item, _ := attributevalue.MarshalMap(actor)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListActors(context.Background(), "serv1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "actor1", result[0].ActorID)
}
