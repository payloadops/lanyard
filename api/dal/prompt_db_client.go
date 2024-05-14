package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_prompt_db_client.go "github.com/payloadops/plato/api/dal" PromptManager

// PromptManager defines the operations available for managing prompts.
type PromptManager interface {
	CreatePrompt(ctx context.Context, prompt Prompt) error
	GetPrompt(ctx context.Context, id string) (*Prompt, error)
	UpdatePrompt(ctx context.Context, prompt Prompt) error
	DeletePrompt(ctx context.Context, id string) error
	ListPrompts(ctx context.Context) ([]Prompt, error)
	ListPromptsByProject(ctx context.Context, projectID string) ([]Prompt, error)
}

// Ensure PromptDBClient implements the PromptManager interface
var _ PromptManager = &PromptDBClient{}

// Prompt represents a prompt in the system.
type Prompt struct {
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// PromptDBClient is a client for interacting with DynamoDB for prompt-related operations.
type PromptDBClient struct {
	service *dynamodb.Client
}

// NewPromptDBClient creates a new PromptDBClient with AWS configuration.
func NewPromptDBClient() (*PromptDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	svc := dynamodb.NewFromConfig(cfg)
	return &PromptDBClient{
		service: svc,
	}, nil
}

// CreatePrompt creates a new prompt in the DynamoDB table.
func (d *PromptDBClient) CreatePrompt(ctx context.Context, prompt Prompt) error {
	now := time.Now().UTC().Format(time.RFC3339)
	prompt.CreatedAt = now
	prompt.UpdatedAt = now

	av, err := attributevalue.MarshalMap(prompt)
	if err != nil {
		return fmt.Errorf("failed to marshal prompt: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Prompts"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetPrompt retrieves a prompt by ID from the DynamoDB table.
func (d *PromptDBClient) GetPrompt(ctx context.Context, id string) (*Prompt, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Prompts"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
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

	return &prompt, nil
}

// UpdatePrompt updates an existing prompt in the DynamoDB table.
func (d *PromptDBClient) UpdatePrompt(ctx context.Context, prompt Prompt) error {
	prompt.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(prompt)
	if err != nil {
		return fmt.Errorf("failed to marshal prompt: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Prompts"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeletePrompt deletes a prompt by ID from the DynamoDB table.
func (d *PromptDBClient) DeletePrompt(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Prompts"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	_, err := d.service.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from DynamoDB: %v", err)
	}

	return nil
}

// ListPrompts retrieves all prompts from the DynamoDB table.
func (d *PromptDBClient) ListPrompts(ctx context.Context) ([]Prompt, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Prompts"),
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var prompts []Prompt
	err = attributevalue.UnmarshalListOfMaps(result.Items, &prompts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return prompts, nil
}

// ListPromptsByProject retrieves all prompts belonging to a specific project from the DynamoDB table.
func (d *PromptDBClient) ListPromptsByProject(ctx context.Context, projectID string) ([]Prompt, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("Prompts"),
		FilterExpression: aws.String("projectId = :projectId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":projectId": &types.AttributeValueMemberS{Value: projectID},
		},
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var prompts []Prompt
	err = attributevalue.UnmarshalListOfMaps(result.Items, &prompts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return prompts, nil
}
