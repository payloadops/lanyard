package client

import (
	"context"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/payloadops/plato/app/config"
)

// localEndpointResolver resolves AWS service endpoints for local development.
type localEndpointResolver struct {
	cfg *config.Config
}

// ResolveEndpoint resolves the endpoint for a given AWS service.
func (r localEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	if r.cfg.Environment != config.Local {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	switch service {
	case dynamodb.ServiceID:
		return aws.Endpoint{URL: r.cfg.AWS.DynamoDBEndpoint}, nil
	case s3.ServiceID:
		return aws.Endpoint{URL: r.cfg.AWS.S3Endpoint}, nil
	default:
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}

// LoadAWSConfig loads AWS configuration based on the environment.
func LoadAWSConfig(cfg *config.Config) (aws.Config, error) {
	options := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.AWS.Region),
	}

	// Add credentials if provided (used for local development or specific scenarios)
	if cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
		options = append(options, awsconfig.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, "")),
		))
	}

	// Add local endpoint resolver for local development environment
	if cfg.Environment == config.Local {
		options = append(options, awsconfig.WithEndpointResolverWithOptions(localEndpointResolver{
			cfg: cfg,
		}))
	}

	// Load default AWS config which includes support for ECS task roles
	return awsconfig.LoadDefaultConfig(context.TODO(), options...)
}
