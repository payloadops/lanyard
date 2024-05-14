package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// InitializeDynamoDBClient initializes and returns a DynamoDB client
func InitializeDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	// Create the DynamoDB service client
	svc := dynamodb.NewFromConfig(cfg)
	return svc, nil
}
