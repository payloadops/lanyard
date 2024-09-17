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
	setEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
	setEnv("OTEL_EXPORTER_OTLP_CA_CERT", "test-ca-cert")
	setEnv("BIND_ADDRESS", ":8080")
	setEnv("JWT_SECRET", "test-jwt-secret")
	setEnv("PROMPT_BUCKET", "test-prompt-bucket")

	defer unsetEnv("AWS_DEFAULT_REGION")
	defer unsetEnv("AWS_ACCESS_KEY_ID")
	defer unsetEnv("AWS_SECRET_ACCESS_KEY")
	defer unsetEnv("ENVIRONMENT")
	defer unsetEnv("DYNAMODB_ENDPOINT")
	defer unsetEnv("S3_ENDPOINT")
	defer unsetEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer unsetEnv("OTEL_EXPORTER_OTLP_CA_CERT")
	defer unsetEnv("BIND_ADDRESS")
	defer unsetEnv("JWT_SECRET")
	defer unsetEnv("PROMPT_BUCKET")

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "us-west-2", cfg.AWS.Region)
	assert.Equal(t, "test-access-key-id", cfg.AWS.AccessKeyID)
	assert.Equal(t, "test-secret-access-key", cfg.AWS.SecretAccessKey)
	assert.Equal(t, "local", string(cfg.AWS.Environment))
	assert.Equal(t, "http://localhost:8000", cfg.AWS.DynamoDBEndpoint)
	assert.Equal(t, "http://localhost:4566", cfg.AWS.S3Endpoint)
	assert.Equal(t, "local", string(cfg.Environment))
	assert.Equal(t, ":8080", cfg.BindAddress)
	assert.Equal(t, "test-jwt-secret", cfg.JWTSecret)
	assert.Equal(t, "http://localhost:4317", cfg.OpenTelemetry.ProviderEndpoint)
	assert.Equal(t, "test-ca-cert", cfg.OpenTelemetry.CACert)
}

func TestLoadConfigMissingEnvVars(t *testing.T) {
	// Unset the environment variables
	unsetEnv("AWS_DEFAULT_REGION")
	unsetEnv("AWS_ACCESS_KEY_ID")
	unsetEnv("AWS_SECRET_ACCESS_KEY")
	unsetEnv("ENVIRONMENT")
	unsetEnv("DYNAMODB_ENDPOINT")
	unsetEnv("S3_ENDPOINT")
	unsetEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	unsetEnv("OTEL_EXPORTER_OTLP_CA_CERT")
	unsetEnv("BIND_ADDRESS")
	unsetEnv("JWT_SECRET")
	unsetEnv("PROMPT_BUCKET")

	cfg, err := LoadConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadConfigDefaultValues(t *testing.T) {
	// Set up the environment variables
	setEnv("AWS_DEFAULT_REGION", "us-west-2")
	setEnv("AWS_ACCESS_KEY_ID", "test-access-key-id")
	setEnv("AWS_SECRET_ACCESS_KEY", "test-secret-access-key")
	setEnv("JWT_SECRET", "test-jwt-secret")
	setEnv("PROMPT_BUCKET", "test-prompt-bucket")

	defer unsetEnv("AWS_DEFAULT_REGION")
	defer unsetEnv("AWS_ACCESS_KEY_ID")
	defer unsetEnv("AWS_SECRET_ACCESS_KEY")

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "us-west-2", cfg.AWS.Region)
	assert.Equal(t, "test-access-key-id", cfg.AWS.AccessKeyID)
	assert.Equal(t, "test-secret-access-key", cfg.AWS.SecretAccessKey)
	assert.Equal(t, ":8080", cfg.BindAddress)               // default value
	assert.Equal(t, "", cfg.OpenTelemetry.ProviderEndpoint) // default value when not set
	assert.Equal(t, "", cfg.OpenTelemetry.CACert)
}
