package dal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//go:generate mockgen -package=mocks -destination=mocks/mock_dynamo_db_client.go "github.com/payloadops/plato/app/dal" DynamoDBAPI

// Ensure DynamoDBClient implements the DynamoDBAPI interface
var _ DynamoDBAPI = &dynamodb.Client{}

// DynamoDBAPI defines the interface for the DynamoDB client
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}
