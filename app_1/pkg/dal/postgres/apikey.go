package dbdal

import (
	"context"
	"fmt"

	"plato/app/pkg/util"
	awsclient "plato/app_1/pkg/client/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var PROJECT_KEYS_TABLE = "Keys"

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
	return apiKeys, nil
}

// GetApiKey retrieves an Api key by its string value from the database.
func GetApiKey(ctx context.Context, apiKeyString string) (*ApiKeyItem, error) {
	apiKey := &ApiKeyItem{}
	pk := fmt.Sprintf("KEY#%s", apiKeyString)

	resp, err := awsclient.GetDynamoClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &PROJECT_KEYS_TABLE,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
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
		ApiKey:    util.GenIDString(),
		ProjectId: projectId,
		OrgId:     orgId,
		RateLimit: 1000,
		Active:    true,
		Scopes:    scopes,
	}
	pk := fmt.Sprintf("KEY#%s", apiKey.ApiKey)
	sk := fmt.Sprintf("PROJECT#%s", projectId)

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
func UpdateApiKey(ctx context.Context, projectId string, apiKeyId string, newDesc string, newScopes []string) error {
	return nil
}

// DeactivateApiKey deactivates a specific Api key.
func DeactivateApiKey(ctx context.Context, projectId string, apiKeyId string) error {
	return nil
}
