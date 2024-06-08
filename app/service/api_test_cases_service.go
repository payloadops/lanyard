package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/payloadops/plato/app/utils"
	"go.uber.org/zap"

	"github.com/payloadops/plato/app/dal"
	"github.com/payloadops/plato/app/openapi"
)

// TestCasesAPIService is a service that implements the logic for the TestCasesAPIServicer
// This service should implement the business logic for every endpoint for the TestCaesAPI API.
type TestCasesAPIService struct {
	promptClient   dal.PromptManager
	testCaseClient dal.TestCaseManager
	logger         *zap.Logger
}

// NewTestCasesAPIService creates a default app service
func NewTestCasesAPIService(promptClient dal.PromptManager, testCaseClient dal.TestCaseManager, logger *zap.Logger) openapi.TestCasesAPIServicer {
	return &TestCasesAPIService{promptClient: promptClient, testCaseClient: testCaseClient, logger: logger}
}

// CreateTestCase - Create a new test case for a prompt
func (s *TestCasesAPIService) CreateTestCase(ctx context.Context, projectID string, promptID string, testCaseInput openapi.TestCaseInput) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	testCase := &dal.TestCase{
		Name: testCaseInput.Name,
	}

	err = s.testCaseClient.CreateTestCase(ctx, orgID, projectID, testCase)
	if err != nil {
		s.logger.Error("failed to create test case",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(testCase.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(testCase.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	var parameters []openapi.TestCaseParameter
	// TODO create parameters

	response := openapi.TestCase{
		Id:         testCase.TestCaseID,
		Name:       testCase.Name,
		Parameters: parameters,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	return openapi.Response(http.StatusCreated, response), nil
}

// DeleteTestCase - Delete a specific test case for a prompt
func (s *TestCasesAPIService) DeleteTestCase(ctx context.Context, projectID, promptID, testCaseID string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, promptID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	// Check if the test case exists
	testCase, err := s.testCaseClient.GetTestCase(ctx, orgID, promptID, testCaseID)
	if err != nil {
		s.logger.Error("failed to get test case",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if testCase == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("testCase not found")
	}

	// TODO delete test case params

	err = s.testCaseClient.DeleteTestCase(ctx, orgID, promptID, testCaseID)
	if err != nil {
		s.logger.Error("failed to delete test case",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetTestCase - Retrieve a specific test case within a prompt
func (s *TestCasesAPIService) GetTestCase(ctx context.Context, projectID, promptID, testCaseID string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, projectID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	testCase, err := s.testCaseClient.GetTestCase(ctx, orgID, promptID, testCaseID)
	if err != nil {
		s.logger.Error("failed to get testCase",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if testCase == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("test case not found")
	}

	createdAt, err := utils.ParseTimestamp(testCase.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(testCase.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	var parameters []openapi.TestCaseParameter
	// TODO get parameters

	response := openapi.TestCase{
		Id:         testCase.TestCaseID,
		Name:       testCase.Name,
		Parameters: parameters,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}

// ListTestCases - List all test cases for a prompt
func (s *TestCasesAPIService) ListTestCases(ctx context.Context, projectID, promptID string) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, projectID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	testCases, err := s.testCaseClient.ListTestCases(ctx, orgID, promptID)
	if err != nil {
		s.logger.Error("failed to list test cases by prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	responses := make([]openapi.TestCase, len(testCases))
	for i, testCase := range testCases {
		createdAt, err := utils.ParseTimestamp(testCase.CreatedAt)
		if err != nil {
			s.logger.Error("failed to parse timestamp",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		updatedAt, err := utils.ParseTimestamp(testCase.UpdatedAt)
		if err != nil {
			s.logger.Error("failed to parse timestamp",
				zap.String("requestID", requestID),
				zap.Error(err),
			)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
		}

		responses[i] = openapi.TestCase{
			Id:        testCase.TestCaseID,
			Name:      testCase.Name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return openapi.Response(http.StatusOK, responses), nil
}

// UpdateTestCase - Update a specific test case for a prompt
func (s *TestCasesAPIService) UpdateTestCase(ctx context.Context, projectID, promptID, testCaseID string, testCaseInput openapi.TestCaseInput) (openapi.ImplResponse, error) {
	requestID := middleware.GetReqID(ctx)
	orgID, ok := ctx.Value("orgID").(string)
	if !ok || orgID == "" {
		s.logger.Error("orgID not present in context",
			zap.String("requestID", requestID),
		)
		return openapi.Response(http.StatusNotFound, nil), errors.New("org not found")
	}

	// Check if the prompt exists
	prompt, err := s.promptClient.GetPrompt(ctx, orgID, projectID, projectID)
	if err != nil {
		s.logger.Error("failed to get prompt",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if prompt == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	// Check if the testCase exists
	testCase, err := s.testCaseClient.GetTestCase(ctx, orgID, promptID, testCaseID)
	if err != nil {
		s.logger.Error("failed to get test case",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}
	if testCase == nil {
		return openapi.Response(http.StatusNotFound, nil), errors.New("prompt not found")
	}

	err = s.testCaseClient.UpdateTestCase(ctx, orgID, promptID, testCase)
	if err != nil {
		s.logger.Error("failed to update test cacse",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	createdAt, err := utils.ParseTimestamp(testCase.CreatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	updatedAt, err := utils.ParseTimestamp(testCase.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to parse timestamp",
			zap.String("requestID", requestID),
			zap.Error(err),
		)
		return openapi.Response(http.StatusInternalServerError, nil), errors.New("internal server error")
	}

	var parameters []openapi.TestCaseParameter
	// TODO update parameters

	response := openapi.TestCase{
		Id:         testCase.TestCaseID,
		Name:       testCase.Name,
		Parameters: parameters,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	return openapi.Response(http.StatusOK, response), nil
}
