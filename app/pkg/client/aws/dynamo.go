package aws_client

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetDynamoDBClient() *dynamodb.DynamoDB {
	var once sync.Once
	var singletonClient *dynamodb.DynamoDB
	once.Do(func() {
		singletonClient = dynamodb.New(GetAWSSession())
	})
	return singletonClient
}
