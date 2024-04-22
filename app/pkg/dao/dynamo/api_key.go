package dynamodao

import (
	aws_client "plato/app/pkg/client/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type APIKeyItem struct {
	ApiKey    string   `json:"api_key"`
	ProjectId string   `json:"project_id"`
	RateLimit int      `json:"rate_limit"`
	Active    bool     `json:"active"`
	Scopes    []string `json:"scopes"`
}

var API_KEYS_TABLE_NAME = "ApiKeys"

func GetApiKey(apiKey string) (*APIKeyItem, error) {
	result, err := aws_client.GetDynamoDBClient().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(API_KEYS_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"ApiKey": {
				S: aws.String(apiKey),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var apiKeyData APIKeyItem
	if err := dynamodbattribute.UnmarshalMap(result.Item, &apiKeyData); err != nil {
		return nil, err
	}
	return &apiKeyData, nil
}

func CreateApiKey(data *APIKeyItem) error {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(API_KEYS_TABLE_NAME),
	}

	_, err = aws_client.GetDynamoDBClient().PutItem(input)
	return err
}

func DeleteApiKey(apiKey string) error {
	_, err := aws_client.GetDynamoDBClient().DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(API_KEYS_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"ApiKey": {
				S: aws.String(apiKey),
			},
		},
	})
	return err
}

func UpdateApiKey(data *APIKeyItem) error {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	_, err = aws_client.GetDynamoDBClient().PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(API_KEYS_TABLE_NAME),
	})
	return err
}
