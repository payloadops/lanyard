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

// shutdownTimeout represents the time to wait in graceful-shutdown before force exiting
const shutdownTimeout = 5 * time.Second

func main() {
	// Load config from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize service logging
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Flush buffered logs upon exiting
	defer logger.Sync()

	// Set global logger to use this implementation
	// zap.ReplaceGlobals(logger)

	// Load AWS/localstack config values
	// awsConfig, err := client.LoadAWSConfig(cfg)
	// if err != nil {
	//	logger.Fatal("Failed to initialize aws config", zap.Error(err))
	// }

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
		zap.String("address", cfg.BindAddress),
		zap.String("environment", string(cfg.Environment)))

	// Wait for interrupt signal to shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	// Set a timeout of 5 seconds for graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", zap.Error(err))
	}

	logger.Info("Server exiting")
}
