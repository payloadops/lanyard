package dbdal

import (
	"context"
	"fmt"

	awsclient "plato/app/pkg/client/aws"
	dbClient "plato/app/pkg/client/db"
	"plato/app/pkg/util"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var PROJECT_KEYS_TABLE = "ProjectKeys"

type ApiKeyItem struct {
	ApiKey    string   `dynamodbav:"apiKey,pk" json:"api_key"`
	ProjectId string   `dynamodbav:"projectId" json:"project_id"`
	OrgId     string   `dynamodbav:"orgId" json:"org_id"`
	RateLimit int      `dynamodbav:"rateLimit" json:"rate_limit"`
	Active    bool     `dynamodbav:"active" json:"active"`
	Scopes    []string `dynamodbav:"scopes" json:"scopes"`
}

// Lists active Api keys by project id
func ListApiKeysByProjectId(ctx context.Context, projectId string) (*[]ApiKeyItem, error) {
	apiKeys := &[]ApiKeyItem{}
	err := dbClient.GetClient().NewSelect().Model(apiKeys).Where("project_id = ?", projectId).Where("active = TRUE").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying Api keys: %w", err)
	}
	return apiKeys, nil
}

// GetApiKey retrieves an Api key by its string value from the database.
func GetApiKey(ctx context.Context, apiKeyString string) (*ApiKeyItem, error) {
	apiKey := &ApiKeyItem{}
	pk := fmt.Sprintf("KEY#%s", apiKeyString)
	sk := pk

	resp, err := awsclient.GetDynamoClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &PROJECT_KEYS_TABLE,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &apiKey)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return apiKey, nil
}

// CreateApiKey creates a new Api key in the database.
func CreateApiKey(
	ctx context.Context,
	orgId string,
	projectId string,
	desc string,
	scopes []string,
) (*ApiKeyItem, error) {
	apiKey := &ApiKeyItem{
		ApiKey:    util.GenUUIDString(),
		ProjectId: projectId,
		OrgId:     orgId,
		RateLimit: 1000,
		Active:    true,
		Scopes:    scopes,
	}
	pk := fmt.Sprintf("KEY#%s", apiKey.ApiKey)
	sk := pk

	item, err := attributevalue.MarshalMap(apiKey)
	if err != nil {
		return nil, fmt.Errorf("error marshaling prompt: %w", err)
	}

	item["PK"] = &types.AttributeValueMemberS{Value: pk}
	item["SK"] = &types.AttributeValueMemberS{Value: sk}

	_, err = awsclient.GetDynamoClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &PROJECT_KEYS_TABLE,
		Item:      item,
	})

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return apiKey, nil
}

// UpdateApiKey updates an existing Api key's description and scopes.
func UpdateApiKey(ctx context.Context, apiKeyId, newDesc string, newScopes []string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&ApiKeyItem{}).Set("description = ?", newDesc).Set("scopes = ?", newScopes).Where("api_key = ?", apiKeyId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error updating Api key: %w", err)
	}
	return nil
}

// DeactivateApiKey deactivates a specific Api key.
func DeactivateApiKey(ctx context.Context, apiKeyId string) error {
	_, err := dbClient.GetClient().NewUpdate().Model(&ApiKeyItem{}).Set("active = false").Where("api_key = ?", apiKeyId).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deactivating Api key: %w", err)
	}
	return nil
}
