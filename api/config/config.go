package config

import (
	"github.com/kelseyhightower/envconfig"
)

// AWSConfig holds AWS-specific configuration values.
type AWSConfig struct {
	Region              string `envconfig:"AWS_DEFAULT_REGION"`
	AccessKeyID         string `envconfig:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey     string `envconfig:"AWS_SECRET_ACCESS_KEY"`
	Environment         string `envconfig:"ENVIRONMENT"`
	DynamoDBEndpoint    string `envconfig:"DYNAMODB_ENDPOINT"`
	S3Endpoint          string `envconfig:"S3_ENDPOINT"`
	ElasticacheEndpoint string `envconfig:"ELASTICACHE_ENDPOINT"`
	CloudWatchEndpoint  string `envconfig:"CLOUDWATCH_ENDPOINT"`
}

// Config holds the entire configuration for the application.
type Config struct {
	Environment string `envconfig:"ENVIRONMENT"`
	BindAddress string `envconfig:"BIND_ADDRESS" default:":8080"`
	JWTSecret   string `envconfig:"JWT_SECRET" required:"true"`
	AWS         AWSConfig
}

// LoadConfig loads the configuration from environment variables into the Config struct.
// It returns an error if any required environment variable is missing or has an invalid value.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
