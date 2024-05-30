package client

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/payloadops/plato/app/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadAWSConfig_LocalEnvironment(t *testing.T) {
	cfg := &config.Config{
		Environment: config.Local,
		AWS: config.AWSConfig{
			Region:           "us-west-2",
			AccessKeyID:      "dummyAccessKey",
			SecretAccessKey:  "dummySecretKey",
			DynamoDBEndpoint: "http://localhost:4566",
			S3Endpoint:       "http://localhost:4572",
		},
	}

	awsCfg, err := LoadAWSConfig(cfg)
	assert.NoError(t, err)

	assert.Equal(t, "us-west-2", awsCfg.Region)

	endpointResolver := awsCfg.EndpointResolverWithOptions.(aws.EndpointResolverWithOptions)
	endpoint, err := endpointResolver.ResolveEndpoint(dynamodb.ServiceID, "us-west-2")
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:4566", endpoint.URL)

	endpoint, err = endpointResolver.ResolveEndpoint(s3.ServiceID, "us-west-2")
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:4572", endpoint.URL)
}

func TestLoadAWSConfig_ProductionEnvironment(t *testing.T) {
	cfg := &config.Config{
		Environment: config.Production,
		AWS: config.AWSConfig{
			Region:          "us-west-2",
			AccessKeyID:     "dummyAccessKey",
			SecretAccessKey: "dummySecretKey",
		},
	}

	awsCfg, err := LoadAWSConfig(cfg)
	assert.NoError(t, err)

	assert.Equal(t, "us-west-2", awsCfg.Region)

	endpointResolver := awsCfg.EndpointResolverWithOptions
	assert.Nil(t, endpointResolver)
}
