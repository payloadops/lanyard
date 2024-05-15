package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setEnv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
}

func unsetEnv(key string) {
	err := os.Unsetenv(key)
	if err != nil {
		panic(err)
	}
}

func TestLoadConfig(t *testing.T) {
	// Set up the environment variables
	setEnv("AWS_DEFAULT_REGION", "us-west-2")
	setEnv("AWS_ACCESS_KEY_ID", "test-access-key-id")
	setEnv("AWS_SECRET_ACCESS_KEY", "test-secret-access-key")
	setEnv("ENVIRONMENT", "local")
	setEnv("DYNAMODB_ENDPOINT", "http://localhost:8000")
	setEnv("S3_ENDPOINT", "http://localhost:4566")
	setEnv("ELASTICACHE_ENDPOINT", "http://localhost:6379")
	setEnv("CLOUDWATCH_ENDPOINT", "http://localhost:4582")
	setEnv("BIND_ADDRESS", ":8080")
	setEnv("JWT_SECRET", "test-jwt-secret")

	defer unsetEnv("AWS_DEFAULT_REGION")
	defer unsetEnv("AWS_ACCESS_KEY_ID")
	defer unsetEnv("AWS_SECRET_ACCESS_KEY")
	defer unsetEnv("ENVIRONMENT")
	defer unsetEnv("DYNAMODB_ENDPOINT")
	defer unsetEnv("S3_ENDPOINT")
	defer unsetEnv("ELASTICACHE_ENDPOINT")
	defer unsetEnv("CLOUDWATCH_ENDPOINT")
	defer unsetEnv("BIND_ADDRESS")
	defer unsetEnv("JWT_SECRET")

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "us-west-2", cfg.AWS.Region)
	assert.Equal(t, "test-access-key-id", cfg.AWS.AccessKeyID)
	assert.Equal(t, "test-secret-access-key", cfg.AWS.SecretAccessKey)
	assert.Equal(t, "local", cfg.AWS.Environment)
	assert.Equal(t, "http://localhost:8000", cfg.AWS.DynamoDBEndpoint)
	assert.Equal(t, "http://localhost:4566", cfg.AWS.S3Endpoint)
	assert.Equal(t, "http://localhost:6379", cfg.AWS.ElasticacheEndpoint)
	assert.Equal(t, "http://localhost:4582", cfg.AWS.CloudWatchEndpoint)
	assert.Equal(t, ":8080", cfg.BindAddress)
	assert.Equal(t, "test-jwt-secret", cfg.JWTSecret)
}

func TestLoadConfigMissingEnvVars(t *testing.T) {
	// Unset the environment variables
	unsetEnv("AWS_DEFAULT_REGION")
	unsetEnv("AWS_ACCESS_KEY_ID")
	unsetEnv("AWS_SECRET_ACCESS_KEY")
	unsetEnv("ENVIRONMENT")
	unsetEnv("DYNAMODB_ENDPOINT")
	unsetEnv("S3_ENDPOINT")
	unsetEnv("ELASTICACHE_ENDPOINT")
	unsetEnv("CLOUDWATCH_ENDPOINT")
	unsetEnv("BIND_ADDRESS")
	unsetEnv("JWT_SECRET")

	cfg, err := LoadConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadConfigDefaultValues(t *testing.T) {
	// Set up the environment variables
	setEnv("AWS_DEFAULT_REGION", "us-west-2")
	setEnv("AWS_ACCESS_KEY_ID", "test-access-key-id")
	setEnv("AWS_SECRET_ACCESS_KEY", "test-secret-access-key")
	setEnv("JWT_SECRET", "foo")

	defer unsetEnv("AWS_DEFAULT_REGION")
	defer unsetEnv("AWS_ACCESS_KEY_ID")
	defer unsetEnv("AWS_SECRET_ACCESS_KEY")

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "us-west-2", cfg.AWS.Region)
	assert.Equal(t, "test-access-key-id", cfg.AWS.AccessKeyID)
	assert.Equal(t, "test-secret-access-key", cfg.AWS.SecretAccessKey)
	assert.Equal(t, ":8080", cfg.BindAddress) // default value
}
