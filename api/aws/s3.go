package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// InitializeS3Client initializes and returns an S3 client
func InitializeS3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	// Create the S3 service client
	client := s3.NewFromConfig(cfg)
	return client, nil
}
