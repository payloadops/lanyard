package awsclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Once   sync.Once
	s3Client *s3.Client
	initErr  error
)

func InitS3Client() (*s3.Client, error) {
	s3Once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"), config.WithClientLogMode(aws.LogRetries|aws.LogRequestWithBody|aws.LogResponseWithBody))
		if err != nil {
			initErr = fmt.Errorf("failed to load AWS config: %w", err)
			return
		}
		s3Client = s3.NewFromConfig(cfg)
	})
	return s3Client, initErr
}

func GetS3Client() *s3.Client {
	return s3Client
}
