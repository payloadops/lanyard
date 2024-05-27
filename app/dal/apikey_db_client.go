package dal

/*
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

//go:generate mockgen -package=mocks -destination=mocks/mock_apikey_db_client.go "github.com/payloadops/plato/api/dal" APIKeyManager

// APIKeyManager defines the operations available for managing API keys.
type APIKeyManager interface {
	CreateAPIKey(ctx context.Context, apiKey APIKey) error
	GetAPIKey(ctx context.Context, id string) (*APIKey, error)
	UpdateAPIKey(ctx context.Context, apiKey APIKey) error
	DeleteAPIKey(ctx context.Context, id string) error
	ListAPIKeys(ctx context.Context, projectId string) ([]APIKey, error)
}

// Ensure APIKeyDBClient implements the APIKeyManager interface
var _ APIKeyManager = &APIKeyDBClient{}

// APIKey represents an API key associated with a project.
type APIKey struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"projectId"`
	Key       string   `json:"key"`
	Scopes    []string `json:"scopes"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

// APIKeyDBClient is a client for interacting with DynamoDB for API key-related operations.
type APIKeyDBClient struct {
	service *dynamodb.Client
}

// NewAPIKeyDBClient creates a new APIKeyDBClient with the given AWS configuration.
func NewAPIKeyDBClient() (*APIKeyDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	svc := dynamodb.NewFromConfig(cfg)
	return &APIKeyDBClient{
		service: svc,
	}, nil
}

// CreateAPIKey creates a new API key in the DynamoDB table.
func (d *APIKeyDBClient) CreateAPIKey(ctx context.Context, apiKey APIKey) error {
	now := time.Now().UTC().Format(time.RFC3339)
	apiKey.CreatedAt = now
	apiKey.UpdatedAt = now

	av, err := attributevalue.MarshalMap(apiKey)
	if err != nil {
		return fmt.Errorf("failed to marshal API key: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("APIKeys"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetAPIKey retrieves an API key by ID from the DynamoDB table.
func (d *APIKeyDBClient) GetAPIKey(ctx context.Context, id string) (*APIKey, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("APIKeys"),
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

	var apiKey APIKey
	err = attributevalue.UnmarshalMap(result.Item, &apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from DynamoDB: %v", err)
	}

	return &apiKey, nil
}

// UpdateAPIKey updates an existing API key in the DynamoDB table.
func (d *APIKeyDBClient) UpdateAPIKey(ctx context.Context, apiKey APIKey) error {
	apiKey.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(apiKey)
	if err != nil {
		return fmt.Errorf("failed to marshal API key: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("APIKeys"),
		Item:      av,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteAPIKey deletes an API key by ID from the DynamoDB table.
func (d *APIKeyDBClient) DeleteAPIKey(ctx context.Context, id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("APIKeys"),
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

// ListAPIKeys retrieves all API keys for a specific project from the DynamoDB table.
func (d *APIKeyDBClient) ListAPIKeys(ctx context.Context, projectId string) ([]APIKey, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String("APIKeys"),
		FilterExpression: aws.String("projectId = :projectId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":projectId": &types.AttributeValueMemberS{Value: projectId},
		},
	}

	result, err := d.service.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan items in DynamoDB: %v", err)
	}

	var apiKeys []APIKey
	err = attributevalue.UnmarshalListOfMaps(result.Items, &apiKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	return apiKeys, nil
}
*/
