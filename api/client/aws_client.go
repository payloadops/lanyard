package client

import (
	"context"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/payloadops/plato/api/config"
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
	case cloudwatch.ServiceID:
		return aws.Endpoint{URL: r.cfg.AWS.CloudWatchEndpoint}, nil
	default:
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}

// LoadAWSConfig loads AWS configuration based on the environment.
func LoadAWSConfig(cfg *config.Config) (aws.Config, error) {
	options := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.AWS.Region),
	}

	if cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
		options = append(options, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, ""),
		))
	}

	if cfg.Environment == config.Local {
		options = append(options, awsconfig.WithEndpointResolverWithOptions(localEndpointResolver{
			cfg: cfg,
		}))
	}

	return awsconfig.LoadDefaultConfig(context.TODO(), options...)
}
