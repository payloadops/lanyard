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

func TestCreatePrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewPromptDBClient(mockSvc)

	prompt := &dal.Prompt{
		Name:        "Prompt1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreatePrompt(context.Background(), "org1", "project1", prompt)
	assert.NoError(t, err)
}

func TestGetPrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewPromptDBClient(mockSvc)

	prompt := dal.Prompt{
		PromptID:    "prompt1",
		Name:        "Prompt1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(prompt)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetPrompt(context.Background(), "org1", "project1", "prompt1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Prompt1", result.Name)
}

func TestUpdatePrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewPromptDBClient(mockSvc)

	prompt := &dal.Prompt{
		PromptID:    "prompt1",
		Name:        "Prompt1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *dynamodb.UpdateItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
			assert.Equal(t, "Org#org1Project#project1", input.Key["pk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Prompt#prompt1", input.Key["sk"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Prompt1", input.ExpressionAttributeValues[":name"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "Description1", input.ExpressionAttributeValues[":description"].(*types.AttributeValueMemberS).Value)
			assert.NotEmpty(t, input.ExpressionAttributeValues[":updatedAt"].(*types.AttributeValueMemberS).Value)
			assert.Equal(t, "SET #name = :name, #description = :description, #updatedAt = :updatedAt", *input.UpdateExpression)
			assert.Equal(t, "Name", input.ExpressionAttributeNames["#name"])
			assert.Equal(t, "Description", input.ExpressionAttributeNames["#description"])
			assert.Equal(t, "UpdatedAt", input.ExpressionAttributeNames["#updatedAt"])
			return &dynamodb.UpdateItemOutput{}, nil
		})

	err := client.UpdatePrompt(context.Background(), "org1", "project1", prompt)
	assert.NoError(t, err)
}

func TestDeletePrompt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewPromptDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeletePrompt(context.Background(), "org1", "project1", "prompt1")
	assert.NoError(t, err)
}

func TestListPromptsByProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewPromptDBClient(mockSvc)

	prompt := dal.Prompt{
		PromptID:    "prompt1",
		Name:        "Prompt1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(prompt)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListPromptsByProject(context.Background(), "org1", "project1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Prompt1", result[0].Name)
}
