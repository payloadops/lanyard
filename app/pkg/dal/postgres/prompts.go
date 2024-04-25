package dbdal

import (
	"context"
	"fmt"

	awsclient "plato/app/pkg/client/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

var PROMPT_TABLE_NAME = aws.String("ProjectsPrompts") // Updated to the correct table name

// Prompt represents the structure of a prompt record in the database.
type Prompt struct {
	Id           string `dynamodbav:"promptId" json:"prompt_id"`
	ProjectId    string `dynamodbav:"projectId" json:"project_id"`
	PromptS3Path string `dynamodbav:"promptS3Path" json:"prompt_s3_path"`
	Deleted      bool   `dynamodbav:"deleted" json:"deleted"`
	Version      string `dynamodbav:"version" json:"version"`
	Stub         string `dynamodbav:"stub" json:"stub"`
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
func AddPrompt(ctx context.Context, stub string, projectId string, promptId string, promptS3Path string, version string) (*Prompt, error) {
	prompt := &Prompt{
		ProjectId:    projectId,
		Id:           promptId,
		PromptS3Path: promptS3Path,
		Version:      version,
		Deleted:      false,
		Stub:         stub,
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
func UpdatePromptDeletedStatus(ctx context.Context, id string, deleted bool) error {
	pk := fmt.Sprintf("PROMPT#%s", id)
	sk := pk // Assuming SK is the same as PK

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROMPT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("set deleted = :deleted"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":deleted": &types.AttributeValueMemberBOOL{Value: deleted},
		},
	})
	if err != nil {
		return fmt.Errorf("error updating prompt deleted status: %w", err)
	}

	return nil
}

func UpdatePrompt(ctx context.Context, id string, stub string, version string) error {
	pk := fmt.Sprintf("PROMPT#%s", id)
	sk := pk // Assuming SK is the same as PK

	_, err := awsclient.GetDynamoClient().UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: PROMPT_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET stub = :stub, version = :version"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":stub":    &types.AttributeValueMemberS{Value: stub},
			":version": &types.AttributeValueMemberS{Value: version},
		},
	})
	if err != nil {
		return fmt.Errorf("error updating prompt deleted status: %w", err)
	}

	return nil
}
