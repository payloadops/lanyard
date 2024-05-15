package main

import (
	"context"
	"errors"
	"github.com/payloadops/plato/api/config"
	"github.com/payloadops/plato/api/logging"
	"github.com/payloadops/plato/api/openapi"
	"github.com/payloadops/plato/api/service"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
// LoadAWSConfig loads AWS configuration based on the environment
func LoadAWSConfig() (aws.Config, error) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")

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
*/

func main() {
	// Load config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize service logging
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}

	// Set global logger to use this implementation
	// zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Load AWS/localstack config values
	/*
		awsConfig, err := LoadAWSConfig()
		if err != nil {
			sugar.Fatalf("Unable to load SDK config, %v", err)
		}
	*/

	/*
		// Create AWS clients
		dynamoClient := dynamodb.NewFromConfig(awsConfig)
		s3Client := s3.NewFromConfig(awsConfig)
		elastiCacheClient := elasticache.NewFromConfig(awsConfig)
		cloudwatchClient := cloudwatch.NewFromConfig(awsConfig)

		// Initialize Redis client for ElastiCache
		redisClient := redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ENDPOINT"),
		})

		// Create cache instance
		cache := cache.NewRedisCache(redisClient)
	*/

	// Initialize database clients
	/*
		commitDBClient := dal.NewCommitDBClient(dynamoClient, s3Client, cache)
		branchDBClient := &dal.BranchDBClient{service: dynamoClient}
		orgDBClient := &dal.OrgDBClient{service: dynamoClient}
		projectDBClient := &dal.ProjectDBClient{service: dynamoClient}
		promptDBClient := &dal.PromptDBClient{service: dynamoClient}
		teamDBClient := &dal.TeamDBClient{service: dynamoClient}
		userDBClient := &dal.UserDBClient{service: dynamoClient}
		apiKeyDBClient := &dal.APIKeyDBClient{service: dynamoClient}
	*/

	// Initialize the healtcheck service
	HealthCheckAPIService := service.NewHealthCheckAPIService()

	// Initialize services with injected dependencies
	/*
		APIKeysAPIService := service.NewAPIKeysAPIService(apiKeyDBClient, projectDBClient)
		BranchesAPIService := service.NewBranchesAPIService(branchDBClient, promptDBClient)
		CommitsAPIService := service.NewCommitsAPIService(commitDBClient, branchDBClient)
		OrganizationsAPIService := service.NewOrganizationsAPIService(orgDBClient)
		ProjectsAPIService := service.NewProjectsAPIService(projectDBClient, orgDBClient)
		PromptsAPIService := service.NewPromptsAPIService(promptDBClient, projectDBClient)
		TeamsAPIService := service.NewTeamsAPIService(teamDBClient, orgDBClient)
		UsersAPIService := service.NewUsersAPIService(userDBClient)
	*/

	// Initialize controllers
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)
	/*
		APIKeysAPIController := openapi.NewAPIKeysAPIController(APIKeysAPIService)
		BranchesAPIController := openapi.NewBranchesAPIController(BranchesAPIService)
		CommitsAPIController := openapi.NewCommitsAPIController(CommitsAPIService)
		OrganizationsAPIController := openapi.NewOrganizationsAPIController(OrganizationsAPIService)
		ProjectsAPIController := openapi.NewProjectsAPIController(ProjectsAPIService)
		PromptsAPIController := openapi.NewPromptsAPIController(PromptsAPIService)
		TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)
		UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)
	*/

	// Initialize router
	router := openapi.NewRouter(
		HealthCheckAPIController,
		/*
			APIKeysAPIController,
			BranchesAPIController,
			CommitsAPIController,
			OrganizationsAPIController,
			ProjectsAPIController,
			PromptsAPIController,
			TeamsAPIController,
			UsersAPIController,
		*/
	)

	// Initialize server
	srv := &http.Server{
		Addr:    cfg.BindAddress,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Listen", zap.Error(err))
		}
	}()

	logger.Info("Server started",
		zap.String("bind_address", cfg.BindAddress),
		zap.String("environment", cfg.Environment))

	// Wait for interrupt signal to shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	// Set a timeout of 5 seconds for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", zap.Error(err))
	}

	logger.Info("Server exiting")
}
