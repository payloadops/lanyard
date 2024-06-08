package service

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"net/http"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/dal/mocks"
	"github.com/payloadops/plato/app/openapi"
	"github.com/stretchr/testify/assert"
)

func TestTestCaseAPIService_CreateTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockTestCaseClient := mocks.NewMockTestCaseManager(ctrl)
	service := NewTestCasesAPIService(mockPromptClient, mockTestCaseClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	testCaseInput := openapi.TestCaseInput{
		Name: "TestCase1",
		Parameters: []openapi.TestCaseParameter{
			{Key: "paramKey1", Value: "paramVal1"},
			{Key: "paramKey2", Value: "paramVal2"},
		},
	}

	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(nil)

	response, err := service.CreateTestCase(ctx, projectID, promptID, testCaseInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.NotNil(t, response.Body)
	testCase, ok := response.Body.(openapi.TestCase)
	assert.True(t, ok)
	assert.Equal(t, testCaseInput.Name, testCase.Name)
	assert.Equal(t, testCaseInput.Parameters, testCase.Parameters)
}

func TestTestCaseAPIService_DeleteTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockTestCaseClient := mocks.NewMockTestCaseManager(ctrl)
	service := NewTestCasesAPIService(mockPromptClient, mockTestCaseClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	testCaseID := "testCase1"

	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(nil)
	mockTestCaseClient.EXPECT().GetTestCase(ctx, "org1", promptID, testCaseID).Return(&dal.TestCase{}, nil)
	mockTestCaseClient.EXPECT().DeleteTestCase(ctx, "org1", promptID, testCaseID).Return(nil)

	response, err := service.DeleteTestCase(ctx, projectID, promptID, testCaseID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestTestCaseAPIService_GetTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockTestCaseClient := mocks.NewMockTestCaseManager(ctrl)
	service := NewTestCasesAPIService(mockPromptClient, mockTestCaseClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	testCaseID := "testCase1"
	testCase := openapi.TestCase{
		Id:   "testCase1",
		Name: "TestCase1",
		Parameters: []openapi.TestCaseParameter{
			{Key: "paramKey1", Value: "paramVal1"},
			{Key: "paramKey2", Value: "paramVal2"},
		},
	}

	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(nil)
	mockTestCaseClient.EXPECT().GetTestCase(ctx, "org1", promptID, testCaseID).Return(testCase, nil)

	response, err := service.GetTestCase(ctx, projectID, promptID, testCaseID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	testCase, ok := response.Body.(openapi.TestCase)
	assert.True(t, ok)
	assert.Equal(t, promptID, testCase.Id)
}

func TestTestCaseAPIService_ListTestCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockTestCaseClient := mocks.NewMockTestCaseManager(ctrl)
	service := NewTestCasesAPIService(mockPromptClient, mockTestCaseClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	testCase1 := dal.TestCase{
		TestCaseID: "testCase1",
		Name:       "TestCase1",
		UpdatedAt:  "",
		CreatedAt:  "",
	}
	testCase2 := dal.TestCase{
		TestCaseID: "testCase2",
		Name:       "TestCase2",
		UpdatedAt:  "",
		CreatedAt:  "",
	}

	testCases := []dal.TestCase{
		testCase1,
		testCase2,
	}

	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(nil)
	mockTestCaseClient.EXPECT().ListTestCases(ctx, "org1", promptID).Return(testCases, nil)

	response, err := service.ListTestCases(ctx, projectID, promptID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	listedTestCases, ok := response.Body.([]openapi.TestCase)
	assert.True(t, ok)
	assert.Equal(t, 2, len(listedTestCases))
	assert.Equal(t, "testCase1", listedTestCases[0].Id)
	assert.Equal(t, "testCase2", listedTestCases[1].Id)
}

func TestTestCaseAPIService_UpdateTestCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPromptClient := mocks.NewMockPromptManager(ctrl)
	mockTestCaseClient := mocks.NewMockTestCaseManager(ctrl)
	service := NewTestCasesAPIService(mockPromptClient, mockTestCaseClient, zap.NewNop())

	ctx := context.WithValue(context.Background(), "orgID", "org1")
	projectID := "proj1"
	promptID := "prompt1"
	testCaseID := "testCase1"
	testCaseInput := openapi.TestCaseInput{
		Name: "TestCase1",
		Parameters: []openapi.TestCaseParameter{
			{Key: "paramKey1", Value: "paramVal1"},
			{Key: "paramKey2", Value: "paramVal2"},
		},
	}
	testCase := openapi.TestCase{
		Id:   "testCase1",
		Name: "TestCase1",
		Parameters: []openapi.TestCaseParameter{
			{Key: "paramKey1", Value: "paramVal1"},
			{Key: "paramKey2", Value: "paramVal2"},
		},
	}

	mockPromptClient.EXPECT().GetPrompt(ctx, "org1", projectID, promptID).Return(nil)
	mockTestCaseClient.EXPECT().GetTestCase(ctx, "org1", promptID, testCaseID).Return(testCase, nil)
	mockTestCaseClient.EXPECT().UpdateTestCase(ctx, "org1", projectID, gomock.Any()).Return(nil)

	response, err := service.UpdateTestCase(ctx, projectID, promptID, "testCase1", testCaseInput)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	updatedPrompt, ok := response.Body.(openapi.TestCase)
	assert.True(t, ok)
	assert.Equal(t, testCaseInput.Name, updatedPrompt.Name)
}
