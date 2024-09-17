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

func TestCreateTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTierDBClient(mockSvc)

	Tier := &dal.Tier{
		TierID:              "12342341234",
		Name:                "Steve",
		DefaultRequestLimit: 1000000,
		OveragePrice:        1,
		Interval:            1,
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateTier(context.Background(), "org1", "serv1", Tier)
	assert.NoError(t, err)
}

func TestGetTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTierDBClient(mockSvc)

	Tier := dal.Tier{
		TierID:              "12342341234",
		Name:                "Steve",
		DefaultRequestLimit: 1000000,
		OveragePrice:        1,
		Interval:            1,
	}

	item, _ := attributevalue.MarshalMap(Tier)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetTier(context.Background(), "org1", "serv1", "Tier1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Tier1", result.TierID)
}

func TestUpdateTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTierDBClient(mockSvc)

	Tier := &dal.Tier{
		TierID:              "12342341234",
		Name:                "Steve",
		DefaultRequestLimit: 1000000,
		OveragePrice:        1,
		Interval:            1,
	}

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, "Service#serv1", input.Key["pk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Tier#Tier1", input.Key["sk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "12342341234", input.ExpressionAttributeValues[":externalId"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "1000000", input.ExpressionAttributeValues[":monthlyRequestLimit"].(*types.AttributeValueMemberN).Value)
			assert.Equal(t, "SET #externalId = :externalId, #monthlyRequestLimit = :monthlyRequestLimit", *input.UpdateExpression)
			assert.Equal(t, "ExternalId", input.ExpressionAttributeNames["#externalId"])
			assert.Equal(t, "MonthlyRequestLimit", input.ExpressionAttributeNames["#monthlyRequestLimit"])
			return &dynamodb.UpdateItemOutput{}, nil
		})

	err := client.UpdateTier(context.Background(), "org1", "serv1", Tier)
	assert.NoError(t, err)
}

func TestDeleteTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTierDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteTier(context.Background(), "org1", "serv1", "Tier1")
	assert.NoError(t, err)
}

func TestListTiersByServiceanization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewTierDBClient(mockSvc)

	Tier := dal.Tier{
		TierID:              "12342341234",
		Name:                "Steve",
		DefaultRequestLimit: 1000000,
		OveragePrice:        1,
		Interval:            1,
	}

	item, _ := attributevalue.MarshalMap(Tier)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListTiers(context.Background(), "org1", "serv1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Tier1", result[0].TierID)
}
