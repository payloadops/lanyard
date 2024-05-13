package dbdal

import (
	"context"
	"fmt"
	"time"

	awsclient "plato/app/go/client/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

var PROMPT_TABLE_NAME = aws.String("ProjectsPrompts") // Updated to the correct table name

// Prompt represents the structure of a prompt record in the database.
type Prompt struct {
	Id            string `dynamodbav:"promptId" json:"prompt_id"`
	Name          string `dynamodbav:"name" json:"name"`
	ProjectId     string `dynamodbav:"projectId" json:"project_id"`
	PromptS3Path  string `dynamodbav:"promptS3Path" json:"prompt_s3_path"`
	Deleted       bool   `dynamodbav:"deleted" json:"deleted"`
	Version       string `dynamodbav:"version" json:"version"`
	Stub          string `dynamodbav:"stub" json:"stub"`
	DefaultBranch string `dynamodbav:"default_branch" json:"default_branch"`
	ModifiedAt    string `dynamodbav:"modifiedAt" json:"modified_at"`
}

// ListPromptsByProjectId fetches prompts for a given project Id
func ListPromptsByProjectId(ctx context.Context, projectId string) ([]Prompt, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId)
	skPrefix := "PROMPT#"

	params := &dynamodb.QueryInput{
		TableName:              PROMPT_TABLE_NAME,
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":    &types.AttributeValueMemberS{Value: pk},
			":sk":    &types.AttributeValueMemberS{Value: skPrefix},
			":false": &types.AttributeValueMemberBOOL{Value: false},
		},
		FilterExpression: aws.String("deleted = :false"),
	}

	resp, err := awsclient.GetDynamoClient().Query(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error querying prompts: %w", err)
	}

	var prompts []Prompt
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &prompts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling prompts: %w", err)
	}

	return prompts, nil
}

// GetPromptById retrieves a prompt by its Id
func GetPromptById(ctx context.Context, projectId string, promptId string) (*Prompt, error) {
	prompt := &Prompt{Id: promptId}
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := fmt.Sprintf("PROMPT#%s", promptId)

	resp, err := awsclient.GetDynamoClient().GetItem(ctx, &dynamodb.GetItemInput{
		TableName: PROMPT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting prompt: %w", err)
	}

	err = attributevalue.UnmarshalMap(resp.Item, &prompt)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling prompt: %w", err)
	}

	return prompt, nil
}

// AddPrompt adds a new prompt to the database
func AddPrompt(ctx context.Context, name string, stub string, projectId string, promptId string, promptS3Path string, version string, branch string) (*Prompt, error) {
	prompt := &Prompt{
		ProjectId:     projectId,
		Name:          name,
		Id:            promptId,
		DefaultBranch: branch,
		PromptS3Path:  promptS3Path,
		Version:       version,
		Deleted:       false,
		Stub:          stub,
		ModifiedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	pk := fmt.Sprintf("PROJECT#%s", projectId)
	sk := fmt.Sprintf("PROMPT#%s", promptId)

	item, err := attributevalue.MarshalMap(prompt)
	if err != nil {
		return nil, fmt.Errorf("error marshaling prompt: %w", err)
	}

	item["PK"] = &types.AttributeValueMemberS{Value: pk}
	item["SK"] = &types.AttributeValueMemberS{Value: sk}

	_, err = awsclient.GetDynamoClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: PROMPT_TABLE_NAME,
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("error adding prompt: %w", err)
	}

	return prompt, nil
}

// UpdatePromptDeletedStatus updates the 'deleted' status of a prompt
func UpdatePromptDeletedStatus(ctx context.Context, projectId string, promptId string, deleted bool) (string, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := fmt.Sprintf("PROMPT#%s", promptId)

	modifiedAt := time.Now().UTC().Format(time.RFC3339)

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROMPT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("set deleted = :deleted, modifiedAt = :modifiedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":deleted":    &types.AttributeValueMemberBOOL{Value: deleted},
			":modifiedAt": &types.AttributeValueMemberS{Value: modifiedAt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error updating prompt deleted status: %w", err)
	}

	return modifiedAt, nil
}

func UpdatePromptActiveVersion(ctx context.Context, projectId string, promptId string, stub string, version string) (string, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := fmt.Sprintf("PROMPT#%s", promptId)

	modifiedAt := time.Now().UTC().Format(time.RFC3339)

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROMPT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET stub = :stub, version = :version, modifiedAt = :modifiedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":stub":       &types.AttributeValueMemberS{Value: stub},
			":version":    &types.AttributeValueMemberS{Value: version},
			":modifiedAt": &types.AttributeValueMemberS{Value: modifiedAt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error updating prompt deleted status: %w", err)
	}

	return modifiedAt, nil
}

func UpdatePrompt(ctx context.Context, name string, projectId string, promptId string, stub string, version string) (string, error) {
	pk := fmt.Sprintf("PROJECT#%s", projectId) // This assumes that the PK is based on the project Id
	sk := fmt.Sprintf("PROMPT#%s", promptId)

	modifiedAt := time.Now().UTC().Format(time.RFC3339)

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROMPT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:         aws.String("SET stub = :stub, version = :version, modifiedAt = :modifiedAt, #n = :name"),
		ExpressionAttributeNames: map[string]string{"#n": "name"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":       &types.AttributeValueMemberS{Value: name},
			":stub":       &types.AttributeValueMemberS{Value: stub},
			":version":    &types.AttributeValueMemberS{Value: version},
			":modifiedAt": &types.AttributeValueMemberS{Value: modifiedAt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error updating prompt: %w", err)
	}

	return modifiedAt, nil
}
