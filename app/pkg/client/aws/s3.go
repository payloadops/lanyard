package aws_client

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
)

func GetS3Client() *s3.S3 {
	var once sync.Once
	var s3Client *s3.S3
	once.Do(func() {
		s3Client = s3.New(GetAWSSession())
	})
	return s3Client
}
