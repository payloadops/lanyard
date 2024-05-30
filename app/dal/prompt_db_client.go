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

//go:generate mockgen -package=mocks -destination=mocks/mock_prompt_db_client.go "github.com/payloadops/plato/app/dal" PromptManager

// PromptManager defines the operations available for managing prompts.
type PromptManager interface {
	CreatePrompt(ctx context.Context, orgID, projectID string, prompt *Prompt) error
	GetPrompt(ctx context.Context, orgID, projectID, promptID string) (*Prompt, error)
	UpdatePrompt(ctx context.Context, orgID, projectID string, prompt *Prompt) error
	DeletePrompt(ctx context.Context, orgID, projectID, promptID string) error
	ListPromptsByProject(ctx context.Context, orgID, projectID string) ([]Prompt, error)
}

// Ensure PromptDBClient implements the PromptManager interface
var _ PromptManager = &PromptDBClient{}

// Prompt represents a prompt in the system.
type Prompt struct {
	PromptID    string `json:"promptId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deleted     bool   `json:"deleted"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// PromptDBClient is a client for interacting with DynamoDB for prompt-related operations.
type PromptDBClient struct {
	service DynamoDBAPI
}

// NewPromptDBClient creates a new PromptDBClient.
func NewPromptDBClient(service DynamoDBAPI) *PromptDBClient {
	return &PromptDBClient{
		service: service,
	}
}

// createProjectCompositeKeys generates the partition key (pk) and sort key (SK) for a prompt.
func createPromptCompositeKeys(orgID, projectID, promptID string) (string, string) {
	return "Org#" + orgID + "Project#" + projectID, "Prompt#" + promptID
}

// CreatePrompt creates a new prompt in the DynamoDB table.
func (d *PromptDBClient) CreatePrompt(ctx context.Context, orgID, projectID string, prompt *Prompt) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	prompt.PromptID = ksuid
	pk, sk := createPromptCompositeKeys(orgID, projectID, prompt.PromptID)

	now := time.Now().UTC().Format(time.RFC3339)
	prompt.CreatedAt = now
	prompt.UpdatedAt = now

	av, err := attributevalue.MarshalMap(prompt)
	if err != nil {
		return fmt.Errorf("failed to marshal prompt: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Prompts"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetPrompt retrieves a prompt by orgID, project ID, and prompt ID from the DynamoDB table.
func (d *PromptDBClient) GetPrompt(ctx context.Context, orgID, projectID, promptID string) (*Prompt, error) {
	pk, sk := createPromptCompositeKeys(orgID, projectID, promptID)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Prompts"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := d.service.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var prompt Prompt
	err = attributevalue.UnmarshalMap(result.Item, &prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	if prompt.Deleted {
		return nil, nil
	}

	return &prompt, nil
}

// UpdatePrompt updates an existing prompt in the DynamoDB table.
func (d *PromptDBClient) UpdatePrompt(ctx context.Context, orgID string, projectID string, prompt *Prompt) error {
	pk, sk := createPromptCompositeKeys(orgID, projectID, prompt.PromptID)
	prompt.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(prompt)
	if err != nil {
		return fmt.Errorf("failed to marshal prompt: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String("Prompts"),
		Item:                item,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(SK)"),
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeletePrompt marks a prompt as deleted by org ID, project ID, and prompt ID in the DynamoDB table.
func (d *PromptDBClient) DeletePrompt(ctx context.Context, orgID, projectID, promptID string) error {
	pk, sk := createPromptCompositeKeys(orgID, projectID, promptID)
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
		TableName: aws.String("Prompts"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		AttributeUpdates:    update,
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(SK)"),
	}

	_, err := d.service.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in DynamoDB: %v", err)
	}

	return nil
}

// ListPromptsByProject retrieves all prompts belonging to a specific project from the DynamoDB table.
func (d *PromptDBClient) ListPromptsByProject(ctx context.Context, orgID string, projectID string) ([]Prompt, error) {
	pk, _ := createPromptCompositeKeys(orgID, projectID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("Prompts"),
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

	var prompts []Prompt
	err = attributevalue.UnmarshalListOfMaps(result.Items, &prompts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return prompts, nil
}
