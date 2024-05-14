package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

// InitializeAWSConfig loads AWS configuration based on the environment
func InitializeAWSConfig() (aws.Config, error) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	endpointResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if os.Getenv("ENVIRONMENT") == "local" {
				switch service {
				case dynamodb.ServiceID:
					return aws.Endpoint{URL: os.Getenv("DYNAMODB_ENDPOINT")}, nil
				case s3.ServiceID:
					return aws.Endpoint{URL: os.Getenv("S3_ENDPOINT")}, nil
				default:
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				}
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

	options := []func(*config.LoadOptions) error{
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}

	if os.Getenv("ENVIRONMENT") == "local" {
		options = append(options, config.WithEndpointResolverWithOptions(endpointResolver))
	}

	return config.LoadDefaultConfig(context.TODO(), options...)
}
