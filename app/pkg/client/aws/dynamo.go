package awsclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dynamoOnce    sync.Once
	dynamoClient  *dynamodb.Client
	dynamoInitErr error
)

func InitDynamoClient() (*dynamodb.Client, error) {
	dynamoOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"), config.WithClientLogMode(aws.LogRetries|aws.LogRequestWithBody|aws.LogResponseWithBody))
		if err != nil {
			dynamoInitErr = fmt.Errorf("failed to load AWS config: %w", err)
			return
		}
		dynamoClient = dynamodb.NewFromConfig(cfg)
	})
	return dynamoClient, dynamoInitErr
}

func GetDynamoClient() *dynamodb.Client {
	if dynamoClient == nil {
		panic("Client not initialized")
	}
	return dynamoClient
}
