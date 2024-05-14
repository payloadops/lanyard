package main

import (
	"context"
	"github.com/payloadops/plato/api/dal"
	"github.com/payloadops/plato/api/openapi"
	"github.com/payloadops/plato/api/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewLogger initializes a Zap logger suitable for the given environment.
func NewLogger() (*zap.Logger, error) {
	if os.Getenv("ENVIRONMENT") == "local" {
		return newLocalLogger()
	}
	return newProductionLogger()
}

// newLocalLogger initializes a Zap sugared logger for local development.
func newLocalLogger() (*zap.Logger, error) {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	return config.Build()
}

// newProductionLogger initializes a Zap logger suitable for production.
func newProductionLogger() (*zap.Logger, error) {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	return config.Build()
}

// LoadAWSConfig loads AWS configuration based on the environment
func LoadAWSConfig() (aws.Config, error) {
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
				case "elasticache":
					return aws.Endpoint{URL: os.Getenv("ELASTICACHE_ENDPOINT")}, nil
				case cloudwatch.ServiceID:
					return aws.Endpoint{URL: os.Getenv("CLOUDWATCH_ENDPOINT")}, nil
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

func main() {
	logger, err := NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	awsConfig, err := LoadAWSConfig()
	if err != nil {
		sugar.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create AWS clients
	dynamoClient := dynamodb.NewFromConfig(awsConfig)
	s3Client := s3.NewFromConfig(awsConfig)
	elastiCacheClient := elasticache.NewFromConfig(awsConfig)
	cloudwatchClient := cloudwatch.NewFromConfig(awsConfig)

	// Initialize database clients
	commitDBClient := &dal.CommitDBClient{dynamoDb: dynamoClient, s3: s3Client}
	branchDBClient := &dal.BranchDBClient{service: dynamoClient}
	orgDBClient := &dal.OrgDBClient{service: dynamoClient}
	projectDBClient := &dal.ProjectDBClient{service: dynamoClient}
	promptDBClient := &dal.PromptDBClient{service: dynamoClient}
	teamDBClient := &dal.TeamDBClient{service: dynamoClient}
	userDBClient := &dal.UserDBClient{service: dynamoClient}
	apiKeyDBClient := &dal.APIKeyDBClient{service: dynamoClient}

	// Initialize services with injected dependencies
	APIKeysAPIService := service.NewAPIKeysAPIService(apiKeyDBClient, projectDBClient)
	BranchesAPIService := service.NewBranchesAPIService(branchDBClient, promptDBClient)
	CommitsAPIService := service.NewCommitsAPIService(commitDBClient, branchDBClient)
	HealthCheckAPIService := service.NewHealthCheckAPIService()
	OrganizationsAPIService := service.NewOrganizationsAPIService(orgDBClient)
	ProjectsAPIService := service.NewProjectsAPIService(projectDBClient, orgDBClient)
	PromptsAPIService := service.NewPromptsAPIService(promptDBClient, projectDBClient)
	TeamsAPIService := service.NewTeamsAPIService(teamDBClient, orgDBClient)
	UsersAPIService := service.NewUsersAPIService(userDBClient)

	// Initialize controllers
	APIKeysAPIController := openapi.NewAPIKeysAPIController(APIKeysAPIService)
	BranchesAPIController := openapi.NewBranchesAPIController(BranchesAPIService)
	CommitsAPIController := openapi.NewCommitsAPIController(CommitsAPIService)
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)
	OrganizationsAPIController := openapi.NewOrganizationsAPIController(OrganizationsAPIService)
	ProjectsAPIController := openapi.NewProjectsAPIController(ProjectsAPIService)
	PromptsAPIController := openapi.NewPromptsAPIController(PromptsAPIService)
	TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)
	UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)

	// Initialize router
	router := openapi.NewRouter(
		APIKeysAPIController,
		BranchesAPIController,
		CommitsAPIController,
		HealthCheckAPIController,
		OrganizationsAPIController,
		ProjectsAPIController,
		PromptsAPIController,
		TeamsAPIController,
		UsersAPIController,
	)

	// Initialize server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Listen: %s\n", err)
		}
	}()
	sugar.Infof("Server started on :8080")

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server forced to shutdown: %v", err)
	}

	sugar.Infof("Server exiting")
}
