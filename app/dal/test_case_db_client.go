package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/payloadops/plato/app/utils"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_test_case_db_client.go "github.com/payloadops/plato/app/dal" TestCaseManager

// TestCaseManager defines the operations available for managing test cases.
type TestCaseManager interface {
	CreateTestCase(ctx context.Context, orgID, promptID string, testCase *TestCase) error
	GetTestCase(ctx context.Context, orgID, promptID, testCaseID string) (*TestCase, error)
	UpdateTestCase(ctx context.Context, orgID, promptID string, testCase *TestCase) error
	DeleteTestCase(ctx context.Context, orgID, promptID, testCaseID string) error
	ListTestCases(ctx context.Context, orgID, promptID string) ([]TestCase, error)

	CreateTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID string, parameter *Parameter) error
	GetTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID, parameterID string) (*TestCase, error)
	UpdateTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID, parameter *Parameter) error
	DeleteTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID, parameterID string) error
	ListTestCaseParameters(ctx context.Context, orgID, promptID, testCaseID string) ([]TestCase, error)
}

// Ensure TestCaseDBClient implements the TestCaseManager interface
var _ TestCaseManager = &TestCaseDBClient{}

// TestCase represents a test case in the system.
type TestCase struct {
	TestCaseID string `json:"testCaseId"`
	Name       string `json:"name"`
	Deleted    bool   `json:"deleted"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// Parameter represents a test case parameter in the system.
type Parameter struct {
	ParameterID string `json:"testCaseId"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

// TestCaseDBClient is a client for interacting with DynamoDB for prompt-related operations.
type TestCaseDBClient struct {
	service DynamoDBAPI
}

// NewTestCaseDBClient creates a new TestCaseDBClient.
func NewTestCaseDBClient(service DynamoDBAPI) *TestCaseDBClient {
	return &TestCaseDBClient{
		service: service,
	}
}

// createProjectCompositeKeys generates the partition key (pk) and sort key (sk) for a prompt.
func createTestCaseCompositeKeys(orgID, promptID, testCaseID string) (string, string) {
	return fmt.Sprintf("Org#%sPrompt#%s", orgID, promptID), fmt.Sprintf("TestCase#%s", testCaseID)
}

func createParameterCompositeKeys(orgID, promptID, testCaseID, parameterID string) (string, string) {
	return fmt.Sprintf("Org#%sPrompt#%sTestCase#%s", orgID, promptID, testCaseID), fmt.Sprintf("Parameter#%s", parameterID)
}

// CreateTestCase creates a new test case in the DynamoDB table.
func (d *TestCaseDBClient) CreateTestCase(ctx context.Context, orgID, promptID string, testCase *TestCase) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	testCase.TestCaseID = ksuid
	pk, sk := createTestCaseCompositeKeys(orgID, promptID, testCase.TestCaseID)

	now := time.Now().UTC().Format(time.RFC3339)
	testCase.CreatedAt = now
	testCase.UpdatedAt = now

	av, err := attributevalue.MarshalMap(testCase)
	if err != nil {
		return fmt.Errorf("failed to marshal test case: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("TestCases"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetTestCase retrieves a test case by orgID, prompt ID, and test case ID from the DynamoDB table.
func (d *TestCaseDBClient) GetTestCase(ctx context.Context, orgID, promptID, testCaseID string) (*TestCase, error) {
	pk, sk := createTestCaseCompositeKeys(orgID, promptID, testCaseID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("TestCases"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := d.service.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var testCase TestCase
	err = attributevalue.UnmarshalMap(result.Item, &testCase)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	if testCase.Deleted {
		return nil, nil
	}

	return &testCase, nil
}

// UpdateTestCase updates the name, and updatedAt fields of an existing test case in the DynamoDB table.
func (d *TestCaseDBClient) UpdateTestCase(ctx context.Context, orgID, promptID string, testCase *TestCase) error {
	pk, sk := createTestCaseCompositeKeys(orgID, promptID, testCase.TestCaseID)
	testCase.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	updateExpr := "SET #name = :name, #description = :description, #updatedAt = :updatedAt"
	exprAttrNames := map[string]string{
		"#name":        "Name",
		"#description": "Description",
		"#updatedAt":   "UpdatedAt",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":name":      &types.AttributeValueMemberS{Value: testCase.Name},
		":updatedAt": &types.AttributeValueMemberS{Value: testCase.UpdatedAt},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("TestCases"),
		Key:                       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}, "sk": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ConditionExpression:       aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteTestCase marks a testCase as deleted by org ID, prompt ID, and test case ID in the DynamoDB table.
func (d *TestCaseDBClient) DeleteTestCase(ctx context.Context, orgID, promptID, testCaseID string) error {
	pk, sk := createTestCaseCompositeKeys(orgID, promptID, testCaseID)
	now := time.Now().UTC().Format(time.RFC3339)

	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
		"UpdatedAt": {
			Value:  &types.AttributeValueMemberS{Value: now},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("TestCases"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListTestCasesByProject retrieves all test cases belonging to a specific prompt from the DynamoDB table.
func (d *TestCaseDBClient) ListTestCasesByProject(ctx context.Context, orgID string, promptID string) ([]TestCase, error) {
	pk, _ := createTestCaseCompositeKeys(orgID, promptID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("TestCases"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.service.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var testCases []TestCase
	err = attributevalue.UnmarshalListOfMaps(result.Items, &testCases)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	results := []TestCase{}
	for _, testCase := range testCases {
		if testCase.Deleted {
			continue
		}
		results = append(results, testCase)
	}

	return testCases, nil
}

func (d *TestCaseDBClient) CreateTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID string, parameter *Parameter) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	parameter.ParameterID = ksuid
	pk, sk := createParameterCompositeKeys(orgID, promptID, testCaseID, parameter.ParameterID)

	av, err := attributevalue.MarshalMap(parameter)
	if err != nil {
		return fmt.Errorf("failed to marshal test case parameter: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"sk": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("TestCases"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetTestCase retrieves a test case by orgID, prompt ID, and test case ID from the DynamoDB table.
func (d *TestCaseDBClient) GetTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID, parameterID string) (*TestCase, error) {
	pk, sk := createParameterCompositeKeys(orgID, promptID, testCaseID, parameterID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("TestCases"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := d.service.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var testCase TestCase
	err = attributevalue.UnmarshalMap(result.Item, &testCase)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	if testCase.Deleted {
		return nil, nil
	}

	return &testCase, nil
}

// UpdateTestCase updates the name, and updatedAt fields of an existing test case in the DynamoDB table.
func (d *TestCaseDBClient) UpdateTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID string, parameter *Parameter) error {
	pk, sk := createParameterCompositeKeys(orgID, promptID, testCaseID, parameter.ParameterID)

	updateExpr := "SET #name = :name, #description = :description, #updatedAt = :updatedAt"
	exprAttrNames := map[string]string{
		"#key":   "Key",
		"#value": "Value",
	}

	exprAttrValues := map[string]types.AttributeValue{
		":key":   &types.AttributeValueMemberS{Value: parameter.Key},
		":value": &types.AttributeValueMemberS{Value: parameter.Value},
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String("TestCases"),
		Key:                       map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}, "sk": &types.AttributeValueMemberS{Value: sk}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  exprAttrNames,
		ExpressionAttributeValues: exprAttrValues,
		ConditionExpression:       aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteTestCase marks a testCase as deleted by org ID, prompt ID, and test case ID in the DynamoDB table.
func (d *TestCaseDBClient) DeleteTestCaseParameter(ctx context.Context, orgID, promptID, testCaseID, parameterID string) error {
	pk, sk := createParameterCompositeKeys(orgID, promptID, testCaseID, parameterID)
	now := time.Now().UTC().Format(time.RFC3339)

	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
		"UpdatedAt": {
			Value:  &types.AttributeValueMemberS{Value: now},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("TestCases"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListTestCasesByProject retrieves all test cases belonging to a specific prompt from the DynamoDB table.
func (d *TestCaseDBClient) ListTestCaseParameters(ctx context.Context, orgID, promptID, testCaseID string) ([]TestCase, error) {
	pk, _ := createParameterCompositeKeys(orgID, promptID, testCaseID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("TestCases"),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: pk,
			},
		},
	}

	result, err := d.service.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var testCases []TestCase
	err = attributevalue.UnmarshalListOfMaps(result.Items, &testCases)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	results := []TestCase{}
	for _, testCase := range testCases {
		if testCase.Deleted {
			continue
		}
		results = append(results, testCase)
	}

	return testCases, nil
}
