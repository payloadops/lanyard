package dal_test

import (
	"context"
	"github.com/payloadops/plato/app/dal"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewProjectDBClient(mockSvc)

	project := &dal.Project{
		OrgID:       "org1",
		Name:        "Project1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.CreateProject(context.Background(), project)
	assert.NoError(t, err)
}

func TestGetProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewProjectDBClient(mockSvc)

	project := dal.Project{
		OrgID:       "org1",
		ProjectID:   "proj1",
		Name:        "Project1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(project)
	mockSvc.EXPECT().
		GetItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.GetItemOutput{Item: item}, nil)

	result, err := client.GetProject(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Project1", result.Name)
}

func TestUpdateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewProjectDBClient(mockSvc)

	project := &dal.Project{
		OrgID:       "org1",
		ProjectID:   "proj1",
		Name:        "Project1",
		Description: "Description1",
	}

	mockSvc.EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := client.UpdateProject(context.Background(), project)
	assert.NoError(t, err)
}

func TestDeleteProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewProjectDBClient(mockSvc)

	mockSvc.EXPECT().
		UpdateItem(gomock.Any(), gomock.Any()).
		Return(&dynamodb.UpdateItemOutput{}, nil)

	err := client.DeleteProject(context.Background(), "org1", "proj1")
	assert.NoError(t, err)
}

func TestListProjectsByOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockDynamoDBAPI(ctrl)
	client := dal.NewProjectDBClient(mockSvc)

	project := dal.Project{
		OrgID:       "org1",
		ProjectID:   "proj1",
		Name:        "Project1",
		Description: "Description1",
		Deleted:     false,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	item, _ := attributevalue.MarshalMap(project)
	mockSvc.EXPECT().
		Query(gomock.Any(), gomock.Any()).
		Return(&dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{item}}, nil)

	result, err := client.ListProjectsByOrganization(context.Background(), "org1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Project1", result[0].Name)
}
