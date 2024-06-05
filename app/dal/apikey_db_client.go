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

// SecretLength represents the length of the secret to generate for API keys.
const SecretLength = 32

//go:generate mockgen -package=mocks -destination=mocks/mock_apikey_db_client.go "github.com/payloadops/plato/app/dal" APIKeyManager

// APIKeyManager defines the operations available for managing API keys.
type APIKeyManager interface {
	CreateAPIKey(ctx context.Context, apiKey *APIKey) error
	GetAPIKey(ctx context.Context, orgId, projectID, apiKeyID string) (*APIKey, error)
	GetAPIKeyByID(ctx context.Context, apiKeyID string) (*APIKey, error)
	UpdateAPIKey(ctx context.Context, apiKey *APIKey) error
	DeleteAPIKey(ctx context.Context, orgID, projectID, apiKeyID string) error
	ListAPIKeysByProject(ctx context.Context, orgID, projectID string) ([]APIKey, error)
}

// Ensure APIKeyDBClient implements the APIKeyManager interface
var _ APIKeyManager = &APIKeyDBClient{}

// APIKey represents an API key associated with a project.
type APIKey struct {
	OrgID     string   `json:"orgId"`
	ProjectID string   `json:"projectId"`
	APIKeyID  string   `json:"apiKeyId"`
	Secret    string   `json:"secret"`
	Scopes    []string `json:"scopes"`
	Deleted   bool
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// APIKeyDBClient is a client for interacting with DynamoDB for API key-related operations.
type APIKeyDBClient struct {
	service DynamoDBAPI
}

// NewAPIKeyDBClient creates a new APIKeyDBClient.
func NewAPIKeyDBClient(service DynamoDBAPI) *APIKeyDBClient {
	return &APIKeyDBClient{
		service: service,
	}
}

// createAPIKeyGSICompositeKeys generates the partition key (pk) and sort key (sk) for an API key.
func createAPIKeyGSICompositeKeys(orgID, projectID, apiKeyID string) (string, string) {
	return "Org#" + orgID + "Project#" + projectID, "APIKey#" + apiKeyID
}

// createAPIKeyCompositeKeys generates the partition key (pk) for an API key.
func createAPIKeyCompositeKeys(apiKeyID string) string {
	return "APIKey#" + apiKeyID
}

// CreateAPIKey creates a new API key in the DynamoDB table.
func (d *APIKeyDBClient) CreateAPIKey(ctx context.Context, apiKey *APIKey) error {
	ksuid, err := utils.GenerateKSUID()
	if err != nil {
		return fmt.Errorf("failed to create ksuid: %v", err)
	}

	apiKey.APIKeyID = ksuid
	pk := createAPIKeyCompositeKeys(apiKey.APIKeyID)
	gsiPK, gsiSK := createAPIKeyGSICompositeKeys(apiKey.OrgID, apiKey.ProjectID, apiKey.APIKeyID)

	now := time.Now().UTC().Format(time.RFC3339)
	apiKey.CreatedAt = now
	apiKey.UpdatedAt = now

	av, err := attributevalue.MarshalMap(apiKey)
	if err != nil {
		return fmt.Errorf("failed to marshal API key: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk":     &types.AttributeValueMemberS{Value: pk},
		"GSI1PK": &types.AttributeValueMemberS{Value: gsiPK},
		"GSISK":  &types.AttributeValueMemberS{Value: gsiSK},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("APIKeys"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %v", err)
	}

	return nil
}

// GetAPIKey retrieves an API key by org ID, project ID and API key ID from the DynamoDB table.
func (d *APIKeyDBClient) GetAPIKey(ctx context.Context, orgID, projectID, apiKeyID string) (*APIKey, error) {
	pk := createAPIKeyCompositeKeys(apiKeyID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("APIKeys"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
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

	if apiKey.Deleted {
		return nil, nil
	}

	return &apiKey, nil
}

func (d *APIKeyDBClient) GetAPIKeyByID(ctx context.Context, apiKeyID string) (*APIKey, error) {
	pk := createAPIKeyCompositeKeys(apiKeyID)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("APIKeys"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
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

	if apiKey.Deleted {
		return nil, nil
	}

	return &apiKey, nil
}

// UpdateAPIKey updates an existing API key in the DynamoDB table.
func (d *APIKeyDBClient) UpdateAPIKey(ctx context.Context, apiKey *APIKey) error {
	pk := createAPIKeyCompositeKeys(apiKey.APIKeyID)
	apiKey.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	av, err := attributevalue.MarshalMap(apiKey)
	if err != nil {
		return fmt.Errorf("failed to marshal API key: %v", err)
	}

	item := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: pk},
	}
	for k, v := range av {
		item[k] = v
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("APIKeys"),
		Item:      item,
	}

	_, err = d.service.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in DynamoDB: %v", err)
	}

	return nil
}

// DeleteAPIKey marks an API key as deleted by org ID, project ID, and API key ID in the DynamoDB table.
func (d *APIKeyDBClient) DeleteAPIKey(ctx context.Context, orgID, projectID, apiKeyID string) error {
	pk := createAPIKeyCompositeKeys(apiKeyID)
	update := map[string]types.AttributeValueUpdate{
		"Deleted": {
			Value:  &types.AttributeValueMemberBOOL{Value: true},
			Action: types.AttributeActionPut,
		},
		"UpdatedAt": {
			Value:  &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
			Action: types.AttributeActionPut,
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("APIKeys"),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
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

// ListAPIKeysByProject retrieves all API keys for a specific project from the DynamoDB table.
func (d *APIKeyDBClient) ListAPIKeysByProject(ctx context.Context, orgID, projectID string) ([]APIKey, error) {
	gsiPK, _ := createAPIKeyGSICompositeKeys(orgID, projectID, "")
	input := &dynamodb.QueryInput{
		TableName:              aws.String("APIKeys"),
		KeyConditionExpression: aws.String("GSI1PK = :gsi1PK"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1PK": &types.AttributeValueMemberS{
				Value: gsiPK,
			},
		},
	}

	result, err := d.service.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items in DynamoDB: %v", err)
	}

	var apiKeys []APIKey
	err = attributevalue.UnmarshalListOfMaps(result.Items, &apiKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items from DynamoDB: %v", err)
	}

	results := []APIKey{}
	for _, apiKey := range apiKeys {
		if apiKey.Deleted {
			continue
		}
		results = append(results, apiKey)
	}

	return results, nil
}
