package main

import (
	"github.com/payloadops/plato/api/openapi"
	"github.com/payloadops/plato/api/service"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	APIKeysAPIService := service.NewAPIKeysAPIService()
	APIKeysAPIController := openapi.NewAPIKeysAPIController(APIKeysAPIService)

	BranchesAPIService := service.NewBranchesAPIService()
	BranchesAPIController := openapi.NewBranchesAPIController(BranchesAPIService)

	CommitsAPIService := service.NewCommitsAPIService()
	CommitsAPIController := openapi.NewCommitsAPIController(CommitsAPIService)

	HealthCheckAPIService := service.NewHealthCheckAPIService()
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)

	OrganizationsAPIService := service.NewOrganizationsAPIService()
	OrganizationsAPIController := openapi.NewOrganizationsAPIController(OrganizationsAPIService)

	ProjectsAPIService := service.NewProjectsAPIService()
	ProjectsAPIController := openapi.NewProjectsAPIController(ProjectsAPIService)

	PromptsAPIService := service.NewPromptsAPIService()
	PromptsAPIController := openapi.NewPromptsAPIController(PromptsAPIService)

	TeamsAPIService := service.NewTeamsAPIService()
	TeamsAPIController := openapi.NewTeamsAPIController(TeamsAPIService)

	UsersAPIService := service.NewUsersAPIService()
	UsersAPIController := openapi.NewUsersAPIController(UsersAPIService)

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

	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

func main() {
	awsConfig, err := LoadAWSConfig()
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(awsConfig)
	// Create S3 client
	s3Client := s3.NewFromConfig(awsConfig)
	// Create ElastiCache client
	elastiCacheClient := elasticache.NewFromConfig(awsConfig)

	// Example usage of the clients
	fmt.Println("DynamoDB, S3, and ElastiCache clients created successfully")

	// Your application logic here...
}

// LoadAWSConfig loads AWS configuration based on the environment
func LoadAWSConfig() (aws.Config, error) {
	region := os.Getenv("AWS_DEFAULT_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	endpointResolver := aws.EndpointResolverFunc(
		func(service, region string) (aws.Endpoint, error) {
			if os.Getenv("ENVIRONMENT") == "local" {
				switch service {
				case dynamodb.ServiceID:
					return aws.Endpoint{URL: os.Getenv("DYNAMODB_ENDPOINT")}, nil
				case s3.ServiceID:
					return aws.Endpoint{URL: os.Getenv("S3_ENDPOINT")}, nil
				case elasticache.ServiceID:
					return aws.Endpoint{URL: os.Getenv("ELASTICACHE_ENDPOINT")}, nil
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
		options = append(options, config.WithEndpointResolver(endpointResolver))
	}

	return config.LoadDefaultConfig(context.TODO(), options...)
}
*/
